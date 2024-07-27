package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/TezzBhandari/frs"
	"github.com/TezzBhandari/frs/utils"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

func (s *Server) registerFundRaiserRoutes(r *mux.Router) {
	r.HandleFunc("/fund-raiser", s.handleCreateFundRaiser).Methods(http.MethodPost)
	r.HandleFunc("/fund-raiser", s.handleFindFundRaisers).Methods(http.MethodGet)
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
		return
	}
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusCreated)
}

func (s *Server) handleFindFundRaisers(rw http.ResponseWriter, r *http.Request) {
	filterFundRaiser := &frs.FilterFundRaiser{}
	if err := ReadJsonBody(r.Body, filterFundRaiser); err != nil {
		Error(rw, r, err)
		return
	}

	fundRaisers, _, err := s.FundRaiserService.FindFundRaiser(r.Context(), filterFundRaiser)
	if err != nil {
		Error(rw, r, err)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)

	if err = json.NewEncoder(rw).Encode(SuccessResponse{
		Data: map[string]any{
			"fundraisers": fundRaisers,
		},
	}); err != nil {
		log.Error().Err(fmt.Errorf("%s %w", utils.FailedResponseMsg(), err)).Msg("")
	}

}
