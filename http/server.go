package http

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/TezzBhandari/frs"
	"github.com/TezzBhandari/frs/utils"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

type Server struct {
	server *http.Server
	router *mux.Router
	Addr   string

	ln net.Listener

	UserService       frs.UserService
	FundRaiserService frs.FundRaiserService
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
	s.router.Use(trackMetrics)

	s.server.Handler = s.router

	s.router.NotFoundHandler = s.handleNotFound()
	router := s.router.PathPrefix("/api/v1").Subrouter()

	s.registerUserRoutes(router)
	s.registerFundRaiserRoutes(router)

	return s
}

func (s *Server) Open() error {
	var err error
	if s.Addr == "" {
		return fmt.Errorf("addr required")
	}
	s.server.Addr = s.Addr

	if s.ln, err = net.Listen("tcp", s.Addr); err != nil {
		return err
	}

	// Begin serving requests on the listener. We use Serve() instead of
	// ListenAndServe() because it allows us to check for listen errors (such
	// as trying to use an already open port) synchronously.
	go s.server.Serve(s.ln)

	return nil
}

func (s *Server) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer func() {
		defer cancel()
		defer log.Info().Msg("server gracefully shutdown")
	}()
	return s.server.Shutdown(ctx)
}

func (s *Server) handleNotFound() http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusNotFound)
		err := json.NewEncoder(rw).Encode(ErrorResponse{Error: "path not found"})
		if err != nil {
			log.Error().Err(fmt.Errorf("%s %w", utils.FailedResponseMsg(), err)).Msg("")
		}
	})
}

func (s *Server) Url() string {
	return fmt.Sprintf("http://localhost:%d", s.Port())
}

func (s *Server) Port() int {
	if s.ln == nil {
		return 0
	}
	return s.ln.Addr().(*net.TCPAddr).Port
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

func trackMetrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		t := time.Now()
		next.ServeHTTP(rw, r)
		timeTaken := float64(time.Since(t).Seconds())
		log.Info().Float64("request time", timeTaken).Str("method", r.Method).Str("path", r.RequestURI).Str("remote address", r.RemoteAddr).Msg("")
	})
}

func Error(rw http.ResponseWriter, r *http.Request, err error) {
	errCode, errMessage := frs.ErrorCode(err), frs.ErrorMessage(err)

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(ErrorStatusCode(errCode))

	if err := json.NewEncoder(rw).Encode(ErrorResponse{Error: errMessage}); err != nil {
		log.Error().Err(fmt.Errorf("%s %w", utils.FailedResponseMsg(), err)).Msg("")
	}
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type SuccessResponse struct {
	Data map[string]any `json:"data"`
}

var codes = map[string]int{
	frs.EBADREQUEST:   http.StatusBadRequest,
	frs.EINVALID:      http.StatusBadRequest,
	frs.EINTERNAL:     http.StatusInternalServerError,
	frs.ENOTFOUND:     http.StatusNotFound,
	frs.EUNAUTHORIZED: http.StatusUnauthorized,
}

func ErrorStatusCode(code string) int {
	if v, ok := codes[code]; ok {
		return v
	}
	return codes[frs.EINTERNAL]
}

func ReadJsonBody(body io.ReadCloser, v any) error {
	if err := json.NewDecoder(body).Decode(v); err != nil {
		switch err {
		case io.EOF:
			return nil
		default:
			return frs.Errorf(frs.EBADREQUEST, utils.InvalidJsonMsg())

		}
	}
	return nil
}
