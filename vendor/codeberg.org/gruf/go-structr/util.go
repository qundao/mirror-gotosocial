package structr

import "sync"

// mutex wraps a sync.Mutex{} to provide
// an additional helper function for safer
// lock-deferred-unlock handling.
type mutex struct{ sync.Mutex }

func (m *mutex) SafeLock() func() {
	m.Lock()
	var once bool
	return func() {
		if !once {
			once = true
			m.Unlock()
		}
	}
}
