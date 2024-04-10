package lfu

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLFUSet(t *testing.T) {
	type Element struct {
		Key       int
		Val       int
		Frequency int
	}
	tcases := []struct {
		name          string
		cap           int
		init          func(cap int) *LFUCache[int, int]
		expectedLen   int
		expectedCap   int
		expectedPairs []Element
	}{
		{
			name: "Negative cap",
			cap:  -1,
			init: func(cap int) *LFUCache[int, int] {
				return New[int, int](cap)
			},
			expectedCap: DefaultCacheCapacity,
		},
		{
			name: "Empty cache",
			cap:  10,
			init: func(cap int) *LFUCache[int, int] {
				return New[int, int](cap)
			},
			expectedCap: 10,
		},
		{
			name: "Correct",
			cap:  10,
			init: func(cap int) *LFUCache[int, int] {
				c := New[int, int](cap)
				c.Set(1, 1)
				return c
			},
			expectedLen:   1,
			expectedCap:   10,
			expectedPairs: []Element{{1, 1, 1}},
		},
		{
			name: "Full",
			cap:  5,
			init: func(cap int) *LFUCache[int, int] {
				c := New[int, int](cap)
				for i := 0; i < cap; i++ {
					c.Set(i, i)
				}
				return c
			},
			expectedLen:   5,
			expectedCap:   5,
			expectedPairs: []Element{{0, 0, 1}, {1, 1, 1}, {2, 2, 1}, {3, 3, 1}, {4, 4, 1}},
		},
		{
			name: "Displacing the first",
			cap:  5,
			init: func(cap int) *LFUCache[int, int] {
				c := New[int, int](cap)
				for i := 0; i < cap+1; i++ {
					c.Set(i, i)
				}
				return c
			},
			expectedLen:   5,
			expectedCap:   5,
			expectedPairs: []Element{{1, 1, 1}, {2, 2, 1}, {3, 3, 1}, {4, 4, 1}, {5, 5, 1}},
		},
		{
			name: "Updating existing",
			cap:  5,
			init: func(cap int) *LFUCache[int, int] {
				c := New[int, int](cap)
				for i := 0; i < cap; i++ {
					c.Set(i, i)
				}
				c.Set(3, 100)
				return c
			},
			expectedLen:   5,
			expectedCap:   5,
			expectedPairs: []Element{{0, 0, 1}, {1, 1, 1}, {2, 2, 1}, {3, 100, 2}, {4, 4, 1}},
		},
	}

	for _, tc := range tcases {
		t.Run(tc.name, func(t *testing.T) {
			c := tc.init(tc.cap)
			assert.Equal(t, tc.expectedLen, c.Len())
			assert.Equal(t, tc.expectedCap, c.Cap())

			for _, e := range tc.expectedPairs {
				assert.Contains(t, c.data, e.Key)
				assert.Equal(t, c.data[e.Key].Value, e.Val)
				assert.Equal(t, c.data[e.Key].Frequency, e.Frequency)
			}
		})
	}
}

func TestLFUGet(t *testing.T) {
	tcases := []struct {
		name              string
		cap               int
		key               int
		init              func(cap int) *LFUCache[int, int]
		expectedVal       int
		expectedExst      bool
		expectedFrequency int
	}{
		{
			name: "Non-existent element",
			cap:  5,
			key:  1000,
			init: func(cap int) *LFUCache[int, int] {
				return New[int, int](cap)
			},
		},
		{
			name: "Correct",
			cap:  5,
			key:  3,
			init: func(cap int) *LFUCache[int, int] {
				c := New[int, int](cap)
				for i := 0; i < cap; i++ {
					c.Set(i, i)
				}
				return c
			},
			expectedVal:       3,
			expectedExst:      true,
			expectedFrequency: 2,
		},
		{
			name: "After accessing",
			cap:  5,
			key:  3,
			init: func(cap int) *LFUCache[int, int] {
				c := New[int, int](cap)
				for i := 0; i < cap; i++ {
					c.Set(i, i)
				}
				c.Set(3, 100)
				c.Set(3, 100)
				c.Set(3, 1000)

				return c
			},
			expectedVal:       1000,
			expectedExst:      true,
			expectedFrequency: 5,
		},
	}

	for _, tc := range tcases {
		t.Run(tc.name, func(t *testing.T) {
			c := tc.init(tc.cap)
			actualVal, actualExst := c.Get(tc.key)
			assert.Equal(t, tc.expectedVal, actualVal)
			assert.Equal(t, tc.expectedExst, actualExst)

			if tc.expectedExst {
				assert.Contains(t, c.data, tc.key)
				assert.Equal(t, c.data[tc.key].Frequency, tc.expectedFrequency)
			}
		})
	}
}

