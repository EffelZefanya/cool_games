package domain

import (
	"context"
	"time"
)

type User struct {
	ID             int        `json:"id"`
	Email          string     `json:"email" binding:"required,email"`
	HashedPassword string     `json:"-"`
	Password       string     `json:"password,omitempty" binding:"required,min=6"`
	Role           string     `json:"role" binding:"required,oneof=admin customer publisher"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type UserRepository interface {
    Create(ctx context.Context, user *User) error
    GetByEmail(ctx context.Context, email string) (User, error)
    CreatePublisher(ctx context.Context, userID int, name string) error
}

type AuthUsecase interface {
	Register(ctx context.Context, user *User) (AuthResponse, error)
	Login(ctx context.Context, req LoginRequest) (AuthResponse, error)
}