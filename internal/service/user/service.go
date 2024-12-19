package user

import (
	"errors"
	"time"

	"http/internal/domain"
)

type userRepository interface {
	Insert(user *domain.User) (*domain.User, error)
	Get(userID string) (*domain.User, error)
	Update(user *domain.User) (*domain.User, error)
	GetAll(returnDeleted bool) ([]domain.User, error)
}

//go:generate go run github.com/vektra/mockery/v2 --name=accountService --structname=AccountService --output=mocks/
type accountService interface {
	Create(userID string) error
	GetUserAccounts(userID string) ([]domain.Account, error)
	DeleteUserAccounts(userID string)
}

type Service struct {
	userRepository userRepository
	accountService accountService
}

func NewService(userRepository userRepository, accountService accountService) *Service {
	return &Service{
		userRepository: userRepository,
		accountService: accountService,
	}
}

func (service Service) CreateUser(name string) (*domain.User, error) {
	u, err := domain.NewUser(name)
	if err != nil {
		return nil, errors.Join(failedToCreateUser, err)
	}

	if err = service.accountService.Create(u.ID); err != nil {
		return nil, errors.Join(failedToCreateAccount, err)
	}

	u, err = service.userRepository.Insert(u)
	if err != nil {
		return nil, errors.Join(failedToPersistUser, err)
	}

	return u, nil
}

func (service Service) DeleteUser(userID string) error {
	u, err := service.userRepository.Get(userID)
	if err != nil {
		return errors.Join(failedToGetUser, err)
	}

	service.accountService.DeleteUserAccounts(u.ID)

	now := time.Now()
	u.DeletedAt = &now

	_, err = service.userRepository.Update(u)
	if err != nil {
		return errors.Join(failedToUpdateUser, err)
	}

	return nil
}

func (service Service) GetUsers(returnDeleted bool) ([]domain.User, error) {
	return service.userRepository.GetAll(returnDeleted)
}

func (service Service) GetAccounts(userID string) ([]domain.Account, error) {
	u, err := service.userRepository.Get(userID)
	if err != nil {
		return nil, errors.Join(failedToGetUser, err)
	}

	accs, err := service.accountService.GetUserAccounts(u.ID)
	if err != nil {
		return nil, errors.Join(failedToGetUserAccounts, err)
	}

	return accs, nil
}
