package http

import (
	"encoding/json"
	"net/http"

	"github.com/TezzBhandari/frs"
	"github.com/gorilla/mux"
)

func (s *Server) registerUserRoutes(r *mux.Router) {
	r.PathPrefix("/users").Subrouter()
	// userRouter :=  r.PathPrefix("/users").SubRouter()
	r.HandleFunc("/", s.handleCreateUser).Methods(http.MethodGet)

}

func (s *Server) handleCreateUser(rw http.ResponseWriter, r *http.Request) {
	var user *frs.User

	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		Error(rw, r, frs.Errorf(frs.EINVALID, "invalid json body"))
		return
	}

	if err := s.userService.CreateUser(r.Context(), user); err != nil {
		Error(rw, r, err)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusCreated)
}
