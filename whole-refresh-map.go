package ttlmap

import (
	"sync"
	"time"
)

// WholeRefreshMap is a generic map with time-to-live functionality
type WholeRefreshMap[K comparable, V any] struct {
	mu          sync.RWMutex
	data        map[K]V
	ttl         time.Duration
	populate    func() (map[K]V, error)
	lastRefresh time.Time
}

// NewWholeRefreshMap initializes the TTL map with a populate function and TTL duration
func NewWholeRefreshMap[K comparable, V any](populate func() (map[K]V, error), ttl time.Duration) *WholeRefreshMap[K, V] {
	t := &WholeRefreshMap[K, V]{}
	t.populate = populate
	t.ttl = ttl
	return t
}

// Get retrieves a value from the map
func (t *WholeRefreshMap[K, V]) Get(key K) (V, bool, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	if time.Since(t.lastRefresh) > t.ttl {
		var err error
		t.data, err = t.populate()
		if err != nil {
			return *new(V), false, err
		}
		t.lastRefresh = time.Now()
	}

	val, ok := t.data[key]
	return val, ok, nil
}
