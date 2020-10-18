// Package mock implements storage interfaces that
// can be used in unit tests or for running long-season
// with volatile data storage.
package mock

import (
	"context"
	"sync"

	"github.com/hakierspejs/long-season/pkg/models"
	serrors "github.com/hakierspejs/long-season/pkg/storage/errors"
)

// Factory returns mock interfaces specific
// to stored data. Implements storage.Factory interface.
type Factory struct{}

// New returns new mock factory.
func New() *Factory {
	return new(Factory)
}

// Users returns storage interface for manipulating
// users data.
func (Factory) Users(ctx context.Context) *UsersStorage {
	return &UsersStorage{
		data:    make(map[int]models.User),
		counter: 0,
		mutex:   new(sync.Mutex),
	}
}

// UsersStorage implements storage.Users interface
// for mocking purposes.
type UsersStorage struct {
	data    map[int]models.User
	counter int
	mutex   *sync.Mutex
}

// New stores given user data in database and returns
// assigned id.
func (s *UsersStorage) New(ctx context.Context, newUser models.User) (int, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for _, u := range s.data {
		if u.Nickname == newUser.Nickname {
			return 0, serrors.ErrNicknameTaken(u.Nickname)
		}
	}

	newUser.ID = s.counter

	s.data[s.counter] = newUser
	defer func() { s.counter += 1 }()

	return s.counter, nil
}

// Read returns single user data with given ID.
func (s *UsersStorage) Read(ctx context.Context, id int) (*models.User, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	user, ok := s.data[id]
	if !ok {
		return nil, serrors.ErrNoID(id)
	}

	return &user, nil
}

// All returns slice with all users from storage.
func (s *UsersStorage) All(ctx context.Context) ([]models.User, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	res := make([]models.User, len(s.data), len(s.data))
	for _, u := range s.data {
		res = append(res, u)
	}

	return res, nil
}

// Update overwrites existing user data.
func (s *UsersStorage) Update(ctx context.Context, u models.User) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	_, ok := s.data[u.ID]
	if !ok {
		return serrors.ErrNoID(u.ID)
	}

	s.data[u.ID] = u
	return nil
}