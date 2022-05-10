package user

import (
	"fmt"
	"reflect"
	"testing"

	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

// go test -coverprofile=cover.out && go tool cover -html=cover.out -o cover.html

func TestAdd(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	repo := &UsersMemoryRepository{
		data: db,
	}

	login := "12"
	password := "123456789"
	testItem := &User{
		Login:    login,
		Password: password,
		ID:       1,
	}

	//ok query
	mock.
		ExpectExec(`INSERT INTO users`).
		WithArgs(login, password, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Add(testItem)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	// query error
	mock.
		ExpectExec(`INSERT INTO users`).
		WithArgs(1, password, 1).
		WillReturnError(fmt.Errorf("bad query"))

	err = repo.Add(testItem)
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}

}

type LoginCase struct {
	login    string
	password string
}

func TestAuthorize(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	// good query
	rows := sqlmock.NewRows([]string{"id", "login", "password"})
	expect := []*User{
		{1, "ayta", "123456789"},
		{2, "qwe", "123456789"},
	}
	for _, user := range expect {
		rows = rows.AddRow(user.ID, user.Login, user.Password)
	}

	//testCase := []LoginCase{
	//	{"ayta", "123456789"},
	//	{"ayta", "123456789"},
	//	//{"ayt", "123456789"},
	//}
	mock.
		ExpectQuery("SELECT id, login, password FROM users WHERE").
		WithArgs("ayta").
		WillReturnRows(rows)
	repo := &UsersMemoryRepository{
		data:   db,
		LastID: 1,
	}

	// good request
	item, err1 := repo.Authorize("ayta", "123456789")
	if err1 != nil {
		fmt.Println("error: ", err1)
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if !reflect.DeepEqual(item, expect[0]) {
		t.Errorf("results not match, want %v, have %v", expect[0], item)
		return
	}

	// invalid password
	mock.
		ExpectQuery("SELECT id, login, password FROM users WHERE").
		WithArgs("ayta").
		WillReturnRows(rows)

	item, err1 = repo.Authorize("ayta", "12345679")
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	//no user found
	mock.
		ExpectQuery("SELECT id, login, password FROM users WHERE").
		WithArgs("aya").
		WillReturnError(fmt.Errorf(" No user found"))

	item, err1 = repo.Authorize("aya", "123456789")
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
}

func TestFind(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	// good query
	rows := sqlmock.NewRows([]string{"id", "login", "password"})
	expect := []*User{
		{1, "ayta", "123456789"},
		{2, "aya", "123456789"},
	}
	for _, user := range expect {
		rows = rows.AddRow(user.ID, user.Login, user.Password)
	}

	mock.
		ExpectQuery("SELECT id, login, password FROM users WHERE").
		WithArgs("ayta").
		WillReturnRows(rows)

	repo := &UsersMemoryRepository{
		data:   db,
		LastID: 1,
	}

	// good request
	item, err := repo.FindUser(expect[0].Login)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if !reflect.DeepEqual(item, expect[0]) {
		t.Errorf("results not match, want %v, have %v", item, expect[0])
		return
	}

	mock.
		ExpectQuery("SELECT id, login, password FROM users WHERE").
		WithArgs("wqqew").
		WillReturnError(fmt.Errorf(" No user found"))

	_, err = repo.FindUser("wqqew")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
}

func TestMemoryRepo(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	// good query
	rows := sqlmock.NewRows([]string{"id", "login", "password"})
	expect := []*User{
		{1, "ayta", "123456789"},
		{2, "aya", "123456789"},
	}
	for _, user := range expect {
		rows = rows.AddRow(user.ID, user.Login, user.Password)
	}

	repo := &UsersMemoryRepository{
		data:   db,
		LastID: 1,
	}

	// good request
	item := NewMemoryRepo(db)
	if !reflect.DeepEqual(item, repo) {
		t.Errorf("results not match, want %v, have %v", item, repo)
		return
	}

}
