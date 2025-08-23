package service

import (
	"strings"
	"user-management/internal/errors"
	"user-management/internal/model"
	"user-management/internal/repository"
)

type UserService struct {
	repo repository.UserRepo
}

func NewUserService(repo repository.UserRepo) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetAllUsers() ([]model.User, error) {
	return s.repo.GetAll()
}

func (s *UserService) GetUser(id int) (*model.User, error) {
	return s.repo.GetByID(id)
}
func (s *UserService) CreateUser(user *model.User) error {
	if err := s.validateUser(*user); err != nil {
		return err
	}
	if s.repo.ExistsByEmail(user.Email){
		return errors.NewDuplicateError(user.ID,"already existed with the email")
	}
	return s.repo.Create(user)
}
func (s *UserService) UpdateUser(id int, user model.User) error {
	if id < 0 {
		return errors.NewValidationError(id, "id cannot be negative")
	}
	 existingUser, err := s.repo.GetByID(id)
    if err != nil {
        return err
    }
    
    // Check if username is changing and if new username already exists
    if user.Username != existingUser.Username {
        if s.repo.ExistsByUsername(user.Username) {
            return errors.NewDuplicateError("username", user.Username)
        }
    }
    
    // Check if email is changing and if new email already exists  
    if user.Email != existingUser.Email {
        if s.repo.ExistsByEmail(user.Email) {
            return errors.NewDuplicateError("email", user.Email)
        }
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
