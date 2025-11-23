package service

import (
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"

	"github.com/abneribeiro/goapi/internal/model"
	"github.com/abneribeiro/goapi/internal/pkg/jwt"
	"github.com/abneribeiro/goapi/internal/pkg/validator"
	"github.com/abneribeiro/goapi/internal/repository"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrEmailAlreadyExists = errors.New("email already exists")
)

type AuthService struct {
	userRepo   *repository.UserRepository
	jwtManager *jwt.Manager
}

func NewAuthService(userRepo *repository.UserRepository, jwtManager *jwt.Manager) *AuthService {
	return &AuthService{
		userRepo:   userRepo,
		jwtManager: jwtManager,
	}
}

func (s *AuthService) Register(ctx context.Context, req *model.CreateUserRequest) (*model.AuthResponse, error) {
	v := validator.New()
	v.Required("email", req.Email).Email("email", req.Email)
	v.Required("password", req.Password).Password("password", req.Password)
	v.Required("name", req.Name)

	if req.Role != "" {
		v.InList("role", string(req.Role), []string{string(model.RoleOwner), string(model.RoleRenter)})
	}

	if v.Errors().HasErrors() {
		return nil, v.Errors()
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	role := req.Role
	if role == "" {
		role = model.RoleRenter
	}

	user := &model.User{
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		Name:         req.Name,
		Phone:        req.Phone,
		Role:         role,
		Verified:     false,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		if errors.Is(err, repository.ErrEmailExists) {
			return nil, ErrEmailAlreadyExists
		}
		return nil, err
	}

	token, err := s.jwtManager.Generate(user.ID, user.Email, string(user.Role))
	if err != nil {
		return nil, err
	}

	return &model.AuthResponse{
		Token: token,
		User:  *user,
	}, nil
}

func (s *AuthService) Login(ctx context.Context, req *model.LoginRequest) (*model.AuthResponse, error) {
	v := validator.New()
	v.Required("email", req.Email).Email("email", req.Email)
	v.Required("password", req.Password)

	if v.Errors().HasErrors() {
		return nil, v.Errors()
	}

	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	token, err := s.jwtManager.Generate(user.ID, user.Email, string(user.Role))
	if err != nil {
		return nil, err
	}

	return &model.AuthResponse{
		Token: token,
		User:  *user,
	}, nil
}
