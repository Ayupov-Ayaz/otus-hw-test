package memory

import (
	"context"
	"github.com/ayupov-ayaz/otus-wh-test/hw12_13_14_15_calendar/internal/storage/entity"
	"sync"
)

type Storage struct {
	mu    sync.RWMutex
	users map[int64]entity.User
}

func New() *Storage {
	return &Storage{
		users: make(map[int64]entity.User),
		mu:    sync.RWMutex{},
	}
}

func (s *Storage) Create(_ context.Context, user entity.User) (id int64, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	id = int64(len(s.users) + 1)
	user.ID = id
	s.users[id] = user

	return id, nil
}

func (s *Storage) Get(_ context.Context, id int64) (*entity.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, ok := s.users[id]
	if !ok {
		return nil, nil
	}

	return &user, nil
}
