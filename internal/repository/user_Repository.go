package repository

import (
	"database/sql"
	"fmt"
	"log"
	"prothomuse-server/internal/model"
	"strings"
)

type userRepository struct {
	db *sql.DB
}

// UserRepository defines the methods implemented by the user repository
type UserRepository interface {
	CreateTable() error
	CreateUser(user *model.User) error
	GetUserByEmail(email string) (*model.User, error)
	GetUserByAPIKey(apiKey string) (*model.User, error)
	GetUserByID(id int) (*model.User, error)
	UpdateUser(user *model.User) error
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}
func (r *userRepository) CreateTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		username VARCHAR(255) NOT NULL,
		email VARCHAR(255) UNIQUE NOT NULL,
		password VARCHAR(255) NOT NULL,
		api_key VARCHAR(255) UNIQUE NOT NULL,
		is_active BOOLEAN DEFAULT TRUE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	create index if not exists idx_email on users(email);
	create index if not exists idx_api_key on users(api_key);
	`
	_, err := r.db.Exec(query)
	return err
}
func (r *userRepository) CreateUser(user *model.User) error {
	query := `
		INSERT INTO users (username, email, password, api_key, is_active)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`
	if err := r.db.QueryRow(query,
		user.Username,
		user.Email,
		user.Password,
		user.APIKey,
		user.IsActive,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt); err != nil {
		return err
	}
	log.Println("User created with ID:", user.ID)
	return nil
}
func (r *userRepository) GetUserByEmail(email string) (*model.User, error) {
	query := `
		SELECT id, username, email, password, api_key, is_active, created_at, updated_at
		FROM users
		WHERE email = $1
			`
	user := &model.User{}
	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.APIKey,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	) 
	log.Println(user);
	if err != nil {
		log.Println("Error fetching user by email:", err)
		return nil, err
	}
	return user, nil
}
func (r *userRepository) GetUserByAPIKey(apiKey string) (*model.User, error) {
	query := `
	SELECT id, username, email, password, api_key, is_active, created_at, updated_at
	FROM users	
	WHERE api_key = $1
	`
	user := &model.User{}
	err := r.db.QueryRow(query, apiKey).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.APIKey,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		log.Println("Error fetching user by API key:", err)
		return nil, err
	}
	return user, nil
}

// GetUserByID returns a user by their numeric ID
func (r *userRepository) GetUserByID(id int) (*model.User, error) {
	query := `
	SELECT id, username, email, password, api_key, is_active, created_at, updated_at
	FROM users
	WHERE id = $1
	`
	user := &model.User{}
	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.APIKey,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		log.Println("Error fetching user by ID:", err)
		return nil, err
	}
	return user, nil
}

// UpdateUser updates fields on an existing user. Fields with zero values are updated as provided
func (r *userRepository) UpdateUser(user *model.User) error {
	// Build dynamic SET clause to avoid overwriting password when empty
	setParts := []string{}
	args := []interface{}{}
	idx := 1
	if user.Username != "" {
		setParts = append(setParts, fmt.Sprintf("username = $%d", idx))
		args = append(args, user.Username)
		idx++
	}
	if user.Email != "" {
		setParts = append(setParts, fmt.Sprintf("email = $%d", idx))
		args = append(args, user.Email)
		idx++
	}
	if user.Password != "" {
		setParts = append(setParts, fmt.Sprintf("password = $%d", idx))
		args = append(args, user.Password)
		idx++
	}
	if user.APIKey != "" {
		setParts = append(setParts, fmt.Sprintf("api_key = $%d", idx))
		args = append(args, user.APIKey)
		idx++
	}
	// is_active is a bool; we always include it
	setParts = append(setParts, fmt.Sprintf("is_active = $%d", idx))
	args = append(args, user.IsActive)
	idx++

	// updated_at
	setParts = append(setParts, "updated_at = CURRENT_TIMESTAMP")

	if len(setParts) == 0 {
		return nil // nothing to update
	}

	// join parts
	setClause := strings.Join(setParts, ", ")

	// final query
	query := "UPDATE users SET " + setClause + " WHERE id = $" + fmt.Sprintf("%d", idx)
	args = append(args, user.ID)

	// execute
	_, err := r.db.Exec(query, args...)
	if err != nil {
		log.Println("Error updating user:", err)
		return err
	}
	return nil
}
