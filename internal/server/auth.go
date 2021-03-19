package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/okhomin/security/internal/models/user"

	"github.com/okhomin/security/internal/service/auther"
)

type Credentials struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func generateToken(u user.User, key string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": u.ID,
		"login":   u.Login,
	})

	return token.SignedString([]byte(key))
}

func (s *Server) Login(w http.ResponseWriter, h *http.Request) {
	decoder := json.NewDecoder(h.Body)
	creds := new(Credentials)

	if err := decoder.Decode(creds); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	u, err := s.auther.Login(ctx, creds.Password, creds.Login)
	if err != nil {
		if err == auther.ErrInvalidLoginOrPassword {
			http.Error(w, "invalid login or password", http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	token, err := generateToken(*u, s.config.JWTKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Authorization", "Bearer "+token)
	fmt.Fprint(w, "")
}

func (s *Server) Signup(w http.ResponseWriter, h *http.Request) {
	decoder := json.NewDecoder(h.Body)
	creds := new(Credentials)

	if err := decoder.Decode(creds); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	u, err := s.auther.Signup(ctx, creds.Password, creds.Login)
	if err != nil {
		if err == auther.ErrAlreadyExist {
			http.Error(w, "user with such login already exist", http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	token, err := generateToken(*u, s.config.JWTKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Authorization", "Bearer "+token)
	fmt.Fprint(w, "")
}
