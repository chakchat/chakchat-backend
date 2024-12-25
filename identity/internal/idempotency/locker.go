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
		// TODO:
		// Every new idempotency key triggers mutex allocation
		// And no cleaning provided
		// So, if 4000 rps is in 1 month then 9 GB whould be allocated
		// Kinda memory leak
		// If you solve it then make PR https://github.com/gofiber/fiber/blob/main/middleware/idempotency/locker.go
		keyed = new(counterMu)
		l.keyMu[key] = keyed
	}
	keyed.locked++
	l.mu.Unlock()

	keyed.mu.Lock()

	l.mu.Lock()
	keyed.locked--
	if keyed.locked <= 0 {
		delete(l.keyMu, key)
	}
	l.mu.Unlock()
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
}
