package user

type service struct {
	userRepository Persister
}

func newService(userRepository Persister) service {
	return service{userRepository: userRepository}
}

func (us service) getByID(userID int64) (user, error) {
	return us.userRepository.selectByID(userID)
}
