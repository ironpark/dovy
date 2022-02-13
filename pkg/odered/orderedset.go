// Original source code - https://github.com/wk8/go-ordered-map
// - change to use generics
// - add Keys(), Values() functions
// edited by @ironpark 2022.02.02

package ordered

type OrderedSet[V KeyAble] struct {
	set *OrderedMap[V, struct{}]
}

// NewSet New creates a new OrderedSet.
func NewSet[V KeyAble]() *OrderedSet[V] {
	return &OrderedSet[V]{
		set: NewMap[V, struct{}](),
	}
}

func (om OrderedSet[V]) Exist(key V) bool {
	_, ok := om.set.Get(key)
	return ok
}

func (om *OrderedSet[V]) Set(key V) bool {
	_, ok := om.set.Set(key, struct{}{})
	return ok
}

func (om *OrderedSet[V]) Delete(key V) bool {
	_, ok := om.set.Delete(key)
	return ok
}

func (om *OrderedSet[K, V]) Len() int {
	return om.set.Len()
}

func (om *OrderedSet[K, V]) Values() (keys []K) {
	return om.set.Keys()
}
