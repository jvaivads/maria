package user

import (
	"errors"
	"fmt"
)

var (
	userNotFoundError          = errors.New("user not found")
	conflictError              = errors.New("conflict internal error")
	userWithSameValueError     = errors.New("common user feature")
	userWithSameValueErrorFunc = func(value string) error {
		return fmt.Errorf("%w: there is already a user with same %s", userWithSameValueError, value)
	}
)

type Service interface {
	getByID(int64) (User, error)
	createUser(NewUserRequest) (User, error)
	modifyUser(ModifyUserRequest, User) (User, error)
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
		return user, userNotFoundError
	}
	return user, err
}

func (us userService) createUser(user NewUserRequest) (User, error) {
	if users, err := us.userRepository.selectByAny(user.UserName, user.Alias, user.Email); err != nil {
		return User{}, err
	} else if len(users) > 0 {
		switch {
		case users[0].UserName == user.UserName:
			return User{}, userWithSameValueErrorFunc("user_name")
		case users[0].Alias == user.Alias:
			return User{}, userWithSameValueErrorFunc("alias")
		case users[0].Email == user.Email:
			return User{}, userWithSameValueErrorFunc("email")
		}
		return User{}, conflictError
	}

	return us.userRepository.createUser(user)
}

func (us userService) modifyUser(request ModifyUserRequest, user User) (User, error) {
	var (
		err error
	)
	if user.ID == 0 {
		users, err := us.userRepository.selectByAny(user.UserName, user.Alias, "")
		if err != nil {
			return User{}, err
		}
		if len(users) == 0 {
			return User{}, userNotFoundError
		}
		if len(users) > 1 {
			return User{}, fmt.Errorf("%w: there is more than one user", conflictError)
		}
		user = users[0]
	} else if user, err = us.getByID(user.ID); err != nil {
		return User{}, err
	}

	if err = us.userRepository.withTransaction(func(tx Transactioner) error {
		if _, err = tx.modifyUser(request, user); err != nil {
			return err
		}

		user, err = tx.selectByID(user.ID)
		return err
	}); err != nil {
		return User{}, err
	}

	return user, nil
}
