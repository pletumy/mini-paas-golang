package repository

import (
	"sync"
	"time"
)

type MemoryRepo struct {
	data map[string]Application
	mu   sync.Mutex
}

func NewMemoryRepo() *MemoryRepo {
	return &MemoryRepo{data: make(map[string]Application)}
}

func (r *MemoryRepo) Create(app Application) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if app.CreatedAt.IsZero() {
		app.CreatedAt = time.Now()
	}
	app.UpdatedAt = app.CreatedAt
	r.data[app.ID] = app
	return nil
}

func (r *MemoryRepo) Update(app Application) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.data[app.ID] = app
	return nil
}

func (r *MemoryRepo) UpdateStatus(id, status string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	app := r.data[id]
	app.Status = status
	app.UpdatedAt = time.Now()
	r.data[id] = app
	return nil
}

func (r *MemoryRepo) GetByID(id string) (*Application, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	app, ok := r.data[id]
	if !ok {
		return nil, nil
	}
	return &app, nil
}

func (r *MemoryRepo) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.data, id)
	return nil
}
