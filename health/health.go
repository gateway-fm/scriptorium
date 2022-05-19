package health

import (
	"sync"
)

type Status struct {
	status bool
	lock sync.Mutex
}

func NewStatus() *Status {
	return &Status{status: false}
}

func (s *Status) SetHealth(healthy bool) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.status = healthy
}

func (s *Status) IsHealthy() bool {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.status
}
