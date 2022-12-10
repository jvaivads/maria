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

func (u User) isEmptyUser() bool {
	return u.ID == 0
}

type NewUserRequest struct {
	UserName string `json:"user_name" binding:"required"`
	Alias    string `json:"alias" binding:"required"`
	Email    string `json:"email" binding:"required"`
}

func (u NewUserRequest) toUser(userID int64, dateCreated time.Time, active bool) User {
	return User{
		ID:          userID,
		UserName:    u.UserName,
		Alias:       u.Alias,
		Email:       u.Email,
		DateCreated: dateCreated,
		Active:      active,
	}
}

type ModifyUserRequest struct {
	Active string `json:"active" binding:"required"`
}
