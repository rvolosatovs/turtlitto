package api

import (
	"math/rand"
	"sync"
)

type lockedSource struct {
	mu     *sync.Mutex
	source rand.Source
}

func (s *lockedSource) Int63() int64 {
	s.mu.Lock()
	v := s.source.Int63()
	s.mu.Unlock()
	return v
}
func (s *lockedSource) Seed(seed int64) {
	s.mu.Lock()
	s.source.Seed(seed)
	s.mu.Unlock()
}

func newLockedSource(s rand.Source) rand.Source {
	return &lockedSource{
		mu:     &sync.Mutex{},
		source: s,
	}
}
