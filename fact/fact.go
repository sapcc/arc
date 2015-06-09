package fact

import (
	"fmt"
	"sync"
	"time"
)

type FactSource interface {
	Name() string
	Facts() (map[string]string, error)
}

type Store struct {
	facts   map[string]string
	sources map[string]FactSource
	mutex   sync.RWMutex
	wg      sync.WaitGroup
}

func NewStore() *Store {
	return &Store{
		facts:   make(map[string]string),
		sources: make(map[string]FactSource),
	}
}

func (fs *Store) AddSource(source FactSource, interval time.Duration) {
	fs.sources[source.Name()] = source
	fs.wg.Add(1)
	//collect facts initially
	go func() {
		fs.update(source.Name())
		fs.wg.Done()
	}()
}

func (fs *Store) Wait() {
	fs.wg.Wait()
}

func (fs *Store) Fact(name string) (string, error) {
	fs.mutex.RLock()
	val, ok := fs.facts[name]
	fs.mutex.RUnlock()
	if !ok {
		return "", fmt.Errorf("fact %s not found", name)
	}
	return val, nil
}

func (fs *Store) Facts() map[string]string {
	fs.Wait()

	//create a copy of the internal map
	facts := make(map[string]string)
	fs.mutex.RLock()
	for fact, val := range fs.facts {
		facts[fact] = val
	}
	fs.mutex.RUnlock()
	return facts
}

func (fs *Store) update(source string) error {

	s, ok := fs.sources[source]
	if !ok {
		return fmt.Errorf("Unknown fact source %s", source)
	}
	facts, err := s.Facts()
	if err != nil {
		return err
	}
	fs.mutex.Lock()
	for key, value := range facts {
		fs.facts[key] = value
	}
	fs.mutex.Unlock()
	return nil
}
