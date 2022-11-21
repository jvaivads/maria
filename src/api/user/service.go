package user

import (
	"errors"
	"fmt"
)

var (
	userNotFoundByIDError      = errors.New("user not found by user_id")
	conflictError              = errors.New("conflict internal error")
	userWithSameValueError     = errors.New("common user feature")
	userWithSameValueErrorFunc = func(value string) error {
		return fmt.Errorf("%w: there is already a user with same %s", userWithSameValueError, value)
	}
)

type Service interface {
	getByID(userID int64) (User, error)
	createUser(user NewUserRequest) (User, error)
}

type userService struct {
	userRepository Persister
}

func NewService(userRepository Persister) Service {
	return userService{userRepository: userRepository}
}

func (us userService) getByID(userID int64) (User, error) {
	user, err := us.userRepository.selectByID(userID)
	if err == nil && user.isEmptyUser() {
		return user, userNotFoundByIDError
	}
	return user, err
}

func (us userService) createUser(user NewUserRequest) (User, error) {
	if users, err := us.userRepository.SelectByAny(user.UserName, user.Alias, user.Email); err != nil {
		return User{}, err
	} else if len(users) > 0 {
		switch {
		case users[0].UserName == user.UserName:
			return User{}, userWithSameValueErrorFunc("user_name")
		case users[0].Alias == user.Alias:
			return User{}, userWithSameValueErrorFunc("alias")
		case users[0].Email == user.Email:
			return User{}, userWithSameValueErrorFunc("user_name")
		}
		return User{}, conflictError
	}

	return us.userRepository.CreateUser(user)
}
