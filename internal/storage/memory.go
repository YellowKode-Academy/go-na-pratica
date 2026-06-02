package storage

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/yellowkode-academy/linkvault/internal/link"
)

// InMemoryRepository e uma implementacao em memoria para testes.
type InMemoryRepository struct {
	mu     sync.RWMutex
	links  map[int64]link.Link
	nextID int64
	urls   map[string]struct{}
}

// NewInMemoryRepository cria um repositorio em memoria vazio.
func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		links:  make(map[int64]link.Link),
		urls:   make(map[string]struct{}),
		nextID: 1,
	}
}

func (r *InMemoryRepository) Save(_ context.Context, l link.Link) (link.Link, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.urls[l.URL]; exists {
		return link.Link{}, ErrURLDuplicada
	}

	l.ID = r.nextID
	r.nextID++
	if l.CreatedAt.IsZero() {
		l.CreatedAt = time.Now()
	}
	r.links[l.ID] = l
	r.urls[l.URL] = struct{}{}
	return l, nil
}

func (r *InMemoryRepository) FindByID(_ context.Context, id int64) (link.Link, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	l, ok := r.links[id]
	if !ok {
		return link.Link{}, ErrLinkNotFound
	}
	return l, nil
}

func (r *InMemoryRepository) List(_ context.Context) ([]link.Link, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]link.Link, 0, len(r.links))
	for _, l := range r.links {
		result = append(result, l)
	}
	return result, nil
}

func (r *InMemoryRepository) Search(_ context.Context, query string) ([]link.Link, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	query = strings.ToLower(query)
	result := make([]link.Link, 0)
	for _, l := range r.links {
		if strings.Contains(strings.ToLower(l.URL), query) ||
			strings.Contains(strings.ToLower(l.Title), query) ||
			strings.Contains(strings.ToLower(l.Tags), query) {
			result = append(result, l)
		}
	}
	return result, nil
}

func (r *InMemoryRepository) Delete(_ context.Context, id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	l, ok := r.links[id]
	if !ok {
		return ErrLinkNotFound
	}
	delete(r.urls, l.URL)
	delete(r.links, id)
	return nil
}
