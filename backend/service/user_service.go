package service

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"

	"posixfy-cloud/backend/auth"
	"posixfy-cloud/backend/models"
)

type UserService struct {
	DB *sql.DB
}

func NewUserService(db *sql.DB) *UserService {
	return &UserService{DB: db}
}

func (s *UserService) Bootstrap(password string) {
	if password == "" {
		return
	}

	var count int
	s.DB.QueryRow("SELECT COUNT(*) FROM users WHERE role = 'admin'").Scan(&count)
	if count > 0 {
		return
	}

	hash, err := auth.HashPassword(password)
	if err != nil {
		log.Fatalf("failed to hash admin password: %v", err)
	}

	_, err = s.DB.Exec(
		"INSERT INTO users (username, password_hash, uid, gid, groups, role) VALUES (?, ?, ?, ?, ?, ?)",
		"admin", hash, 0, 0, "[0]", "admin",
	)
	if err != nil {
		log.Fatalf("failed to bootstrap admin: %v", err)
	}
	log.Println("bootstrapped admin user")
}

func (s *UserService) Authenticate(username, password string) (*models.User, error) {
	user, err := s.GetByUsername(username)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}
	if !auth.CheckPassword(user.PasswordHash, password) {
		return nil, errors.New("invalid credentials")
	}
	return user, nil
}

func (s *UserService) GetByUsername(username string) (*models.User, error) {
	u := &models.User{}
	err := s.DB.QueryRow(
		"SELECT id, username, password_hash, uid, gid, groups, role, created_at FROM users WHERE username = ?",
		username,
	).Scan(&u.ID, &u.Username, &u.PasswordHash, &u.UID, &u.GID, &u.Groups, &u.Role, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (s *UserService) GetByID(id int64) (*models.User, error) {
	u := &models.User{}
	err := s.DB.QueryRow(
		"SELECT id, username, password_hash, uid, gid, groups, role, created_at FROM users WHERE id = ?",
		id,
	).Scan(&u.ID, &u.Username, &u.PasswordHash, &u.UID, &u.GID, &u.Groups, &u.Role, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (s *UserService) List() ([]models.User, error) {
	rows, err := s.DB.Query("SELECT id, username, password_hash, uid, gid, groups, role, created_at FROM users ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Username, &u.PasswordHash, &u.UID, &u.GID, &u.Groups, &u.Role, &u.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func (s *UserService) Create(req models.CreateUserRequest) (*models.User, error) {
	hash, err := auth.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	groups := req.Groups
	if groups == "" {
		groups = "[]"
	}
	// validate groups is valid JSON array
	var arr []int
	if err := json.Unmarshal([]byte(groups), &arr); err != nil {
		return nil, errors.New("groups must be a JSON array of integers")
	}

	role := req.Role
	if role == "" {
		role = "user"
	}

	result, err := s.DB.Exec(
		"INSERT INTO users (username, password_hash, uid, gid, groups, role) VALUES (?, ?, ?, ?, ?, ?)",
		req.Username, hash, req.UID, req.GID, groups, role,
	)
	if err != nil {
		return nil, err
	}

	id, _ := result.LastInsertId()
	return s.GetByID(id)
}

func (s *UserService) Update(id int64, req models.UpdateUserRequest) (*models.User, error) {
	user, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	if req.Password != nil && *req.Password != "" {
		hash, err := auth.HashPassword(*req.Password)
		if err != nil {
			return nil, err
		}
		user.PasswordHash = hash
	}
	if req.GID != nil {
		user.GID = *req.GID
	}
	if req.Groups != nil {
		var arr []int
		if err := json.Unmarshal([]byte(*req.Groups), &arr); err != nil {
			return nil, errors.New("groups must be a JSON array of integers")
		}
		user.Groups = *req.Groups
	}
	if req.Role != nil {
		user.Role = *req.Role
	}

	_, err = s.DB.Exec(
		"UPDATE users SET password_hash = ?, gid = ?, groups = ?, role = ? WHERE id = ?",
		user.PasswordHash, user.GID, user.Groups, user.Role, id,
	)
	if err != nil {
		return nil, err
	}

	return s.GetByID(id)
}

func (s *UserService) Delete(id int64) error {
	result, err := s.DB.Exec("DELETE FROM users WHERE id = ?", id)
	if err != nil {
		return err
	}
	n, _ := result.RowsAffected()
	if n == 0 {
		return errors.New("user not found")
	}
	return nil
}
