package idempotency

import "sync"

type counterMu struct {
	mu     sync.Mutex
	locked int
}

type Locker struct {
	keyMu map[string]*counterMu
	mu    sync.Mutex
}

func NewLocker() *Locker {
	return &Locker{
		keyMu: make(map[string]*counterMu),
		mu:    sync.Mutex{},
	}
}

func (l *Locker) Lock(key string) {
	l.mu.Lock()
	keyed, ok := l.keyMu[key]
	if !ok {
		keyed = new(counterMu)
		l.keyMu[key] = keyed
	}
	keyed.locked++
	l.mu.Unlock()

	keyed.mu.Lock()
}

func (l *Locker) Unlock(key string) {
	l.mu.Lock()
	keyed, ok := l.keyMu[key]
	if !ok {
		// Unknown key
		l.mu.Unlock()
		return
	}
	l.mu.Unlock()

	keyed.mu.Unlock()

	l.mu.Lock()
	keyed.locked--
	if keyed.locked <= 0 {
		delete(l.keyMu, key)
	}
	l.mu.Unlock()
}
