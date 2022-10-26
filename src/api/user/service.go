package user

type service struct {
	userRepository DB
}

func newService(userRepository DB) service {
	return service{userRepository: userRepository}
}

func (us service) getByID(userID int64) (user, error) {
	return us.userRepository.selectByID(userID)
}
