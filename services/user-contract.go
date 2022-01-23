package services

import "time"

type User struct {
	ID         string    `json:"id,omitempty"`
	Name       string    `json:"name,omitempty"`
	Occupation string    `json:"occupation,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type UserService interface {
	Get(id string) (*User, error)
	Delete(id string) error
	DeleteAll() error
	Update(user *User) error
	Create(user *User) error
}
