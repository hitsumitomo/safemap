package safemap

import (
	"cmp"
	"container/list"
	"maps"
	"sync"
)

const Ordered = true

// M is a thread-safe map with additional features.
type M[K cmp.Ordered, V any] struct {
	sync.RWMutex
	m        map[K]V
	ordered  bool
	keys     *list.List
	keyIndex map[K]*list.Element
}

// New creates a new M.
// If the ordered parameter is set to true, the map will be ordered in the order of insertion.
func New[K cmp.Ordered, V any](ordered ...bool) *M[K, V] {
	sm := &M[K, V]{
		m: make(map[K]V),
		keys: list.New(),
		keyIndex: make(map[K]*list.Element),
	}

	if len(ordered) > 0 && ordered[0] {
		sm.ordered = true
	}
	return sm
}

// Exists - checks if a key exists in the map
func (sm *M[K, V]) Exists(key K) (exists bool) {
	sm.RLock()
	defer sm.RUnlock()

	_, exists = sm.m[key]
	return exists
}

// Load - loads a value from the map
func (sm *M[K, V]) Load(key K) (value V, ok bool) {
	sm.RLock()
	defer sm.RUnlock()

	value, ok = sm.m[key]
	return value, ok
}

// Store - stores a value in the map
func (sm *M[K, V]) Store(key K, value V) {
	sm.Lock()
	defer sm.Unlock()

    if elem, exists := sm.keyIndex[key]; !exists {
        sm.keyIndex[key] = sm.keys.PushBack(key)
    } else {
        elem.Value = key
    }
    sm.m[key] = value
}

// Delete - deletes a key from the map
func (sm *M[K, V]) Delete(key K) {
	sm.Lock()
	defer sm.Unlock()

	sm.removeKey(key)
}

// removeKey - removes a key from the map (internal use).
func (sm *M[K, V]) removeKey(key K) {
	if elem, exists := sm.keyIndex[key]; exists {
		sm.keys.Remove(elem)
		delete(sm.keyIndex, key)
	}
	delete(sm.m, key)
}

// Add - adds a value to the existing value for numeric types
func (sm *M[K, V]) Add(key K, value V) {
	sm.Lock()
	defer sm.Unlock()

	switch v := any(sm.m[key]).(type) {
	case int:
		sm.m[key] = any(v + any(value).(int)).(V)
		// sm.m[key] = any(v + any(value).(int)).(V)
	case int8:
		sm.m[key] = any(v + any(value).(int8)).(V)
	case int16:
		sm.m[key] = any(v + any(value).(int16)).(V)
	case int32:
		sm.m[key] = any(v + any(value).(int32)).(V)
	case int64:
		sm.m[key] = any(v + any(value).(int64)).(V)
	case uint:
		sm.m[key] = any(v + any(value).(uint)).(V)
	case uint8:
		sm.m[key] = any(v + any(value).(uint8)).(V)
	case uint16:
		sm.m[key] = any(v + any(value).(uint16)).(V)
	case uint32:
		sm.m[key] = any(v + any(value).(uint32)).(V)
	case uint64:
		sm.m[key] = any(v + any(value).(uint64)).(V)
	case float32:
		sm.m[key] = any(v + any(value).(float32)).(V)
	case float64:
		sm.m[key] = any(v + any(value).(float64)).(V)
	case string:
		sm.m[key] = any(v + any(value).(string)).(V)
	default:
		panic("unsupported type")
	}
}

// Sub - subtracts a value from the existing value for numeric types
func (sm *M[K, V]) Sub(key K, value V) {
	sm.Lock()
	defer sm.Unlock()

	switch v := any(sm.m[key]).(type) {
	case int:
		sm.m[key] = any(v - any(value).(int)).(V)
	case int8:
		sm.m[key] = any(v - any(value).(int8)).(V)
	case int16:
		sm.m[key] = any(v - any(value).(int16)).(V)
	case int32:
		sm.m[key] = any(v - any(value).(int32)).(V)
	case int64:
		sm.m[key] = any(v - any(value).(int64)).(V)
	case uint:
		sm.m[key] = any(v - any(value).(uint)).(V)
	case uint8:
		sm.m[key] = any(v - any(value).(uint8)).(V)
	case uint16:
		sm.m[key] = any(v - any(value).(uint16)).(V)
	case uint32:
		sm.m[key] = any(v - any(value).(uint32)).(V)
	case uint64:
		sm.m[key] = any(v - any(value).(uint64)).(V)
	case float32:
		sm.m[key] = any(v - any(value).(float32)).(V)
	case float64:
		sm.m[key] = any(v - any(value).(float64)).(V)
	default:
		panic("unsupported type")
	}
}

