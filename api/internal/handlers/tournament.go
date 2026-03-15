// Brief   :
// Author  : Nicolas Baeriswyl <nicolas@baeriswyl.dev>
// Created : 13.03.2026

package handlers

import (
	"errors"
	"net/http"

	"baeriswyl.dev/tournament-maker/internal/middleware"
	"baeriswyl.dev/tournament-maker/internal/models"
	"baeriswyl.dev/tournament-maker/internal/repository"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TournamentHandler struct {
	tournaments *repository.TournamentRepository
}

func NewTournamentHandler(t *repository.TournamentRepository) *TournamentHandler {
	return &TournamentHandler{tournaments: t}
}

func (t *TournamentHandler) GetTournaments(w http.ResponseWriter, r *http.Request) {
	requesterID, ok := middleware.GetIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusForbidden, "forbidden")
		return
	}

	part, err := t.tournaments.ListParticipating(requesterID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	owned, err := t.tournaments.ListOwned(requesterID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	resp := models.TournamentsResponse{
		Owning:        owned,
		Participating: part,
	}

	writeJSON(w, http.StatusOK, resp)
}

func (t *TournamentHandler) GetTournamentSingle(w http.ResponseWriter, r *http.Request) {
	requesterID, ok := middleware.GetIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusForbidden, "forbidden")
		return
	}
	id := chi.URLParam(r, "id")
	tournamentId, err := uuid.Parse(id)
	if id == "" || err != nil {
		writeError(w, http.StatusBadRequest, "invalid tournament id")
		return
	}

	tournament, err := t.tournaments.GetByID(tournamentId)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if t.tournaments.IsOwner(tournamentId, requesterID) {
		writeJSON(w, http.StatusOK, tournament)
		return
	}

	if t.tournaments.IsRelated(tournamentId, requesterID) {
		tournamentShort := models.TournamentResponseShort{
			Name:           tournament.Name,
			StartDate:      tournament.StartDate,
			NbParticipants: tournament.NbParticipants,
			Format:         tournament.Format,
			State:          tournament.State,
		}
		writeJSON(w, http.StatusOK, tournamentShort)
		return
	}

	writeError(w, http.StatusForbidden, "forbidden")
}

func (t *TournamentHandler) CreateTournament(w http.ResponseWriter, r *http.Request) {
	requesterID, ok := middleware.GetIDFromContext(r.Context())
	if !ok && middleware.HasLogin(r.Context()) {
		writeError(w, http.StatusForbidden, "forbidden")
		return
	}

	var body models.TournamentRequest
	if err := decode(r, &body); err != nil || *body.Name == "" {
		writeError(w, http.StatusBadRequest, "bad request")
		return
	}

	tournament := models.Tournament{
		OwnerID:    requesterID,
		Name:       *body.Name,
		StartDate:  body.StartDate,
		State:      *body.State,
		Visibility: *body.Visibility,
		Format:     *body.Format,
	}

	if err := t.tournaments.Create(&tournament); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
	}

	writeJSON(w, http.StatusCreated, models.CreateTournamentResponse{ID: tournament.ID})
}

func (t *TournamentHandler) UpdateTournament(w http.ResponseWriter, r *http.Request) {
	requesterID, ok := middleware.GetIDFromContext(r.Context())
	if !ok && middleware.HasLogin(r.Context()) {
		writeError(w, http.StatusForbidden, "forbidden")
		return
	}

	var body models.TournamentRequest
	if err := decode(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "bad request")
		return
	}

	id := chi.URLParam(r, "id")
	tournamentId, err := uuid.Parse(id)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid tournament id")
		return
	}

	tournament, err := t.tournaments.GetByID(tournamentId)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		writeError(w, http.StatusNotFound, "not found")
		return
	} else if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	} else if tournament.OwnerID != requesterID {
		writeError(w, http.StatusForbidden, "forbidden")
		return
	}

	if body.Name != nil {
		tournament.Name = *body.Name
	}
	if body.StartDate != nil {
		tournament.StartDate = body.StartDate
	}
	if body.State != nil {
		tournament.State = *body.State
	}
	if body.Visibility != nil {
		tournament.Visibility = *body.Visibility
	}
	if body.Format != nil {
		tournament.Format = *body.Format
	}

	err = t.tournaments.Update(tournament)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, tournament)
}

func (t *TournamentHandler) DeleteTournament(w http.ResponseWriter, r *http.Request) {
	requesterID, ok := middleware.GetIDFromContext(r.Context())
	if !ok && middleware.HasLogin(r.Context()) {
		writeError(w, http.StatusForbidden, "forbidden")
		return
	}

	id := chi.URLParam(r, "id")
	tournamentId, err := uuid.Parse(id)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid tournament id")
		return
	}

	if !t.tournaments.IsOwner(tournamentId, requesterID) {
		writeError(w, http.StatusForbidden, "forbidden")
		return
	}

	err = t.tournaments.Delete(tournamentId)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
	}

	w.WriteHeader(http.StatusNoContent)
}

func (t *TournamentHandler) ResetTournament(w http.ResponseWriter, r *http.Request) {
	requesterID, ok := middleware.GetIDFromContext(r.Context())
	if !ok && middleware.HasLogin(r.Context()) {
		writeError(w, http.StatusForbidden, "forbidden")
		return
	}

	id := chi.URLParam(r, "id")
	tournamentId, err := uuid.Parse(id)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid tournament id")
		return
	}

	if !t.tournaments.IsOwner(tournamentId, requesterID) {
		writeError(w, http.StatusForbidden, "forbidden")
		return
	}

	err = t.tournaments.Reset(tournamentId)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
	}

	w.WriteHeader(http.StatusNoContent)
}
