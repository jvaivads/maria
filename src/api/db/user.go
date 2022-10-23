package db

import "maria/src/api/domain"

type UserDB interface {
	SelectByID(int64) (domain.User, error)
	Update(domain.User) (domain.User, error)
	Insert(domain.User) (domain.User, error)
}
