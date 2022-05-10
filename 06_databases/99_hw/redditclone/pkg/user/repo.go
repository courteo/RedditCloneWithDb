package user

import (
	"database/sql"
	"errors"
	"fmt"
)

var (
	ErrNoUser  = errors.New(" No user found")
	ErrBadPass = errors.New(" Invald password")
)

type UsersMemoryRepository struct {
	data   *sql.DB
	LastID uint32
}

func (repo *UsersMemoryRepository) NewUserID() uint32 {
	repo.LastID++
	return repo.LastID
}

func NewMemoryRepo(db *sql.DB) *UsersMemoryRepository {
	return &UsersMemoryRepository{
		data:   db,
		LastID: 1,
	}
}

func (repo *UsersMemoryRepository) FindUser(login string) (*User, error) {
	u := &User{}
	row := repo.data.QueryRow("SELECT id, login, password FROM users WHERE login = ?", login)
	err := row.Scan(&u.ID, &u.Login, &u.Password)
	if err != nil {
		return nil, ErrNoUser
	}
	return u, nil
}

func (repo *UsersMemoryRepository) Authorize(login, pass string) (*User, error) {
	u := &User{}
	row := repo.data.QueryRow("SELECT id, login, password FROM users WHERE login = ?", login)
	err := row.Scan(&u.ID, &u.Login, &u.Password)
	fmt.Println("row ", row)
	if err != nil {
		return nil, ErrNoUser
	}

	if u.Password != pass {
		return nil, ErrBadPass
	}

	return u, nil
}

func (repo *UsersMemoryRepository) Add(u *User) error {

	_, err := repo.data.Exec(
		"INSERT INTO users (`login`, `password`, `ID`) VALUES (?, ?, ?)",
		u.Login,
		u.Password,
		repo.NewUserID(),
	)
	return err
}
