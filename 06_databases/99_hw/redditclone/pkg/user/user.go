package user

type User struct {
	ID       uint32
	Login    string
	Password string
}

type UsersRepo interface {
	FindUser(login string) (*User, error)
	Authorize(login, pass string) (*User, error)
	NewUserID() uint32
	Add(u *User) error
}
