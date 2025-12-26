package service

import (
	"errors"
	"laundry-go/internal/config"
	"laundry-go/internal/models"
	"laundry-go/internal/repository"
	"laundry-go/internal/utils"

	"github.com/google/uuid"
)

type AuthService interface {
	Register(req RegisterRequest) (*RegisterResponse, error)
	Login(req LoginRequest) (*LoginResponse, error)
	GetUserByID(userID string) (*models.User, error)
	UpdateLocation(userID string, lat, lng float64) (*models.User, error)
}

type authService struct {
	userRepo repository.UserRepository
	cfg      *config.Config
}

type RegisterRequest struct {
	Name      string   `json:"name"`
	Email     string   `json:"email"`
	Password  string   `json:"password"`
	Phone     string   `json:"phone"`
	Address   string   `json:"address"`
	Latitude  *float64 `json:"latitude,omitempty"`
	Longitude *float64 `json:"longitude,omitempty"`
	Role      string   `json:"role,omitempty"`
}

type RegisterResponse struct {
	User  *UserResponse `json:"user"`
	Token string        `json:"token"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	User  *UserResponse `json:"user"`
	Token string        `json:"token"`
}

type UserResponse struct {
	ID        string   `json:"id"`
	Name      string   `json:"name"`
	Email     string   `json:"email"`
	Phone     string   `json:"phone"`
	Address   string   `json:"address"`
	Latitude  *float64 `json:"latitude,omitempty"`
	Longitude *float64 `json:"longitude,omitempty"`
	Role      string   `json:"role,omitempty"`
}

func NewAuthService(userRepo repository.UserRepository, cfg *config.Config) AuthService {
	return &authService{
		userRepo: userRepo,
		cfg:      cfg,
	}
}

func (s *authService) Register(req RegisterRequest) (*RegisterResponse, error) {
	// Validation
	if utils.IsEmpty(req.Name) {
		return nil, errors.New("name is required")
	}
	if utils.IsEmpty(req.Email) {
		return nil, errors.New("email is required")
	}
	if !utils.ValidateEmail(req.Email) {
		return nil, errors.New("invalid email format")
	}
	if !utils.ValidatePassword(req.Password) {
		return nil, errors.New("password must be at least 8 characters")
	}
	if utils.IsEmpty(req.Phone) {
		return nil, errors.New("phone is required")
	}
	if utils.IsEmpty(req.Address) {
		return nil, errors.New("address is required")
	}

	// Check if email already exists
	existingUser, _ := s.userRepo.FindByEmail(req.Email)
	if existingUser != nil {
		return nil, errors.New("email already registered")
	}

	// Hash password
	passwordHash, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	// Set default role
	role := req.Role
	if role == "" {
		role = "customer"
	}
	if role != "customer" && role != "laundry_owner" {
		role = "customer"
	}

	// Create user
	user := &models.User{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: passwordHash,
		Phone:        req.Phone,
		Address:      req.Address,
		Latitude:     req.Latitude,
		Longitude:    req.Longitude,
		Role:         role,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, errors.New("failed to create user")
	}

	// Generate token
	token, err := utils.GenerateToken(user.ID.String(), user.Email, user.Role, s.cfg.JWT.Secret, s.cfg.JWT.Expiry)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	return &RegisterResponse{
		User:  toUserResponse(user),
		Token: token,
	}, nil
}

func (s *authService) Login(req LoginRequest) (*LoginResponse, error) {
	// Validation
	if utils.IsEmpty(req.Email) {
		return nil, errors.New("email is required")
	}
	if utils.IsEmpty(req.Password) {
		return nil, errors.New("password is required")
	}

	// Find user
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Check password
	if !utils.CheckPasswordHash(req.Password, user.PasswordHash) {
		return nil, errors.New("invalid email or password")
	}

	// Generate token
	token, err := utils.GenerateToken(user.ID.String(), user.Email, user.Role, s.cfg.JWT.Secret, s.cfg.JWT.Expiry)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	return &LoginResponse{
		User:  toUserResponse(user),
		Token: token,
	}, nil
}

func (s *authService) GetUserByID(userID string) (*models.User, error) {
	uuid, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	user, err := s.userRepo.FindByID(uuid)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return user, nil
}

func (s *authService) UpdateLocation(userID string, lat, lng float64) (*models.User, error) {
	// Validate coordinates
	if !utils.ValidateLatitude(lat) {
		return nil, errors.New("invalid latitude (must be between -90 and 90)")
	}
	if !utils.ValidateLongitude(lng) {
		return nil, errors.New("invalid longitude (must be between -180 and 180)")
	}

	uuid, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	user, err := s.userRepo.FindByID(uuid)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Update location (can be called multiple times)
	user.Latitude = &lat
	user.Longitude = &lng

	if err := s.userRepo.Update(user); err != nil {
		return nil, errors.New("failed to update location")
	}

	return user, nil
}

func toUserResponse(user *models.User) *UserResponse {
	return &UserResponse{
		ID:        user.ID.String(),
		Name:      user.Name,
		Email:     user.Email,
		Phone:     user.Phone,
		Address:   user.Address,
		Latitude:  user.Latitude,
		Longitude: user.Longitude,
		Role:      user.Role,
	}
}

