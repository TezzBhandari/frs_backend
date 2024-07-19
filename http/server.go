package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/TezzBhandari/frs"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

type Server struct {
	server *http.Server
	router *mux.Router
	Addr   string

	userService frs.UserService
}

func NewHttpServer() *Server {
	s := &Server{
		server: &http.Server{
			IdleTimeout:  2 * time.Second,
			ReadTimeout:  1 * time.Second,
			WriteTimeout: 1 * time.Second,
		},
		router: mux.NewRouter(),
	}

	s.router.Use(reportPanic)

	s.server.Handler = s.router

	s.router.NotFoundHandler = s.handleNotFound()
	router := s.router.PathPrefix("/").Subrouter()

	s.registerUserRoutes(router)

	return s
}

func (s *Server) Open() error {
	if s.Addr == "" {
		return fmt.Errorf("addr required")
	}
	s.server.Addr = s.Addr
	return s.server.ListenAndServe()
}

func (s *Server) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	return s.server.Shutdown(ctx)
}

func (s *Server) handleNotFound() http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusNotFound)
		err := json.NewEncoder(rw).Encode(ErrorMessage{Error: "path not found"})
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

func Error(rw http.ResponseWriter, r *http.Request, err error) {
	log.Error().Err(err).Msg("")

	errCode, errMessage := frs.ErrorCode(err), frs.ErrorMessage(err)

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(ErrorStatusCode(errCode))

	if err := json.NewEncoder(rw).Encode(ErrorMessage{Error: errMessage}); err != nil {
		log.Error().Err(err).Msg("Failed to write response")
	}
}

type ErrorMessage struct {
	Error string `json:"error"`
}

var codes = map[string]int{
	frs.EBADREQUEST: http.StatusBadRequest,
	frs.EINTERNAL:   http.StatusInternalServerError,
}

func ErrorStatusCode(code string) int {
	if v, ok := codes[code]; ok {
		return v
	}
	return codes[frs.EINTERNAL]
}
