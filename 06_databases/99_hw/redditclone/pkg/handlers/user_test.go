package handlers

import (
	"bytes"
	"fmt"
	"redditclone/pkg/session"
	"redditclone/pkg/user"
	"strings"

	"github.com/golang/mock/gomock"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http/httptest"
	"testing"
)

func TestRegister(t *testing.T) {
	// мы передаём t сюда, это надо чтобы получить корректное сообщение если тесты не пройдут
	ctrl := gomock.NewController(t)

	// Finish сравнит последовательсноть вызовов и выведет ошибку если последовательность другая
	defer ctrl.Finish()

	st := user.NewMockUserRepo(ctrl)
	ses := session.NewMockSessionRepo(ctrl)
	service := &UserHandler{
		UserRepo:       st,
		Logger:         zap.NewNop().Sugar(), // не пишет логи
		SessionManager: ses,
	}

	resultUser := []*user.User{
		{ID: 1, Login: "ayta", Password: "12345678"},
	}

	// тут мы записываем последовтаельность вызовов и результат
	st.EXPECT().FindUser(resultUser[0].Login).Return(nil, user.ErrNoUser)
	st.EXPECT().NewUserID().Return(uint32(2))
	st.EXPECT().Add(&user.User{ID: uint32(2), Login: "ayta", Password: "12345678"}).Return(nil)
	st.EXPECT().FindUser(resultUser[0].Login).Return(&user.User{ID: uint32(2), Login: "ayta", Password: "12345678"}, nil)

	req := httptest.NewRequest("POST", "/api/register", strings.NewReader(`{"password": "12345678","username": "ayta"}`))
	w := httptest.NewRecorder()
	ses.EXPECT().Create(w, uint32(2), "/api/register").Return(session.NewSession(uint32(2)), nil)
	service.Register(w, req)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	title := "s"
	if !bytes.Contains(body, []byte(title)) {
		t.Errorf("no text found")
		return
	}

}

func TestBadRegister(t *testing.T) {
	// мы передаём t сюда, это надо чтобы получить корректное сообщение если тесты не пройдут
	ctrl := gomock.NewController(t)

	// Finish сравнит последовательсноть вызовов и выведет ошибку если последовательность другая
	defer ctrl.Finish()

	st := user.NewMockUserRepo(ctrl)
	st.EXPECT().Add(&user.User{ID: uint32(1), Login: "ayta", Password: "12345678"}).Return(nil)
	st.Add(&user.User{ID: 1, Login: "ayta", Password: "12345678"})
	ses := session.NewMockSessionRepo(ctrl)
	service := &UserHandler{
		UserRepo:       st,
		Logger:         zap.NewNop().Sugar(), // не пишет логи
		SessionManager: ses,
	}

	resultUser := []*user.User{
		{ID: 1, Login: "ayta", Password: "12345678"},
	}

	st.EXPECT().FindUser(resultUser[0].Login).Return(&user.User{ID: uint32(2), Login: "ayta", Password: "12345678"}, nil)

	req := httptest.NewRequest("POST", "/api/register", strings.NewReader(`{"password": "12345678","username": "ayta"}`))
	w := httptest.NewRecorder()
	service.Register(w, req)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	title := "s"
	if !bytes.Contains(body, []byte(title)) {
		t.Errorf("no text found")
		return
	}

	req2 := httptest.NewRequest("POST", "/api/register", strings.NewReader(`{password: "12345678","username": "ayta"}`))
	w2 := httptest.NewRecorder()
	service.Register(w2, req2)

	req = httptest.NewRequest("POST", "/api/register", strings.NewReader(`{"password": "12345678","username": "ayta"}`))
	w = httptest.NewRecorder()
	st.EXPECT().FindUser("ayta").Return(&user.User{ID: uint32(2), Login: "ayta", Password: "12345678"}, user.ErrNoUser)

	st.EXPECT().NewUserID().Return(uint32(2))
	st.EXPECT().Add(&user.User{ID: uint32(2), Login: "ayta", Password: "12345678"}).Return(nil)
	st.EXPECT().FindUser(resultUser[0].Login).Return(&user.User{ID: uint32(2), Login: "ayta", Password: "12345678"}, nil)
	ses.EXPECT().Create(w, uint32(2), "/api/register").Return(nil, fmt.Errorf("no user"))
	service.Register(w, req)
}