// LoadAndDelete - loads a value from the map and deletes the key
func (sm *M[K, V]) LoadAndDelete(key K) (value V, ok bool) {
	sm.Lock()
	defer sm.Unlock()

	value, ok = sm.m[key]
	if ok {
		sm.removeKey(key)
	}
	return value, ok
}

// LoadOrStore - loads a value from the map or stores a new value
func (sm *M[K, V]) LoadOrStore(key K, value V) (actual V, loaded bool) {
	sm.Lock()
	defer sm.Unlock()

	actual, loaded = sm.m[key]
	if loaded {
		return actual, true
	}

	sm.m[key] = value
	if sm.ordered {
		if elem, exists := sm.keyIndex[key]; !exists {
			sm.keyIndex[key] = sm.keys.PushBack(key)
		} else {
			elem.Value = key
		}
	}
	return value, false
}

// Swap - swaps a value in the map
func (sm *M[K, V]) Swap(key K, value V) (previous V, loaded bool) {
	sm.Lock()
	defer sm.Unlock()

	previous, loaded = sm.m[key]
	sm.m[key] = value
	return previous, loaded
}

// Range - ranges over the map by calling a function for each key-value pair
func (sm *M[K, V]) Range(f func(K, V) bool) {
	sm.RLock()
	defer sm.RUnlock()

	if sm.ordered {
		for e := sm.keys.Front(); e != nil; e = e.Next() {
			k := e.Value.(K)
			if !f(k, sm.m[k]) {
				return
			}
		}
		return
	}

	for k, v := range sm.m {
		if !f(k, v) {
			return
		}
	}
}

// RangeDelete - ranges over the map by calling a function for each key-value pair
// and deletes the key if the function returns -1
func (sm *M[K, V]) RangeDelete(f func(K, V) int8) {
	sm.Lock()
	defer sm.Unlock()

	if sm.ordered {
		for e := sm.keys.Front(); e != nil; e = e.Next() {
			k := e.Value.(K)
			if f(k, sm.m[k]) == 0 {
				return

			} else if f(k, sm.m[k]) == -1 {
				sm.removeKey(k)
			}
		}
		return
	}

	keysToDelete := make([]K, 0)
	defer func() {
		for _, k := range keysToDelete {
			sm.removeKey(k)
		}
	}()

	for k, v := range sm.m {
		if f(k, v) == 0 {
			return

		} else if f(k, v) == -1 {
			keysToDelete = append(keysToDelete, k)
		}
	}
}

// Clear - clears the map
func (sm *M[K, V]) Clear() {
	sm.Lock()
	defer sm.Unlock()

	sm.m = make(map[K]V)
	sm.keys.Init()
	sm.keyIndex = make(map[K]*list.Element)
}

// Len - returns the length of the map
func (sm *M[K, V]) Len() int {
	sm.RLock()
	defer sm.RUnlock()

	return len(sm.m)
}

// Keys - returns the keys of the map
func (sm *M[K, V]) Keys() []K {
	sm.RLock()
	defer sm.RUnlock()

	keys := make([]K, sm.keys.Len())
	i := 0
	if sm.ordered {
		for e := sm.keys.Front(); e != nil; e = e.Next() {
			keys[i] = e.Value.(K)
			i++
		}
	} else {
		for k := range sm.m {
			keys[i] = k
			i++
		}
	}
	return keys
}

// Map - returns a clone of the underlying map.
func (sm *M[K, V]) Map() map[K]V {
	sm.RLock()
	defer sm.RUnlock()

	return maps.Clone(sm.m)
}

// KeysInRange - returns the slice of keys of the map in a range
// Note: works only for positive numeric types
func (sm *M[K, V]) KeysInRange(from, to K) []K {
	sm.RLock()
	defer sm.RUnlock()

	var zero K
	filtered := make([]K, 0, len(sm.m))

	if sm.ordered {
		for e := sm.keys.Front(); e != nil; e = e.Next() {
			k := e.Value.(K)
			if (from == zero || k >= from) && (to == zero || k < to) {
				filtered = append(filtered, k)
			}
		}
		return filtered
	}

	for k := range sm.m {
		if (from == zero || k >= from) && (to == zero || k < to) {
			filtered = append(filtered, k)
		}
	}
	return filtered
}
