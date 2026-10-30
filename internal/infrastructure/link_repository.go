package memory

import (
	"sync"

	"shortapi/internal/domain"
	"shortapi/internal/exception"
)

type LinkRepository struct {
	mu    sync.RWMutex
	links map[string]string
}

func NewLinkRepository() *LinkRepository {
	return &LinkRepository{links: make(map[string]string)}
}

func (r *LinkRepository) Save(link domain.Link) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.links[link.Code()]; exists {
		return exception.ErrCodeTaken
	}
	r.links[link.Code()] = link.OriginalURL()
	return nil
}

func (r *LinkRepository) FindByCode(code string) (domain.Link, error) {
	r.mu.RLock()
	originalURL, ok := r.links[code]
	r.mu.RUnlock()

	if !ok {
		return domain.Link{}, exception.ErrNotFound
	}
	return domain.NewLink(code, originalURL)
}
