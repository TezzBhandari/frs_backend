package http

import (
	"net/http"

	"github.com/TezzBhandari/frs"
	"github.com/gorilla/mux"
)

func (s *Server) registerFundRaiserRoutes(r *mux.Router) {
	r.HandleFunc("/fund-raiser", s.handleCreateFundRaiser).Methods(http.MethodGet)
}

func (s *Server) handleCreateFundRaiser(rw http.ResponseWriter, r *http.Request) {
	fundRaiser := &frs.FundRaiser{}
	if err := ReadJsonBody(r.Body, fundRaiser); err != nil {
		Error(rw, r, err)
		return
	}

	err := s.FundRaiserService.CreateFundRaiser(r.Context(), fundRaiser)
	if err != nil {
		Error(rw, r, err)
	}
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusCreated)
}

// func (s *Server) handleFindFundRaisers() {

// }
