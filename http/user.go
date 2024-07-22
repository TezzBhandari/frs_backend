package http

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/TezzBhandari/frs"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

func (s *Server) registerUserRoutes(r *mux.Router) {
	r.HandleFunc("/users", s.handleCreateUser).Methods(http.MethodPost)
	r.HandleFunc("/users", s.handleFindUsers).Methods(http.MethodGet)
	r.HandleFunc("/users/{id}", s.handleFindUserById).Methods(http.MethodGet)
	r.HandleFunc("/users/{id}", s.handleDeleteUser).Methods(http.MethodDelete)
	r.HandleFunc("/users/{id}", s.handleUpdateUser).Methods(http.MethodPut)
}

func (s *Server) handleCreateUser(rw http.ResponseWriter, r *http.Request) {
	user := &frs.User{}

	if err := user.FromJson(r.Body); err != nil {
		switch err {
		case io.EOF:
			break
		default:
			log.Error().Err(err).Msg("")
			Error(rw, r, frs.Errorf(frs.EINVALID, "invalid json body"))
			return

		}
	}

	if err := s.UserService.CreateUser(r.Context(), user); err != nil {
		Error(rw, r, err)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusCreated)
}

func (s *Server) handleFindUsers(rw http.ResponseWriter, r *http.Request) {
	userFilter := &frs.FilterUser{}

	if err := json.NewDecoder(r.Body).Decode(userFilter); err != nil {
		switch err {
		// client didn't send request body
		case io.EOF:
			break
		// client sent invalid json body
		default:
			Error(rw, r, frs.Errorf(frs.EINVALID, "invalid json body"))
			return
		}
	}

	users, _, err := s.UserService.FindUsers(r.Context(), userFilter)
	if err != nil {
		Error(rw, r, err)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)

	if users == nil {
		users = []*frs.User{}
	}

	err = json.NewEncoder(rw).Encode(SuccessMessage{
		Data: map[string]any{
			"users": users,
		},
	})
	if err != nil {
		log.Error().Err(fmt.Errorf("failed to write response %w", err)).Msg("")
	}

}

func (s *Server) handleFindUserById(rw http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	userId, err := strconv.ParseInt(id, 0, 64)
	if err != nil {
		Error(rw, r, frs.Errorf(frs.EINTERNAL, "invalid user id"))
		return
	}

	user, err := s.UserService.FindUserById(r.Context(), userId)
	if err != nil {
		log.Error().Err(err).Msg("")
		Error(rw, r, err)
		return
	}
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	err = json.NewEncoder(rw).Encode(SuccessMessage{
		Data: map[string]any{
			"user": user,
		},
	})
	if err != nil {
		log.Error().Err(err).Msg("failed to write response")
	}
}

func (s *Server) handleDeleteUser(rw http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	userId, err := strconv.ParseInt(id, 0, 64)
	if err != nil {
		Error(rw, r, frs.Errorf(frs.EINVALID, "invlaid user id"))
		return
	}

	err = s.UserService.DeleteUser(r.Context(), userId)
	if err != nil {
		log.Error().Err(err).Msg("")
		Error(rw, r, err)
		return
	}
	rw.WriteHeader(http.StatusOK)
}

func (s *Server) handleUpdateUser(rw http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	userId, err := strconv.ParseInt(id, 0, 64)
	if err != nil {
		Error(rw, r, frs.Errorf(frs.EINVALID, "invalid user id"))
		return
	}

	updateUser := &frs.UpdateUser{}

	err = json.NewDecoder(r.Body).Decode(updateUser)
	if err != nil {
		switch err {
		case io.EOF:
			break
		default:
			Error(rw, r, frs.Errorf(frs.EINVALID, "invalid json body"))
			return
		}
	}

	user, err := s.UserService.UpdateUser(r.Context(), userId, *updateUser)
	if err != nil {
		Error(rw, r, err)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)

	err = json.NewEncoder(rw).Encode(SuccessMessage{
		Data: map[string]any{
			"user": user,
		},
	})

	if err != nil {
		log.Error().Err(err).Msg("failed to write response")
	}
}
