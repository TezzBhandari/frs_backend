package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/TezzBhandari/frs"
	"github.com/TezzBhandari/frs/utils"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

func (s *Server) registerFundRaiserRoutes(r *mux.Router) {
	r.HandleFunc("/fund-raiser", s.handleCreateFundRaiser).Methods(http.MethodPost)
	r.HandleFunc("/fund-raiser", s.handleFindFundRaiser).Methods(http.MethodGet)
	r.HandleFunc("/fund-raiser/{id}", s.handleFindFundRaiserById).Methods(http.MethodGet)
	r.HandleFunc("/fund-raiser/{id}", s.handleDeleteFundRaiser).Methods(http.MethodDelete)
	r.HandleFunc("/fund-raiser/{id}", s.handleUpdateFundRaiser).Methods(http.MethodPut)
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

func (s *Server) handleFindFundRaiser(rw http.ResponseWriter, r *http.Request) {
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

func (s *Server) handleFindFundRaiserById(rw http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	fundRaiserId, err := strconv.ParseInt(id, 0, 64)
	if err != nil {
		Error(rw, r, frs.Errorf(frs.EINVALID, utils.InvalidFundRaiserIdMsg()))
		return
	}
	fundRaiser, err := s.FundRaiserService.FindFundRaiserById(r.Context(), fundRaiserId)
	if err != nil {
		Error(rw, r, err)
		return
	}
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(rw).Encode(SuccessResponse{
		Data: map[string]any{
			"fund-raiser": fundRaiser,
		},
	}); err != nil {
		log.Error().Err(fmt.Errorf("%s %w", utils.FailedResponseMsg(), err)).Msg("")
	}
}

func (s *Server) handleDeleteFundRaiser(rw http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	fundRaiserId, err := strconv.ParseInt(id, 0, 64)
	if err != nil {
		Error(rw, r, frs.Errorf(frs.EINVALID, utils.InvalidFundRaiserIdMsg()))
		return
	}
	err = s.FundRaiserService.DeleteFundRaiser(r.Context(), fundRaiserId)
	if err != nil {
		Error(rw, r, err)
		return
	}

	rw.WriteHeader(http.StatusOK)
}

func (s *Server) handleUpdateFundRaiser(rw http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	fundRaiserId, err := strconv.ParseInt(id, 0, 64)
	if err != nil {
		Error(rw, r, frs.Errorf(frs.EINVALID, utils.InvalidFundRaiserIdMsg()))
		return
	}

	updFundRaiser := &frs.UpdateFundRaiser{}
	if err := ReadJsonBody(r.Body, updFundRaiser); err != nil {
		Error(rw, r, err)
		return
	}
	updatedFundRaiser, err := s.FundRaiserService.UpdateFundRaiser(r.Context(), fundRaiserId, updFundRaiser)
	if err != nil {
		Error(rw, r, err)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(rw).Encode(SuccessResponse{
		Data: map[string]any{
			"fund-raiser": updatedFundRaiser,
		},
	}); err != nil {
		log.Error().Err(fmt.Errorf("%s %w", utils.FailedResponseMsg(), err)).Msg("")
	}

}
