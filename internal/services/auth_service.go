package services

import (
	"errors"
	"log"
	"strings"

	"prothomuse-server/internal/model"
	"prothomuse-server/internal/repository"
	"prothomuse-server/internal/utils"
)

type AuthService struct {
	userRepo repository.UserRepository
}

func NewAuthService(userRepo repository.UserRepository) *AuthService {
	return &AuthService{userRepo: userRepo}
}

func (s *AuthService) RegisterUser(req model.RegisterRequest) (*model.User, error) {
	if err := validateRegisterRequest(req); err != nil {
		return nil, err
	}
	existingUser, err := s.userRepo.GetUserByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("user with this email already exists")
	}
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		log.Printf("error in hashing the password: %v", err)
		return nil, err
	}
	apikey, err := utils.GenerateAPIKey()
	if err != nil {
		log.Printf("error in generating api key: %v", err)
		return nil, err
	}
	user := &model.User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
		APIKey:   apikey,
		IsActive: true,
	}
	if err := s.userRepo.CreateUser(user); err != nil {
		log.Println("error in creating the user on the database in sql code side")
		return nil, err
	}
	return user, nil
}

// login user
func (s *AuthService) Login(req model.LoginRequest) (*model.LoginResponse, error) {
	user, err := s.userRepo.GetUserByEmail(req.Email)
	if err != nil {
		log.Println("error in getting the user by email in login ")
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}
	if !user.IsActive {
		return nil, errors.New("user is not active")
	}
	if !utils.CheckPasswordHash(req.Password, user.Password) {
		return nil, errors.New("invalid password")
	}
	token, err := utils.GenerateJWT(user.ID, user.Email)
	if err != nil {
		log.Println("error in generating the jwt token ")
		return nil, err
	}
	return &model.LoginResponse{
		ID:       user.ID,
		Token:    token,
		APIKey:   user.APIKey,
		Username: user.Username,
	}, nil
}

func (s *AuthService) GetUserByAPIKey(apiKey string) (*model.User, error) {
	if apiKey == "" {
		return nil, errors.New("API key is required")
	}
	user, err := s.userRepo.GetUserByAPIKey(apiKey)
	if err != nil {
		log.Println("error in getting the user by api key ")
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}

// UpdateUser updates an existing user's mutable fields.
// The provided model.User must include the ID of the user to update.
func (s *AuthService) UpdateUser(update model.UpdateUserRequest) (*model.User, error) {
	if update.ID == 0 {
		return nil, errors.New("user id is required")
	}

	// fetch existing user
	existingUser, err := s.userRepo.GetUserByID(update.ID)
	if err != nil {
		log.Println("error fetching user for update:", err)
		return nil, err
	}
	if existingUser == nil {
		return nil, errors.New("user not found")
	}

	// if email is changing, ensure uniqueness
	if update.Email != nil && *update.Email != existingUser.Email {
		other, err := s.userRepo.GetUserByEmail(*update.Email)
		if err != nil && err.Error() != "sql: no rows in result set" {
			return nil, err
		}
		if other != nil {
			return nil, errors.New("another user with this email already exists")
		}
		existingUser.Email = *update.Email
	}

	if update.Username != nil {
		existingUser.Username = *update.Username
	}

	if update.Password != nil {
		hashed, err := utils.HashPassword(*update.Password)
		if err != nil {
			log.Println("error hashing updated password:", err)
			return nil, err
		}
		existingUser.Password = hashed
	}

	if update.APIKey != nil {
		existingUser.APIKey = *update.APIKey
	}

	if update.IsActive != nil {
		existingUser.IsActive = *update.IsActive
	}

	if err := s.userRepo.UpdateUser(existingUser); err != nil {
		log.Println("error updating user in repo:", err)
		return nil, err
	}
	return existingUser, nil
}

func validateRegisterRequest(req model.RegisterRequest) error {
	if req.Email == "" || !strings.Contains(req.Email, "@") {
		
		return errors.New("invalid email address")
		
	}
	if req.Password == "" {
		log.Println(req.Email);
		log.Println(req.Password);
		return errors.New("password is required")
	}
	if len(req.Password) < 6 {
		return errors.New("password must be at least 6 characters long")
	}
	if req.Username == "" {
		return errors.New("username is required")
	}
	return nil
}
