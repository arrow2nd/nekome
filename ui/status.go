package ui

import "sync"

// status ステータスメッセージ
type status struct {
	text string
	mu   sync.Mutex
}

func (s *status) set(t string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.text = t
}

func (s *status) get() string {
	return s.text
}

func (s *status) clear() {
	s.set("")
}