func TestLogin(t *testing.T) {
	// мы передаём t сюда, это надо чтобы получить корректное сообщение если тесты не пройдут
	ctrl := gomock.NewController(t)

	// Finish сравнит последовательсноть вызовов и выведет ошибку если последовательность другая
	defer ctrl.Finish()

	st := user.NewMockUserRepo(ctrl)
	st.EXPECT().Add(&user.User{ID: uint32(1), Login: "ayta", Password: "12345678"}).Return(nil)
	st.Add(&user.User{ID: 1, Login: "ayta", Password: "12345678"})
	ses := session.NewMockSessionRepo(ctrl)
	service := &UserHandler{
		UserRepo:       st,
		Logger:         zap.NewNop().Sugar(), // не пишет логи
		SessionManager: ses,
	}

	resultUser := &user.User{ID: 1, Login: "ayta", Password: "12345678"}

	// тут мы записываем последовтаельность вызовов и результат
	st.EXPECT().Authorize("ayta", "12345678").Return(resultUser, nil)

	req := httptest.NewRequest("POST", "/api/login", strings.NewReader(`{"password": "12345678","username": "ayta"}`))
	w := httptest.NewRecorder()
	ses.EXPECT().Create(w, uint32(1), "/api/login").Return(session.NewSession(uint32(1)), nil)
	service.Login(w, req)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	title := "s"
	if !bytes.Contains(body, []byte(title)) {
		t.Errorf("no text found")
		return
	}

}

func TestLoginNoUser(t *testing.T) {
	// мы передаём t сюда, это надо чтобы получить корректное сообщение если тесты не пройдут
	ctrl := gomock.NewController(t)

	// Finish сравнит последовательсноть вызовов и выведет ошибку если последовательность другая
	defer ctrl.Finish()

	st := user.NewMockUserRepo(ctrl)
	ses := session.NewMockSessionRepo(ctrl)
	service := &UserHandler{
		UserRepo:       st,
		Logger:         zap.NewNop().Sugar(), // не пишет логи
		SessionManager: ses,
	}

	//resultUser := []*user.User{
	//	{ID: 1, Login: "ayta", Password: "12345678"},
	//}

	// тут мы записываем последовтаельность вызовов и результат
	st.EXPECT().Authorize("ayta", "12345678").Return(nil, user.ErrNoUser)

	req := httptest.NewRequest("POST", "/api/login", strings.NewReader(`{"password": "12345678","username": "ayta"}`))
	w := httptest.NewRecorder()
	//ses.EXPECT().Create(w, uint32(2), "/api/login").Return(session.NewSession(uint32(2)), nil)
	service.Login(w, req)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	title := "s"
	if !bytes.Contains(body, []byte(title)) {
		t.Errorf("no text found")
		return
	}

}

func TestLoginBadPassword(t *testing.T) {
	// мы передаём t сюда, это надо чтобы получить корректное сообщение если тесты не пройдут
	ctrl := gomock.NewController(t)

	// Finish сравнит последовательсноть вызовов и выведет ошибку если последовательность другая
	defer ctrl.Finish()

	st := user.NewMockUserRepo(ctrl)
	st.EXPECT().Add(&user.User{ID: uint32(1), Login: "ayta", Password: "12345678"}).Return(nil)
	st.Add(&user.User{ID: 1, Login: "ayta", Password: "12345678"})

	ses := session.NewMockSessionRepo(ctrl)
	service := &UserHandler{
		UserRepo:       st,
		Logger:         zap.NewNop().Sugar(), // не пишет логи
		SessionManager: ses,
	}

	// тут мы записываем последовтаельность вызовов и результат
	st.EXPECT().Authorize("ayta", "1234578").Return(nil, user.ErrBadPass)

	req := httptest.NewRequest("POST", "/api/login", strings.NewReader(`{"password": "1234578","username": "ayta"}`))
	w := httptest.NewRecorder()
	service.Login(w, req)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	title := "s"
	if !bytes.Contains(body, []byte(title)) {
		t.Errorf("no text found")
		return
	}
	req = httptest.NewRequest("POST", "/api/login", strings.NewReader(`{password: "1234578","username": "ayta"}`))
	w = httptest.NewRecorder()
	service.Login(w, req)

	st.EXPECT().Authorize("ayta", "12345678").Return(&user.User{ID: 2, Login: "ayta", Password: "12345678"}, nil)

	req = httptest.NewRequest("POST", "/api/login", strings.NewReader(`{"password": "12345678","username": "ayta"}`))
	w = httptest.NewRecorder()
	ses.EXPECT().Create(w, uint32(2), "/api/login").Return(nil, fmt.Errorf("no user"))
	service.Login(w, req)
}
