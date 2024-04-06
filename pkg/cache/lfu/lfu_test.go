package lfu

import "testing"

func TestLFU(t *testing.T) {
	cache := New[int, int](10)
	_ = cache
}
