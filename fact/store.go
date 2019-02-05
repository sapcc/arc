package fact

import (
	"fmt"
	"reflect"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
)

type FactSource interface {
	Name() string
	Facts() (map[string]interface{}, error)
}

type factSource struct {
	plugin FactSource
	facts  map[string]interface{}
	mutex  sync.RWMutex
	ticker <-chan time.Time
}

type Store struct {
	sources     sync.Map
	wg          sync.WaitGroup
	updateChan  chan (map[string]interface{})
	initialized bool
	mu          sync.Mutex
	interval    time.Duration
}

func NewStore() *Store {
	return &Store{
		initialized: false,
		interval:    30 * time.Minute,
	}
}

func (fs *Store) AddSource(plugin FactSource, interval time.Duration) {
	source := factSource{plugin: plugin, facts: make(map[string]interface{})}

	fs.sources.Store(plugin.Name(), &source)
	fs.wg.Add(1)
	go func() {
		//collect facts initially
		if err := fs.update(plugin.Name()); err != nil {
			log.Error(err)
			return
		}
		fs.wg.Done()
		if interval > 0 {
			tickChan := time.NewTicker(interval)
			source.ticker = tickChan.C
			for range source.ticker {
				if err := fs.update(plugin.Name()); err != nil {
					log.Error(err)
					tickChan.Stop()
				}
			}
		}
	}()
}

func (fs *Store) Wait() {
	fs.wg.Wait()
}

func (fs *Store) Updates() <-chan map[string]interface{} {
	fs.mu.Lock()
	defer fs.mu.Unlock()
	if fs.updateChan == nil {
		// setup update chanel
		fs.updateChan = make(chan map[string]interface{})
		go func() {
			fs.Wait()
			fs.updateChan <- fs.Facts()
			fs.initialized = true
		}()

		// setup the global full update ticker
		if fs.interval > 0 {
			go func() {
				fullUpdateTicker := time.Tick(fs.interval)
				for range fullUpdateTicker {
					fs.FullUpdate()
				}
			}()
		}
	}
	return fs.updateChan
}

func (fs *Store) Facts() map[string]interface{} {
	fs.Wait()

	//create a copy of the internal map
	facts := make(map[string]interface{})
	fs.sources.Range(func(_, s interface{}) bool {
		src := s.(*factSource)
		src.mutex.RLock()
		for fact, val := range src.facts {
			facts[fact] = val
		}
		src.mutex.RUnlock()
		return true
	})
	return facts
}

/*
* Return an error just when the fact is not known
 */
func (fs *Store) update(name string) error {
	s, ok := fs.sources.Load(name)
	if !ok {
		return fmt.Errorf("unknown fact source %s", name)
	}
	source := s.(*factSource)

	start := time.Now()
	facts, err := source.plugin.Facts()
	if err != nil {
		log.Warnf("Failed to update %s facts: %s", name, err)
		return err
	}
	source.mutex.Lock()
	unmodified := len(facts) == len(source.facts)
	for key, value := range facts {
		if unmodified {
			unmodified = reflect.DeepEqual(source.facts[key], value)
		}
		source.facts[key] = value
	}
	source.mutex.Unlock()
	log.Debugf("Updating %s facts took %s.", name, time.Since(start))
	if !unmodified {
		log.Debugf("%s facts changed.", name)
		if fs.initialized {
			fs.updateChan <- facts
		}
	}
	return nil
}

func (fs *Store) FullUpdate() {
	if fs.initialized {
		fs.mu.Lock()
		f := fs.Facts()
		fs.mu.Unlock()
		fs.updateChan <- f
	}
}
