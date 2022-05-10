package middleware

import (
	"github.com/dgrijalva/jwt-go"
	"go.uber.org/zap"
	"net/http"
	"redditclone/pkg/handlers"
	"strings"
)

func Auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		zapLogger, err1 := zap.NewProduction()
		if err1 != nil {
			return
		}

		defer zapLogger.Sync()
		logger := zapLogger.Sugar()
		tokenString := r.Header.Get("Authorization")
		tokenString = tokenString[strings.Index(tokenString, " ")+1:]
		claims := jwt.MapClaims{}
		_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return handlers.TokenSecret, nil
		})
		if err != nil {
			handlers.JsonError(w, http.StatusBadRequest, "Add: cant token parse", logger)
			return
		}
		next.ServeHTTP(w, r)
	}
}
