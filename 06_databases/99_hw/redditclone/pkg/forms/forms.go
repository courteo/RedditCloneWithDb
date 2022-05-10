package forms

import (
	"github.com/dgrijalva/jwt-go"
	//"gitlab.com/mailru-go/lectures-2022-1/05_web_app/99_hw/redditclone/pkg/comments"
)

type Claims struct {
	jwt.StandardClaims
	jwt.MapClaims
}

type LoginForm struct {
	Login    string `json:"username"`
	Password string `json:"password"`
}

type UserForm struct {
	ID    string `json:"id"`
	Login string `json:"username"`
}

type VoteForm struct {
	ID   string `json:"user"`
	Vote int    `json:"vote"`
}

type CommentForm struct {
	Description string `json:"comment"`
}

type PostFirstData struct {
	Url        string `json:"url"`
	Category   string `json:"category"`
	Text       string `json:"text"`
	Title      string `json:"title"`
	TypeOfPost string `json:"type"`
}
