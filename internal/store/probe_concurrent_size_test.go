package store

import (
	"sync"
	"testing"
)

// TestProbeConcurrentStoreAccessSafe asserts that concurrent writers (Put /
// Delete) and readers (Size) are safe under the race detector. A missing lock
// on Size trips `go test -race`.
func TestProbeConcurrentStoreAccessSafe(t *testing.T) {
	s := New()
	const writers = 8
	const iterations = 200
	var wg sync.WaitGroup
	for w := 0; w < writers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < iterations; i++ {
				_ = s.Put("blob", []byte("v"))
				_ = s.Size()
				_ = s.Delete("blob")
				_ = s.Size()
			}
		}()
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < writers*iterations; i++ {
			_ = s.Size()
		}
	}()
	wg.Wait()
}
