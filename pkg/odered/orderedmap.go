// Original source code - https://github.com/wk8/go-ordered-map
// - change to use generics
// - add Keys(), Values() functions
// - modify for thread safe
// edited by @ironpark 2022.02.02

package ordered

import (
	"constraints"
	"container/list"
	"sync"
)

type KeyAble interface {
	constraints.Integer | constraints.Float | string
}

type Set[K KeyAble, V any] map[K]*Pair[K, V]

type Pair[K KeyAble, V any] struct {
	Key     K
	Value   V
	element *list.Element
}

type OrderedMap[K KeyAble, V any] struct {
	pairs Set[K, V]
	list  *list.List
	lock  sync.RWMutex
}

// NewMap New creates a new OrderedMap.
func NewMap[K KeyAble, V any]() *OrderedMap[K, V] {
	return &OrderedMap[K, V]{
		pairs: make(Set[K, V]),
		list:  list.New(),
	}
}

// Get looks for the given key, and returns the value associated with it,
// or nil if not found. The boolean it returns says whether the key is present in the map.
func (om OrderedMap[K, V]) Get(key K) (V, bool) {
	om.lock.RLock()
	defer om.lock.RUnlock()
	if pair, present := om.pairs[key]; present {
		return pair.Value, present
	}
	var emptyValue V
	return emptyValue, false
}

// GetPair looks for the given key, and returns the pair associated with it,
// or nil if not found. The Pair struct can then be used to iterate over the ordered map
// from that point, either forward or backward.
func (om *OrderedMap[K, V]) GetPair(key K) *Pair[K, V] {
	om.lock.RLock()
	defer om.lock.RUnlock()
	return om.pairs[key]
}

// Set sets the key-value pair, and returns what `Get` would have returned
// on that key prior to the call to `Set`.
func (om *OrderedMap[K, V]) Set(key K, value V) (V, bool) {
	om.lock.Lock()
	defer om.lock.Unlock()
	if pair, present := om.pairs[key]; present {
		oldValue := pair.Value
		pair.Value = value
		return oldValue, true
	}

	pair := &Pair[K, V]{
		Key:   key,
		Value: value,
	}
	pair.element = om.list.PushBack(pair)
	om.pairs[key] = pair

	var emptyValue V
	return emptyValue, false
}

// Delete removes the key-value pair, and returns what `Get` would have returned
// on that key prior to the call to `Delete`.
func (om *OrderedMap[K, V]) Delete(key K) (V, bool) {
	om.lock.Lock()
	defer om.lock.Unlock()
	if pair, present := om.pairs[key]; present {
		om.list.Remove(pair.element)
		delete(om.pairs, key)
		return pair.Value, true
	}
	return nil, false
}

// Len returns the length of the ordered map.
func (om *OrderedMap[K, V]) Len() int {
	om.lock.RLock()
	defer om.lock.RUnlock()
	return om.list.Len()
}

// Oldest returns a pointer to the oldest pair. It's meant to be used to iterate on the ordered map's
// pairs from the oldest to the newest, e.g.:
// for pair := orderedMap.Oldest(); pair != nil; pair = pair.Next() { fmt.Printf("%v => %v\n", pair.Key, pair.Value) }
func (om *OrderedMap[K, V]) Oldest() *Pair {
	om.lock.RLock()
	defer om.lock.RUnlock()
	return listElementToPair(om.list.Front())
}

// Newest returns a pointer to the newest pair. It's meant to be used to iterate on the ordered map's
// pairs from the newest to the oldest, e.g.:
// for pair := orderedMap.Oldest(); pair != nil; pair = pair.Next() { fmt.Printf("%v => %v\n", pair.Key, pair.Value) }
func (om *OrderedMap[K, V]) Newest() *Pair {
	om.lock.RLock()
	defer om.lock.RUnlock()
	return listElementToPair(om.list.Back())
}

func (om *OrderedMap[K, V]) Keys() (keys []K) {
	om.lock.RLock()
	defer om.lock.RUnlock()
	keys = make([]K, om.list.Len())
	element := om.list.Front()
	for i := 0; element != nil; i++ {
		keys[i] = element.Value.(*Pair[K, V]).Key
		element = element.Next()
	}
	return keys
}

func (om *OrderedMap[K, V]) Values() (keys []K) {
	om.lock.RLock()
	defer om.lock.RUnlock()
	keys = make([]K, om.list.Len())
	element := om.list.Front()
	for i := 0; element != nil; i++ {
		keys[i] = element.Value.(*Pair[K, V]).Key
		element = element.Next()
	}
	return keys
}

// Next returns a pointer to the next pair.
func (p *Pair) Next() *Pair {
	return listElementToPair(p.element.Next())
}

// Previous returns a pointer to the previous pair.
func (p *Pair) Prev() *Pair {
	return listElementToPair(p.element.Prev())
}

func listElementToPair(element *list.Element) *Pair {
	if element == nil {
		return nil
	}
	return element.Value.(*Pair)
}
