package server

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func (s *Server) ListFiles(w http.ResponseWriter, h *http.Request) {
	encoder := json.NewEncoder(w)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	files, err := s.filer.List(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := encoder.Encode(files); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) CreateFile(w http.ResponseWriter, h *http.Request) {
	log.Println(h.Context().Value("user_id"))
}

func (s *Server) UpdateFile(w http.ResponseWriter, h *http.Request) {

}

func (s *Server) ReadFile(w http.ResponseWriter, h *http.Request) {

}

func (s *Server) DeleteFile(w http.ResponseWriter, h *http.Request) {

}
