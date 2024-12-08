package dog

import (
	"maps"
	"slices"
	"sync"

	"github.com/google/uuid"
)

type Dogs struct {
	mu    sync.RWMutex
	store map[string]Dog
}

func NewStore() *Dogs {
	return &Dogs{
		store: make(map[string]Dog),
	}
}

func (d *Dogs) Add(name, breed string) Dog {
	d.mu.Lock()
	defer d.mu.Unlock()

	id := uuid.New().String()

	newDog := Dog{
		ID:    id,
		Name:  name,
		Breed: breed,
	}

	d.store[id] = newDog

	return newDog
}

func (d *Dogs) Update(id, name, breed string) bool {
	d.mu.Lock()
	defer d.mu.Unlock()

	_, ok := d.store[id]
	if !ok {
		return false
	}

	updated := Dog{
		ID:    id,
		Name:  name,
		Breed: breed,
	}

	d.store[id] = updated

	return true
}

func (d *Dogs) GetAll() []Dog {
	d.mu.RLock()
	defer d.mu.RUnlock()

	return slices.Collect(maps.Values(d.store))
}

func (d *Dogs) Get(id string) (Dog, bool) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	dog, ok := d.store[id]

	return dog, ok
}

func (d *Dogs) Delete(id string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	delete(d.store, id)
}

type Dog struct {
	ID    string
	Name  string
	Breed string
}
