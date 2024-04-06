package dlist

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDListAppend(t *testing.T) {
	tcases := []struct {
		name           string
		init           func() *DLinkedList[int, int]
		expectedValues []int
		expectedLen    int
	}{
		{
			name: "Empty list",
			init: func() *DLinkedList[int, int] {
				return New[int, int]()
			},
			expectedValues: []int{},
		},
		{
			name: "Correct",
			init: func() *DLinkedList[int, int] {
				l := New[int, int]()
				for i := 0; i < 5; i++ {
					l.Append(NewNode(i, i))
				}
				return l
			},
			expectedValues: []int{0, 1, 2, 3, 4},
			expectedLen:    5,
		},
	}

	for _, tc := range tcases {
		t.Run(tc.name, func(t *testing.T) {
			l := tc.init()
			actualValues := l.Values()
			actualLen := l.Len()
			assert.Equal(t, tc.expectedValues, actualValues)
			assert.Equal(t, tc.expectedLen, actualLen)
		})
	}
}

func TestDListPop(t *testing.T) {
	tcases := []struct {
		name           string
		init           func() *DLinkedList[int, int]
		expectedValues []int
		expectedLen    int
	}{
		{
			name: "Empty list",
			init: func() *DLinkedList[int, int] {
				l := New[int, int]()
				l.Pop(nil)
				return l
			},
			expectedValues: []int{},
		},
		{
			name: "Correct",
			init: func() *DLinkedList[int, int] {
				l := New[int, int]()
				for i := 0; i < 5; i++ {
					l.Append(NewNode(i, i))
				}
				l.Pop(nil)
				return l
			},
			expectedValues: []int{1, 2, 3, 4},
			expectedLen:    4,
		},
		{
			name: "All elements",
			init: func() *DLinkedList[int, int] {
				l := New[int, int]()
				for i := 0; i < 5; i++ {
					l.Append(NewNode(i, i))
				}
				for i := 0; i < 5; i++ {
					l.Pop(nil)
				}
				return l
			},
			expectedValues: []int{},
		},
		{
			name: "First element",
			init: func() *DLinkedList[int, int] {
				l := New[int, int]()
				n := NewNode(100, 100)
				l.Append(n)
				for i := 0; i < 5; i++ {
					l.Append(NewNode(i, i))
				}
				l.Pop(n)
				return l
			},
			expectedValues: []int{0, 1, 2, 3, 4},
			expectedLen:    5,
		},
		{
			name: "Middle element",
			init: func() *DLinkedList[int, int] {
				l := New[int, int]()
				n := NewNode(100, 100)
				for i := 0; i < 2; i++ {
					l.Append(NewNode(i, i))
				}
				l.Append(n)
				for i := 2; i < 5; i++ {
					l.Append(NewNode(i, i))
				}
				l.Pop(n)
				return l
			},
			expectedValues: []int{0, 1, 2, 3, 4},
			expectedLen:    5,
		},
		{
			name: "Last element",
			init: func() *DLinkedList[int, int] {
				l := New[int, int]()
				n := NewNode(100, 100)
				for i := 0; i < 5; i++ {
					l.Append(NewNode(i, i))
				}
				l.Append(n)
				l.Pop(n)
				return l
			},
			expectedValues: []int{0, 1, 2, 3, 4},
			expectedLen:    5,
		},
	}

	for _, tc := range tcases {
		t.Run(tc.name, func(t *testing.T) {
			l := tc.init()
			actualValues := l.Values()
			actualLen := l.Len()
			assert.Equal(t, tc.expectedValues, actualValues)
			assert.Equal(t, tc.expectedLen, actualLen)
		})
	}
}
