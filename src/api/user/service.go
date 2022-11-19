package user

import "errors"

var (
	userNotFoundByIDError = errors.New("user not found by user_id")
)

type Service interface {
	getByID(userID int64) (User, error)
}

type userService struct {
	userRepository Persister
}

func NewService(userRepository Persister) Service {
	return userService{userRepository: userRepository}
}

func (us userService) getByID(userID int64) (User, error) {
	user, err := us.userRepository.SelectByID(userID)
	if err == nil && user.isEmptyUser() {
		return user, userNotFoundByIDError
	}
	return user, err
}
