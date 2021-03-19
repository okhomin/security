package server

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/okhomin/security/internal/models/file"
	"github.com/okhomin/security/internal/service/filer"
)

func (s *Server) ListFiles(w http.ResponseWriter, r *http.Request) {
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

func (s *Server) CreateFile(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	encoder := json.NewEncoder(w)
	newFile := new(file.File)
	if err := decoder.Decode(newFile); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if len(newFile.Groups) == 0 && len(newFile.Acls) == 0 {
		http.Error(w, "invalid permissions", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if createdFile, err := s.filer.Create(ctx, *newFile); err != nil {
		if err == filer.ErrFileAlreadyExist {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		if err := encoder.Encode(createdFile); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (s *Server) UpdateFile(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	encoder := json.NewEncoder(w)
	newFile := new(file.File)
	if err := decoder.Decode(newFile); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if len(newFile.Groups) == 0 && len(newFile.Acls) == 0 {
		http.Error(w, "invalid permissions", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if updatedFile, err := s.filer.Update(ctx, r.Context().Value("user_id").(string), *newFile); err != nil {
		if err == filer.ErrPermissionDenied {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		if err := encoder.Encode(updatedFile); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (s *Server) ReadFile(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if readFile, err := s.filer.Read(ctx, r.Context().Value("user_id").(string), mux.Vars(r)["id"]); err != nil {
		if err == filer.ErrPermissionDenied {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		if err := encoder.Encode(readFile); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (s *Server) DeleteFile(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if deletedFile, err := s.filer.Delete(ctx, r.Context().Value("user_id").(string), mux.Vars(r)["id"]); err != nil {
		if err == filer.ErrPermissionDenied {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		if err := encoder.Encode(deletedFile); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
