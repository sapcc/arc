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
	sources     map[string]factSource
	wg          sync.WaitGroup
	updateChan  chan (map[string]interface{})
	initialized bool
	mu          sync.Mutex
	interval    time.Duration
}

func NewStore() *Store {
	return &Store{
		sources:     make(map[string]factSource),
		initialized: false,
		interval:    30 * time.Minute,
	}
}

func (fs *Store) AddSource(plugin FactSource, interval time.Duration) {
	source := factSource{plugin: plugin, facts: make(map[string]interface{})}

	fs.sources[plugin.Name()] = source
	fs.wg.Add(1)
	go func() {
		//collect facts initially
		fs.update(plugin.Name())
		fs.wg.Done()
		if interval > 0 {
			source.ticker = time.Tick(interval)
			for range source.ticker {
				fs.update(plugin.Name())
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
	for _, src := range fs.sources {
		src.mutex.RLock()
		for fact, val := range src.facts {
			facts[fact] = val
		}
		src.mutex.RUnlock()
	}
	return facts
}

func (fs *Store) update(name string) error {
	source, ok := fs.sources[name]
	if !ok {
		return fmt.Errorf("Unknown fact source %s", name)
	}

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
	return
}
