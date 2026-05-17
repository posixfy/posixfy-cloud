package models

import "time"

type User struct {
	ID           int64     `json:"id"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"`
	UID          int       `json:"uid"`
	GID          int       `json:"gid"`
	Groups       string    `json:"groups"` // JSON array string e.g. "[100,101]"
	Role         string    `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type CreateUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	UID      int    `json:"uid" binding:"required"`
	GID      int    `json:"gid" binding:"required"`
	Groups   string `json:"groups"`
	Role     string `json:"role"`
}

type UpdateUserRequest struct {
	Password *string `json:"password"`
	GID      *int    `json:"gid"`
	Groups   *string `json:"groups"`
	Role     *string `json:"role"`
}
