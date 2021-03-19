package server

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/okhomin/security/internal/models/acl"
)

func (s *Server) ListAcls(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	acls, err := s.acler.List(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := encoder.Encode(acls); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) CreateAcl(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	encoder := json.NewEncoder(w)
	newAcl := new(acl.Acl)
	if err := decoder.Decode(newAcl); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if createdAcl, err := s.acler.Create(ctx, *newAcl); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		if err := encoder.Encode(createdAcl); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (s *Server) UpdateAcl(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	encoder := json.NewEncoder(w)
	newAcl := new(acl.Acl)
	if err := decoder.Decode(newAcl); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if updatedAcl, err := s.acler.Update(ctx, *newAcl); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		if err := encoder.Encode(updatedAcl); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (s *Server) ReadAcl(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if readAcl, err := s.acler.Read(ctx, mux.Vars(r)["id"]); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		if err := encoder.Encode(readAcl); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (s *Server) DeleteAcl(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if deletedAcl, err := s.acler.Delete(ctx, mux.Vars(r)["id"]); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		if err := encoder.Encode(deletedAcl); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
