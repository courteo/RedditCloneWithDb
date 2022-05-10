package handlers

import (
	"encoding/json"
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
	"go.uber.org/zap"
	"io"
	"net/http"
	"redditclone/pkg/errorsForProject"
	"redditclone/pkg/forms"
	"redditclone/pkg/session"
	"redditclone/pkg/user"
	"time"
)

var TokenSecret = []byte("my_secret_key")

type UserHandler struct {
	Logger         *zap.SugaredLogger
	UserRepo       user.UsersRepo
	SessionManager session.SessionRepo
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	fd := &forms.LoginForm{}
	err := json.NewDecoder(r.Body).Decode(&fd)
	if err != nil {
		JsonError(w, http.StatusBadRequest, "bad request", h.Logger)
		return
	}
	u, err := h.UserRepo.Authorize(fd.Login, fd.Password)

	if err == user.ErrNoUser {
		resp, errMarshal := json.Marshal(map[string]interface{}{
			"message": "user not found",
		})
		if errMarshal != nil {
			JsonError(w, http.StatusBadRequest, "Login: "+errorsForProject.ErrCantMarshal.Error(), h.Logger)
			return
		}

		w.WriteHeader(http.StatusUnauthorized)
		w.Write(resp)
		return
	}

	if err == user.ErrBadPass {
		resp, errMarshal := json.Marshal(map[string]interface{}{
			"message": "invalid password",
		})
		if errMarshal != nil {
			JsonError(w, http.StatusBadRequest, "JsonError: "+errorsForProject.ErrCantMarshal.Error(), h.Logger)
			return
		}

		w.WriteHeader(http.StatusUnauthorized)
		w.Write(resp)
		return
	}
	_, err = h.SessionManager.Create(w, u.ID, r.URL.Path)
	if err != nil {
		h.Logger.Infof("can't create session")
		JsonError(w, http.StatusBadRequest, "JsonError: "+"can't create session", h.Logger)
		return
	}

	resp := GetToken(w, *fd, fmt.Sprint(u.ID), h.Logger)
	w.Write(resp)
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	fd := &forms.LoginForm{}
	err := json.NewDecoder(r.Body).Decode(&fd)
	if err != nil {
		JsonError(w, http.StatusBadRequest, "bad request", h.Logger)
		return
	}
	_, err = h.UserRepo.FindUser(fd.Login)

	if err == nil {
		resp, errMarshal := json.Marshal(map[string]interface{}{
			"errors": []errorsForProject.RegisterError{{
				Msg:      "already exists",
				Location: "body",
				Value:    fd.Login,
				Param:    "username",
			},
			}})
		if errMarshal != nil {
			JsonError(w, http.StatusBadRequest, "Register: "+errorsForProject.ErrCantMarshal.Error(), h.Logger)
			return
		}

		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write(resp)
		return
	}
	id := h.UserRepo.NewUserID()
	h.UserRepo.Add(&user.User{
		ID:       id,
		Login:    fd.Login,
		Password: fd.Password,
	})
	h.UserRepo.FindUser(fd.Login)

	_, err = h.SessionManager.Create(w, id, r.URL.Path)
	if err != nil {
		h.Logger.Infof("can't create session")
		JsonError(w, http.StatusBadRequest, "JsonError: "+"can't create session", h.Logger)
		return
	}

	resp := GetToken(w, *fd, fmt.Sprint(id), h.Logger)

	w.Write(resp)
}

func JsonError(w io.Writer, status int, msg string, Logger *zap.SugaredLogger) {
	resp, err := json.Marshal(map[string]interface{}{
		"status": status,
		"error":  msg,
	})

	if err != nil {
		w.Write([]byte("bad request"))
		return
	}

	w.Write(resp)
}

func GetToken(w http.ResponseWriter, fd forms.LoginForm, id string, Logger *zap.SugaredLogger) (resp []byte) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": jwt.MapClaims{
			"username": fd.Login,
			"id":       id,
		},
		"iat": time.Now().Local().Unix(),
		"exp": time.Now().Add(24 * time.Hour).Local().Unix(),
	})

	tokenString, err := token.SignedString(TokenSecret)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	//fmt.Println("token ", tokenString)
	resp, err = json.Marshal(map[string]interface{}{
		"token": tokenString,
	})
	if err != nil {
		JsonError(w, http.StatusBadRequest, "Get Token: "+errorsForProject.ErrCantMarshal.Error(), Logger)
		return
	}

	return resp
}
