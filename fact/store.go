package fact

import (
	"fmt"
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
}

type Store struct {
	sources map[string]factSource
	wg      sync.WaitGroup
}

func NewStore() *Store {
	return &Store{
		sources: make(map[string]factSource),
	}
}

func (fs *Store) AddSource(plugin FactSource, interval time.Duration) {
	fs.sources[plugin.Name()] = factSource{plugin: plugin, facts: make(map[string]interface{})}
	fs.wg.Add(1)
	//collect facts initially
	go func() {
		fs.update(plugin.Name())
		fs.wg.Done()
		if interval > 0 {
			for {
				select {
				case <-time.After(interval):
					fs.update(plugin.Name())
				}
			}
		}
	}()
}

func (fs *Store) Wait() {
	fs.wg.Wait()
}

//func (fs *Store) Fact(name string) (string, error) {
//  fs.mutex.RLock()
//  val, ok := fs.facts[name]
//  fs.mutex.RUnlock()
//  if !ok {
//    return "", fmt.Errorf("fact %s not found", name)
//  }
//  return val, nil
//}

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

func (fs *Store) update(source string) error {

	src, ok := fs.sources[source]
	if !ok {
		return fmt.Errorf("Unknown fact source %s", source)
	}

	start := time.Now()
	facts, err := src.plugin.Facts()
	if err != nil {
		log.Warn("Failed to update %s fact source: %s", source, err)
		return err
	}
	src.mutex.Lock()
	//updated := len(facts) != len(src.facts)
	for key, value := range facts {
		src.facts[key] = value
	}
	src.mutex.Unlock()
	log.Debugf("Updating fact source %s took %s", source, time.Since(start))
	return nil
}
