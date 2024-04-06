package dlist

type Node[K comparable, V any] struct {
	Key   K
	Value V
	Prev  *Node[K, V]
	Next  *Node[K, V]
}

func NewNode[K comparable, V any](key K, value V) *Node[K, V] {
	return &Node[K, V]{
		Key:   key,
		Value: value,
	}
}

type DLinkedList[K comparable, V any] struct {
	sentinel *Node[K, V]
	len      int
}

func New[K comparable, V any]() *DLinkedList[K, V] {
	sentinel := &Node[K, V]{}
	sentinel.Prev = sentinel
	sentinel.Next = sentinel

	return &DLinkedList[K, V]{
		sentinel: sentinel,
		len:      0,
	}
}

func (l *DLinkedList[K, V]) Append(node *Node[K, V]) {
	node.Next = l.sentinel
	node.Prev = l.sentinel.Prev
	l.sentinel.Prev = node
	node.Prev.Next = node
	l.len++
}

func (l *DLinkedList[K, V]) Pop(node *Node[K, V]) {
	if l.len == 0 {
		return
	}

	if node == nil {
		node = l.sentinel.Next
	}

	node.Prev.Next = node.Next
	node.Next.Prev = node.Prev
	l.len--
}

func (l *DLinkedList[K, V]) Len() int {
	return l.len
}

func (l *DLinkedList[K, V]) Values() []V {
	v := make([]V, l.len)
	cur := l.sentinel.Next

	for i := 0; i < l.len; i++ {
		v[i] = cur.Value
		cur = cur.Next
	}

	return v
}
