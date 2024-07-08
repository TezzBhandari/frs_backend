package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/TezzBhandari/frs"
	"github.com/gorilla/mux"
)

type Server struct {
	server *http.Server
	router *mux.Router
	addr   string

	userService frs.UserService
}

func NewHttpServer() *Server {
	s := &Server{
		server: &http.Server{},
		router: mux.NewRouter(),
	}

	s.router.Use(reportPanic)

	s.server.Handler = s.router

	s.router.NotFoundHandler = s.handleNotFound()
	router := s.router.PathPrefix("/").Subrouter()

	s.registerUserRoutes(router)

	return s

}

type NotFound struct {
	Error string `json:"error"`
}

func (s *Server) handleNotFound() http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusNotFound)
		err := json.NewEncoder(rw).Encode(NotFound{Error: "path not found"})
		if err != nil {
			fmt.Printf("failed to write response: %q", err)
		}
	})
}

// middleware catches panics and reports to external service
func reportPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				rw.WriteHeader(http.StatusInternalServerError)
				frs.ReportPanic(err)
			}
		}()
		next.ServeHTTP(rw, r)
	})
}
