package user

type db interface {
	selectByID(int64) (user, error)
	update(user) (user, error)
	insert(user) (user, error)
}
