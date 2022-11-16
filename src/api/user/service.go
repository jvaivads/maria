package user

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
	return us.userRepository.SelectByID(userID)
}
