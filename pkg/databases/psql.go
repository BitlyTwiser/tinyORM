package databases

import "sync"

type Postgres struct {
	mu sync.Mutex
}

var _ DatabaseHandler = (*Postgres)(nil)

func (pd *Postgres) Create(model any) error {
	pd.mu.Lock()
	defer pd.mu.Unlock()
	return nil
}

func (pd *Postgres) Update(model any) error {
	pd.mu.Lock()
	defer pd.mu.Unlock()
	return nil
}

func (pd *Postgres) Delete(model any) error {
	pd.mu.Lock()
	defer pd.mu.Unlock()
	return nil
}
func (pd *Postgres) Find(model any, id string) error {
	pd.mu.Lock()
	defer pd.mu.Unlock()
	return nil
}

func (pd *Postgres) Where(stmt string, args ...any) error {
	pd.mu.Lock()
	defer pd.mu.Unlock()
	return nil
}
