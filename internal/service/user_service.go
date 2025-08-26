package service

import (
	"fmt"
	"log/slog"
	"strings"
	"user-management/internal/auth"
	"user-management/internal/config"
	"user-management/internal/errors"
	"user-management/internal/model"
	"user-management/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	cfg  *config.Config
	repo repository.UserRepo
}

func NewUserService(cfg *config.Config, repo repository.UserRepo) *UserService {
	return &UserService{
		cfg:  cfg,
		repo: repo}
}

func (s *UserService) GetAllUsers() ([]model.User, error) {
	return s.repo.GetAll()
}

func (s *UserService) GetUser(id int) (*model.User, error) {
	return s.repo.GetByID(id)
}

func (s *UserService) CreateUser(user *model.User) error {
	if err := s.validateUser(*user); err != nil {
		slog.Warn("user validation failed", "error", err, "email", user.Email)
		return err
	}
	if s.repo.ExistsByEmail(user.Email) {
		slog.Warn("Email already exists", "email", user.Email)
		return errors.NewDuplicateError(user.ID, "already existed with the email")
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		slog.Error("password hashing failed", "error", err, "email", user.Email)
		return fmt.Errorf("error while encrypting password %w", err)
	}
	user.Password = string(hashed)
	return s.repo.Create(user)
}

func (s *UserService) Login(email, password string) (string, error) {
	u, err := s.CheckPassword(email, password)
	if err != nil {
		slog.Error("Invalid password.Try again","error",err)
		return "", err
	}
	slog.Info("Password matched!!","user_email",email)
	token, err := auth.GenerateToken(u, s.cfg)
	if err != nil {
		slog.Error("token generation failed", "error", err)
		return "", err
	}
	return token, nil
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
	return s.repo.Update(id, user)
}
func (s *UserService) DeleteUser(id int) error {
	if id < 0 {
		return errors.NewValidationError(id, "id cannot be negative")
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
	if user.Password == "" {
		return errors.NewValidationError(user.Password, "password can't be null")
	}
	return nil
}

func (s *UserService) CheckPassword(email, plainPassword string) (*model.User, error) {
	user, err := s.repo.GetByEmail(email)
	if err != nil {
		slog.Error("password check failed (user lookup)", "error", err, "email", email)
		return nil, err
	}
	return user, bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(plainPassword))
}
