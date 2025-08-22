package service

import (
	"strings"
	"user-management/internal/errors"
	"user-management/internal/model"
	"user-management/internal/repository"
)

type UserService struct {
	repo repository.PostgresRepository
}

func NewUserService(repo repository.PostgresRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetAllProducts() ([]model.User, error) {
	return s.repo.GetAll()
}

func (s *UserService) GetProduct(id int) (*model.User, error) {
	return s.repo.GetByID(id)
}
func (s *UserService) CreateUser(user model.User) error {
	if err := s.validateUser(user); err != nil {
		return err
	}
	if err:=s.repo.ExistsByID(user.ID);err!=nil{
		return err
	}
	return s.repo.Create(user)
}
func (s *UserService) UpdateUser(id int, user model.User) error {
	if id < 0 {
		return errors.NewValidationError(id, "id cannot be negative")
	}
	return s.repo.Update(id,user)
}
func(s *UserService) DeleteUser(id int)error{
	if id<0{
		return errors.NewValidationError(id,"id cannot be negative")
	}
	return s.repo.Delete(id)
}

func (s *UserService) validateUser(user model.User) error {
	if user.ID < 0 {
		return errors.NewValidationError(user.ID, "id cannot be negative")
	}
	username := strings.TrimSpace(user.Username)
	name := strings.TrimSpace(user.Name)
	email := strings.TrimSpace(user.Email)

	if username == "" {
		return errors.NewValidationError(user.Username, "username can't be null")
	}
	if name == username {
		return errors.NewValidationError(user.Username, "username can't be same as name")
	}
	if email == "" {
		return errors.NewValidationError(user.Email, "mail can't be null")
	}
	return nil
}
