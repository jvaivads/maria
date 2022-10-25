package service

import (
	"maria/src/api/db"
	"maria/src/api/domain"
)

type User interface {
	GetByID(int64) (domain.User, error)
}

type UserService struct {
	userRepository db.UserDB
}

func NewUserService(userRepository db.UserDB) UserService {
	return UserService{userRepository: userRepository}
}

func (us UserService) GetByID(userID int64) (domain.User, error) {
	return us.userRepository.SelectByID(userID)
}
