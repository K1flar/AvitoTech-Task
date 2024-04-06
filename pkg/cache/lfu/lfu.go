package lfu

import "banner_service/pkg/cache/internal/dlist"

const DefaultCacheCapacity = 1024

type LFUCache[K comparable, V any] struct {
	data map[K]dlist.Node[K, V]
	len  int
	cap  int
}

func New[K comparable, V any](cap int) *LFUCache[K, V] {
	if cap <= 0 {
		cap = DefaultCacheCapacity
	}

	return &LFUCache[K, V]{
		data: make(map[K]dlist.Node[K, V], cap),
		len:  0,
		cap:  cap,
	}
}

func (c *LFUCache[K, V]) Set(key K, val V) {

}
