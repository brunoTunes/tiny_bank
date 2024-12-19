package memory

import (
	"errors"
	"maps"
	"sync"

	"http/internal/domain"
)

type UserRepository struct {
	users      map[string]*domain.User
	usersMutex sync.RWMutex
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		users:      make(map[string]*domain.User),
		usersMutex: sync.RWMutex{},
	}
}

func (repo *UserRepository) Get(userID string) (*domain.User, error) {
	repo.usersMutex.RLock()
	defer repo.usersMutex.RUnlock()

	user := repo.users[userID]
	if user == nil {
		return nil, errors.New("user does not exist")
	}

	return user, nil
}

func (repo *UserRepository) Insert(user *domain.User) (*domain.User, error) {
	repo.usersMutex.Lock()
	defer repo.usersMutex.Unlock()

	if repo.users[user.ID] != nil {
		return nil, errors.New("user already exists")
	}

	repo.users[user.ID] = user

	return repo.users[user.ID], nil
}

func (repo *UserRepository) Update(user *domain.User) (*domain.User, error) {
	repo.usersMutex.Lock()
	defer repo.usersMutex.Unlock()

	if repo.users[user.ID] == nil {
		return nil, errors.New("user does not exist")
	}

	repo.users[user.ID] = user

	return repo.users[user.ID], nil
}

func (repo *UserRepository) GetAll(returnDeleted bool) ([]domain.User, error) {
	repo.usersMutex.Lock()
	defer repo.usersMutex.Unlock()

	var users []domain.User

	for user := range maps.Values(repo.users) {
		if user.DeletedAt != nil && !returnDeleted {
			continue
		}

		users = append(users, *user)
	}

	return users, nil
}
