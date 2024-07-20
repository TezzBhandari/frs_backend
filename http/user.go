package http

import (
	"fmt"
	"net/http"

	"github.com/TezzBhandari/frs"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

func (s *Server) registerUserRoutes(r *mux.Router) {
	r.HandleFunc("/users", s.handleCreateUser).Methods(http.MethodPost)

}

func (s *Server) handleCreateUser(rw http.ResponseWriter, r *http.Request) {
	user := &frs.User{}

	if err := user.FromJson(r.Body); err != nil {
		log.Error().Err(err).Msg("")
		Error(rw, r, frs.Errorf(frs.EINVALID, "invalid json body"))
		return
	}

	log.Debug().Msg(fmt.Sprintf("%v", user))

	if err := s.UserService.CreateUser(r.Context(), user); err != nil {
		Error(rw, r, err)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusCreated)
}
