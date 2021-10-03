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

// Search finds all app metadata records selected by the matcher function.
// The function f should NOT modify the app record (in particular, other
// goroutines might be reading entries in the Maintainers slice).
func (s *Store) Search(m Matcher, f func(App)) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, app := range s.apps {
		if m(app) {
			f(app)
		}
	}
}
