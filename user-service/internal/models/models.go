package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID          uuid.UUID  `json:"id"`
	Username    string     `json:"name"`
	Name        string     `json:"username"`
	Phone       *string    `json:"phone,omitempty"`
	DateOfBirth *time.Time `json:"dateOfBirth,omitempty"`
	PhotoURL    string     `json:"photo"`
	CreatedAt   int64
}

type SearchUsersRequest struct {
	Name     *string
	Username *string
	Offset   *int
	Limit    *int
}

type SearchUsersResponse struct {
	Users  []User `json:"users"`
	Offset int
	Count  int
}

type GetUserByUsernameRequest struct {
	Username string `json:"username"`
}

type UserResponse struct {
	ID          uuid.UUID  `json:"id"`
	Username    string     `json:"username"`
	Name        string     `json:"name"`
	Phone       *string    `json:"phone"`
	DateOfBirth *time.Time `json:"dateOfBirth"`
	PhotoURL    string     `json:"photo"`
}
