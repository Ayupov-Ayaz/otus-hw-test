package entity

type User struct {
	ID       int64
	Username string `required:"true"`
}

func NewUser(id int64, userName string) User {
	return User{
		ID:       id,
		Username: userName,
	}
}
