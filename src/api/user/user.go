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
