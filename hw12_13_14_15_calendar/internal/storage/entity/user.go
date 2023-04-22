package entity

type User struct {
	ID       int
	Username string `required:"true"`
}
