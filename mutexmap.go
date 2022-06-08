package mutexmap

import (
	"sync"
)

// MutexMap provides a go standard map protected by a RWMutex. It exposes the same methods
// as sync.Map but with better type safety and performance in scenarios with up to four cores.
// Type parameter K for the map's key and type parameter V for the map's values.
// For benchmarks compare: https://medium.com/@deckarep/the-new-kid-in-town-gos-sync-map-de24a6bf7c2c
type MutexMap[K comparable, V any] struct {
	sync.RWMutex
	store map[K]V
}

func NewMutexMap[K comparable, V any]() *MutexMap[K, V] {
	m := &MutexMap[K, V]{store: make(map[K]V)}

	return m
}

func (m *MutexMap[K, V]) Store(key K, val V) {
	m.Lock()
	defer m.Unlock()

	m.store[key] = val
}

func (m *MutexMap[K, V]) Load(key K) (V, bool) {
	m.RLock()
	defer m.RUnlock()

	val, ok := m.store[key]
	return val, ok
}

func (m *MutexMap[K, _]) Delete(key K) {
	m.Lock()
	defer m.Unlock()

	delete(m.store, key)
}

func (m *MutexMap[K, V]) LoadAndDelete(key K) (V, bool) {
	m.Lock()
	defer m.Unlock()

	val, loaded := m.store[key]
	if loaded {
		delete(m.store, key)
	}
	return val, loaded
}

func (m *MutexMap[K, V]) LoadOrStore(key K, val V) (V, bool) {
	m.Lock()
	defer m.Unlock()

	if val, ok := m.store[key]; ok {
		return val, ok
	}

	m.store[key] = val
	return val, false
}

func (m *MutexMap[K, V]) Range(f func(key K, value V) bool) {
	m.RLock()
	keys := make([]K, 0, len(m.store))
	for k := range m.store {
		keys = append(keys, k)
	}
	m.RUnlock()

	for _, k := range keys {
		v, ok := m.Load(k)
		if !ok {
			continue
		}
		if !f(k, v) {
			break
		}
	}
}
