package user

import (
	"time"
)

type User struct {
	ID          int64     `json:"user_id"`
	UserName    string    `json:"user_name"`
	Alias       string    `json:"alias"`
	Email       string    `json:"email"`
	DateCreated time.Time `json:"date_created"`
	Active      bool      `json:"active"`
}

type NewUserRequest struct {
	UserName string `json:"user_name" binding:"required"`
	Alias    string `json:"alias" binding:"required"`
	Email    string `json:"email" binding:"required"`
}

func (u User) isEmptyUser() bool {
	return u.ID == 0
}
