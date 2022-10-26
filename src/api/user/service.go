package user

type service struct {
	userRepository db
}

func newService(userRepository db) service {
	return service{userRepository: userRepository}
}

func (us service) getByID(userID int64) (user, error) {
	return us.userRepository.selectByID(userID)
}
