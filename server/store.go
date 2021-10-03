package server

import (
	"sync"
)

// Store is an in-memory database of app metadata records.
// All operations on the store are thread-safe.
type Store struct {
	mu   sync.RWMutex // Protects the apps slice.
	apps []App        // Assume that entries in the apps slice are immutable.
}

// Insert inserts a new app metadata record.
func (s *Store) Insert(app App) {
	s.mu.Lock()
	s.apps = append(s.apps, app)
	s.mu.Unlock()
}

// Search retrieves all app metadata records selected by the matcher function.
// The results are copied into a new slice that is safe for the caller to read.
// The caller should NOT modify the contents of the slice (in particular,
// other goroutines might be reading entries in the Maintainers slice).
func (s *Store) Search(m Matcher) []App {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var results []App
	for _, app := range s.apps {
		if m(app) {
			results = append(results, app)
		}
	}
	return results
}
