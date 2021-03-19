package server

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/okhomin/security/internal/models/group"

	"github.com/gorilla/mux"
)

func (s *Server) CreateGroup(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	encoder := json.NewEncoder(w)
	newGroup := new(group.Group)
	if err := decoder.Decode(newGroup); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if createdGroup, err := s.grouper.Create(ctx, *newGroup); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		if err := encoder.Encode(createdGroup); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

}

func (s *Server) DeleteGroup(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if deletedGroup, err := s.grouper.Delete(ctx, mux.Vars(r)["id"]); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		if err := encoder.Encode(deletedGroup); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (s *Server) ListGroups(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	groups, err := s.grouper.List(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := encoder.Encode(groups); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) ReadGroup(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if readGroup, err := s.grouper.Read(ctx, mux.Vars(r)["id"]); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		if err := encoder.Encode(readGroup); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (s *Server) UpdateGroup(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	encoder := json.NewEncoder(w)
	newGroup := new(group.Group)
	if err := decoder.Decode(newGroup); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if updatedGroup, err := s.grouper.Update(ctx, *newGroup); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		if err := encoder.Encode(updatedGroup); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
