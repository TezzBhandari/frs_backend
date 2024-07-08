package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/TezzBhandari/frs"
	"github.com/gorilla/mux"
)

func (s *Server) registerUserRoutes(r *mux.Router) {
	r.PathPrefix("/users").Subrouter()
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

}

type ErrorMessage struct {
	Error string `json:"error"`
}

func Error(rw http.ResponseWriter, r *http.Request, err error) {
	errCode, errMessage := frs.ErrorCode(err), frs.ErrorMessage(err)

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(ErrorStatusCode(errCode))
	if err := json.NewEncoder(rw).Encode(ErrorMessage{Error: errMessage}); err != nil {
		fmt.Printf("failed to write reponse: %q", err)
	}

}

var codes = map[string]int{
	frs.EBADREQUEST: http.StatusBadRequest,
}

func ErrorStatusCode(code string) int {
	if v, ok := codes[code]; ok {
		return v
	}
	return http.StatusInternalServerError
}
