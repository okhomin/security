package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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
}

func New(auther auther.Auther, acler acler.Acler, filer filer.Filer, grouper grouper.Grouper) *Server {
	return &Server{
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

	s.router.HandleFunc("/files", s.ListFiles).Methods(http.MethodGet)
	s.router.HandleFunc("/file/{id}", s.UpdateFile).Methods(http.MethodPut)
	s.router.HandleFunc("/file/{id}", s.ReadFile).Methods(http.MethodGet)
	s.router.HandleFunc("/file/{id}", s.DeleteFile).Methods(http.MethodDelete)
	s.router.HandleFunc("/file", s.CreateFile).Methods(http.MethodPost)

	s.router.HandleFunc("/groups", s.ListGroups).Methods(http.MethodGet)
	s.router.HandleFunc("/group/{id}", s.DeleteGroup).Methods(http.MethodDelete)
	s.router.HandleFunc("/group/{id}", s.ReadGroup).Methods(http.MethodGet)
	s.router.HandleFunc("/group/{id}", s.UpdateGroup).Methods(http.MethodPut)
	s.router.HandleFunc("/group", s.CreateGroup).Methods(http.MethodPost)

	s.router.HandleFunc("/acls", s.ListGroups).Methods(http.MethodGet)
	s.router.HandleFunc("/acl/{id}", s.UpdateAcl).Methods(http.MethodPut)
	s.router.HandleFunc("/acl/{id}", s.DeleteAcl).Methods(http.MethodDelete)
	s.router.HandleFunc("/acl/{id}", s.ReadAcl).Methods(http.MethodGet)
	s.router.HandleFunc("/acl", s.CreateAcl).Methods(http.MethodPost)
}

func (s *Server) Run(cfg config.Config) {
	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		WriteTimeout: time.Second * 10,
		ReadTimeout:  time.Second * 15,
		Handler:      s.router,
	}

	go func() {
		log.Println("Launching on address: " + cfg.Port)
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
