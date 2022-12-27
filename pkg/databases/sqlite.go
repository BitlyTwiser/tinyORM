package databases

import "sync"

type SQLite struct {
	mu sync.Mutex
}

var _ DatabaseHandler = (*SQLite)(nil)

func (s *SQLite) Create(model any) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return nil
}

func (s *SQLite) Update(model any) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return nil
}

func (s *SQLite) Delete(model any) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return nil
}

func (s *SQLite) Find(model any, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return nil
}

func (s *SQLite) Where(stmt string, args ...any) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return nil
}
