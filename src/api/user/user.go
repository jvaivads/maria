package user

type User struct {
	ID int64 `json:"user_id"`
}

func (u User) isEmptyUser() bool {
	return u.ID == 0
}
