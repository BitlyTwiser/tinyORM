package databases

import "sync"

type Mysql struct {
	mu sync.Mutex
}

var _ DatabaseHandler = (*Mysql)(nil)

func (m *Mysql) Create(model any) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	return nil
}

func (m *Mysql) Update(model any) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	return nil
}

func (m *Mysql) Delete(model any) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	return nil
}

func (m *Mysql) Find(model any, id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	return nil
}

func (m *Mysql) Where(stmt string, args ...any) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	return nil
}
