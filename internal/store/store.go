package store

import (
	"errors"
	"github.com/google/uuid"
	"sync"
)

type Store struct {
	mu sync.Mutex
	n  []News
}

func New() *Store {
	return &Store{
		n: []News{},
	}
}

func (s *Store) Create(news News) (News, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	news.ID = uuid.New()
	s.n = append(s.n, news)
	return news, nil
}

func (s *Store) FindAll() ([]News, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.n, nil
}

func (s *Store) FindByID(id uuid.UUID) (News, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, n := range s.n {
		if n.ID == id {
			return n, nil
		}
	}

	return News{}, errors.New("news not found")
}

func (s *Store) DeleteByID(id uuid.UUID) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	idx := func(id uuid.UUID) int {
		for i, n := range s.n {
			if n.ID == id {
				return i
			}
		}
		return -1
	}(id)

	if idx == -1 {
		return errors.New("news not found")
	}

	s.n = append(s.n[:idx], s.n[idx+1:]...)
	return nil
}

func (s *Store) UpdateByID(news News) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for idx, n := range s.n {
		if n.ID == news.ID {
			s.n[idx] = news
			return nil
		}
	}

	return errors.New("news not found")
}
