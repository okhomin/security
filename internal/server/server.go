package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/okhomin/security/internal/service/acler"
	"github.com/okhomin/security/internal/service/filer"
	"github.com/okhomin/security/internal/service/grouper"

	"github.com/gorilla/mux"
	"github.com/okhomin/security/internal/config"
	"github.com/okhomin/security/internal/service/auther"
)

type Server struct {
	router  *mux.Router
	auther  auther.Auther
	acler   acler.Acler
	filer   filer.Filer
	grouper grouper.Grouper
	config  config.Config
}

func New(config config.Config, auther auther.Auther, acler acler.Acler, filer filer.Filer, grouper grouper.Grouper) *Server {
	return &Server{
		config:  config,
		auther:  auther,
		acler:   acler,
		filer:   filer,
		grouper: grouper,
		router:  mux.NewRouter(),
	}
}

func (s *Server) Setup() {
	s.router.HandleFunc("/login", s.Login).Methods(http.MethodPost)
	s.router.HandleFunc("/signup", s.Signup).Methods(http.MethodPost)

	userSubRouter := s.router.NewRoute().Subrouter()
	userSubRouter.Use(func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenStrings := strings.Split(r.Header.Get("Authorization"), " ")
			if len(tokenStrings) != 2 {
				http.Error(w, "invalid token", http.StatusBadRequest)
				return
			}
			tokenString := tokenStrings[1]
			var claims = struct {
				UserID string `json:"user_id"`
				Login  string `json:"login"`
				jwt.StandardClaims
			}{}

			_, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
				return []byte(s.config.JWTKey), nil
			})
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}
			r.WithContext(context.WithValue(r.Context(), "user_id", claims.UserID))
		})
	})

	userSubRouter.HandleFunc("/files", s.ListFiles).Methods(http.MethodGet)
	userSubRouter.HandleFunc("/file/{id}", s.UpdateFile).Methods(http.MethodPut)
	userSubRouter.HandleFunc("/file/{id}", s.ReadFile).Methods(http.MethodGet)
	userSubRouter.HandleFunc("/file/{id}", s.DeleteFile).Methods(http.MethodDelete)
	userSubRouter.HandleFunc("/file", s.CreateFile).Methods(http.MethodPost)

	rootSubRouter := s.router.NewRoute().Subrouter()
	rootSubRouter.Use(func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenStrings := strings.Split(r.Header.Get("Authorization"), " ")
			if len(tokenStrings) != 2 {
				http.Error(w, "invalid token", http.StatusBadRequest)
				return
			}
			tokenString := tokenStrings[1]
			var claims = struct {
				UserID string `json:"user_id"`
				Login  string `json:"login"`
				jwt.StandardClaims
			}{}

			_, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
				return []byte(s.config.JWTKey), nil
			})
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}
			if claims.Login != s.config.RootLogin {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}
		})
	})

	rootSubRouter.HandleFunc("/groups", s.ListGroups).Methods(http.MethodGet)
	rootSubRouter.HandleFunc("/group/{id}", s.DeleteGroup).Methods(http.MethodDelete)
	rootSubRouter.HandleFunc("/group/{id}", s.ReadGroup).Methods(http.MethodGet)
	rootSubRouter.HandleFunc("/group/{id}", s.UpdateGroup).Methods(http.MethodPut)
	rootSubRouter.HandleFunc("/group", s.CreateGroup).Methods(http.MethodPost)

	rootSubRouter.HandleFunc("/acls", s.ListGroups).Methods(http.MethodGet)
	rootSubRouter.HandleFunc("/acl/{id}", s.UpdateAcl).Methods(http.MethodPut)
	rootSubRouter.HandleFunc("/acl/{id}", s.DeleteAcl).Methods(http.MethodDelete)
	rootSubRouter.HandleFunc("/acl/{id}", s.ReadAcl).Methods(http.MethodGet)
	rootSubRouter.HandleFunc("/acl", s.CreateAcl).Methods(http.MethodPost)
}

func (s *Server) Run() {
	srv := &http.Server{
		Addr:         ":" + s.config.Port,
		WriteTimeout: time.Second * 10,
		ReadTimeout:  time.Second * 15,
		Handler:      s.router,
	}

	go func() {
		log.Println("Launching on address: " + s.config.Port)
		log.Println(srv.ListenAndServe())
	}()
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	<-signalChan

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	log.Println("Shutting down")
	if err := srv.Shutdown(ctx); err != nil {
		log.Println(err)
	}
}
