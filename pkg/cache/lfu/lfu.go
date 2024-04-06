package lfu

import (
	"banner_service/pkg/cache/internal/dlist"
	"sync"
	"time"
)

const DefaultCacheCapacity = 1024

type LFUCache[K comparable, V any] struct {
	sync.Mutex
	data         map[K]*dlist.Node[K, V]
	freq         map[int]*dlist.DLinkedList[K, V]
	minFrequency int
	len          int
	cap          int
}

// New создает менеджер кеша на базе алгоритма LFU с ёмкостью cap
func New[K comparable, V any](cap int) *LFUCache[K, V] {
	if cap <= 0 {
		cap = DefaultCacheCapacity
	}

	return &LFUCache[K, V]{
		data: make(map[K]*dlist.Node[K, V], cap),
		freq: make(map[int]*dlist.DLinkedList[K, V]),
		cap:  cap,
	}
}

type Updater[K comparable, V any] interface {
	GetNewValue(K) (V, bool)
}

// New создает менеджер кеша на базе алгоритма LFU с ёмкостью cap, который будет обновляться,
// получая данные из updater раз в интервал updateInterval
func NewWithUpdateInterval[K comparable, V any](cap int, updater Updater[K, V], updateInterval time.Duration) *LFUCache[K, V] {
	c := New[K, V](cap)
	go c.updateWorker(updater, updateInterval)
	return c
}

// Set обновляет значение по ключу или создает новую пару ключ-значение, если такой не существует.
// Если размер кеша равен его ёмкости, то вытесняется наименее часто используемый элемент,
// в случае, если таких элементов несколько, то вытесняется элемент по правилу FIFO
func (c *LFUCache[K, V]) Set(key K, val V) {
	c.Lock()
	defer c.Unlock()

	if node, ok := c.data[key]; ok {
		node.Value = val
		c.update(node)
		return
	}

	node := dlist.NewNode(key, val)
	node.Frequency = 1

	if c.len == c.cap {
		deleted := c.freq[c.minFrequency].Pop(nil)
		delete(c.data, deleted.Key)
		c.len--
	}

	c.add(node)
	c.minFrequency = 1
	c.data[key] = node
	c.len++
}

// Get возвращает элемент и обновляет его частоту запросов, если он находится в кеше
func (c *LFUCache[K, V]) Get(key K) (value V, ok bool) {
	c.Lock()
	defer c.Unlock()

	if _, ok := c.data[key]; !ok {
		return value, false
	}

	node := c.data[key]
	c.update(node)

	return node.Value, true
}

// GetFrequency возвращает частоту обращения к элементу
func (c *LFUCache[K, V]) GetFrequency(key K) (freq int, ok bool) {
	c.Lock()
	defer c.Unlock()

	if _, ok := c.data[key]; !ok {
		return 0, false
	}

	return c.data[key].Frequency, true
}

func (c *LFUCache[K, V]) Delete(key K) bool {
	if _, ok := c.data[key]; !ok {
		return false
	}

	deleted := c.freq[c.minFrequency].Pop(c.data[key])
	delete(c.data, deleted.Key)
	c.len--

	return true
}

// Len возвращает текущий размер кеша
func (c *LFUCache[K, V]) Len() int {
	return c.len
}

// Cap возвращает ёмкость кеша
func (c *LFUCache[K, V]) Cap() int {
	return c.cap
}

func (c *LFUCache[K, V]) update(node *dlist.Node[K, V]) {
	c.freq[node.Frequency].Pop(node)
	if c.freq[node.Frequency].Len() == 0 {
		c.minFrequency++
	}
	node.Frequency++
	c.add(node)
}

func (c *LFUCache[K, V]) add(node *dlist.Node[K, V]) {
	if _, ok := c.freq[node.Frequency]; !ok {
		c.freq[node.Frequency] = dlist.New[K, V]()
	}
	c.freq[node.Frequency].Append(node)
}

func (c *LFUCache[K, V]) updateWorker(updater Updater[K, V], updateInterval time.Duration) {
	update := time.After(updateInterval)

	for {
		select {
		case <-update:
			for k := range c.data {
				v, ok := updater.GetNewValue(k)
				if !ok {
					c.Delete(k)
					continue
				}
				c.data[k].Value = v
			}
			update = time.After(updateInterval)

		}
	}
}