func TestLFUGetFrequency(t *testing.T) {
	tcases := []struct {
		name              string
		cap               int
		key               int
		init              func(cap int) *LFUCache[int, int]
		expectedFrequency int
	}{
		{
			name: "Non-existent element",
			cap:  5,
			key:  1000,
			init: func(cap int) *LFUCache[int, int] {
				return New[int, int](cap)
			},
		},
		{
			name: "Correct",
			cap:  5,
			key:  3,
			init: func(cap int) *LFUCache[int, int] {
				c := New[int, int](cap)
				for i := 0; i < cap; i++ {
					c.Set(i, i)
				}
				return c
			},
			expectedFrequency: 1,
		},
		{
			name: "After accessing",
			cap:  5,
			key:  3,
			init: func(cap int) *LFUCache[int, int] {
				c := New[int, int](cap)
				for i := 0; i < cap; i++ {
					c.Set(i, i)
				}
				c.Get(3)
				c.Get(3)
				c.Get(3)

				return c
			},
			expectedFrequency: 4,
		},
	}

	for _, tc := range tcases {
		t.Run(tc.name, func(t *testing.T) {
			c := tc.init(tc.cap)
			actualFrequency := c.GetFrequency(tc.key)
			assert.Equal(t, actualFrequency, tc.expectedFrequency)
		})
	}
}

func TestLFUDelete(t *testing.T) {
	tcases := []struct {
		name        string
		cap         int
		key         int
		init        func(cap int) *LFUCache[int, int]
		expectedRes bool
	}{
		{
			name: "Non-existent element",
			cap:  5,
			key:  1000,
			init: func(cap int) *LFUCache[int, int] {
				return New[int, int](cap)
			},
		},
		{
			name: "Correct",
			cap:  5,
			key:  3,
			init: func(cap int) *LFUCache[int, int] {
				c := New[int, int](cap)
				for i := 0; i < cap; i++ {
					c.Set(i, i)
				}
				return c
			},
			expectedRes: true,
		},
	}

	for _, tc := range tcases {
		t.Run(tc.name, func(t *testing.T) {
			c := tc.init(tc.cap)
			actualRes := c.Delete(tc.key)
			assert.Equal(t, tc.expectedRes, actualRes)
			assert.NotContains(t, c.data, tc.key)
		})
	}
}

func TestLFUUpdateWorker(t *testing.T) {
	type Pair struct {
		Key int
		Val int
	}
	tcases := []struct {
		name             string
		cap              int
		valueLifeTime    time.Duration
		timeForSleep     time.Duration
		init             func(cap int, valueLifeTime time.Duration) *LFUCacheWithLifeCycle[int, int]
		expectedLen      int
		expectedElements []Pair
	}{
		{
			name:          "Correct",
			cap:           5,
			valueLifeTime: time.Millisecond,
			timeForSleep:  2 * time.Millisecond,
			init: func(cap int, valueLifeTime time.Duration) *LFUCacheWithLifeCycle[int, int] {
				c := NewWithLifeCycle[int, int](cap, valueLifeTime)
				for i := 0; i < cap; i++ {
					c.Set(i, i)
				}
				return c
			},
			expectedLen:      0,
			expectedElements: []Pair{},
		},
		{
			name:          "Without one",
			cap:           5,
			valueLifeTime: time.Second,
			timeForSleep:  time.Second/2 + time.Millisecond,
			init: func(cap int, valueLifeTime time.Duration) *LFUCacheWithLifeCycle[int, int] {
				c := NewWithLifeCycle[int, int](cap, valueLifeTime)
				for i := 0; i < cap-1; i++ {
					c.Set(i, i)
				}
				time.Sleep(time.Second / 2)
				c.Set(5, 5)
				return c
			},
			expectedLen:      1,
			expectedElements: []Pair{{5, 5}},
		},
	}

	for _, tc := range tcases {
		t.Run(tc.name, func(t *testing.T) {
			c := tc.init(tc.cap, tc.valueLifeTime)
			time.Sleep(tc.timeForSleep)

			assert.Equal(t, tc.expectedLen, c.Len())
			for _, p := range tc.expectedElements {
				assert.Contains(t, c.data, p.Key)
				assert.Equal(t, c.data[p.Key].Value, p.Val)
			}
		})
	}
}
