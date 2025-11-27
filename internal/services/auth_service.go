package services

import (
	"campus-lost-and-found/internal/dto"
	"campus-lost-and-found/internal/models"
	"campus-lost-and-found/internal/repository"
	"campus-lost-and-found/internal/utils"
	"errors"
	"strings"
)

type AuthService struct {
	UserRepo *repository.UserRepository
}

func NewAuthService(userRepo *repository.UserRepository) *AuthService {
	return &AuthService{UserRepo: userRepo}
}

func (s *AuthService) Register(req dto.RegisterRequest) (*dto.AuthResponse, error) {
	existingUser, _ := s.UserRepo.FindByEmail(req.Email)
	if existingUser != nil {
		return nil, errors.New("email already registered")
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	role := models.RoleUser
	if strings.HasSuffix(req.Email, "@students.uii.ac.id") {
		role = models.RoleStudent
	} else if strings.HasSuffix(req.Email, "@uii.ac.id") {
		role = models.RoleStaff
	} else if req.Role != "" {
		// Optional: Allow manual override if needed, or remove this if strict
		// Keeping it for Admin/Security creation if exposed, but usually public reg ignores it.
		// Let's prioritize email logic for these domains, but allow others?
		// For safety, let's say if email matches, enforce it.
		// If not, default to USER (or whatever req.Role is if we trust it? No, public reg shouldn't trust req.Role)
		// Let's stick to the plan:
		// Others -> USER.
		// If we want to allow ADMIN creation, that should be a separate seeded process or admin-only endpoint.
		// So for public register:
		role = models.RoleUser
	}

	user := &models.User{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: hashedPassword,
		Phone:        req.Phone,
		Role:         role,
	}

	if err := s.UserRepo.Create(user); err != nil {
		return nil, err
	}

	token, err := utils.GenerateToken(user.ID, string(user.Role))
	if err != nil {
		return nil, err
	}

	refreshToken, err := utils.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}

	return &dto.AuthResponse{
		Token:        token,
		RefreshToken: refreshToken,
		User: dto.UserResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			Role:  string(user.Role),
		},
	}, nil
}

func (s *AuthService) Login(req dto.LoginRequest) (*dto.AuthResponse, error) {
	user, err := s.UserRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if !utils.CheckPasswordHash(req.Password, user.PasswordHash) {
		return nil, errors.New("invalid credentials")
	}

	token, err := utils.GenerateToken(user.ID, string(user.Role))
	if err != nil {
		return nil, err
	}

	refreshToken, err := utils.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}

	return &dto.AuthResponse{
		Token:        token,
		RefreshToken: refreshToken,
		User: dto.UserResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			Role:  string(user.Role),
		},
	}, nil
}

func (s *AuthService) RefreshToken(req dto.RefreshTokenRequest) (*dto.AuthResponse, error) {
	claims, err := utils.ValidateToken(req.RefreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	if claims.Role != "REFRESH" {
		return nil, errors.New("invalid token type")
	}

	user, err := s.UserRepo.FindByID(claims.UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	newToken, err := utils.GenerateToken(user.ID, string(user.Role))
	if err != nil {
		return nil, err
	}

	// Optionally rotate refresh token here
	newRefreshToken, err := utils.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}

	return &dto.AuthResponse{
		Token:        newToken,
		RefreshToken: newRefreshToken,
		User: dto.UserResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			Role:  string(user.Role),
		},
	}, nil
}
