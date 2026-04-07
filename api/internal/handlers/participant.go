// Brief   :
// Author  : Nicolas Baeriswyl <nicolas@baeriswyl.dev>
// Created : 13.03.2026

package handlers

import (
	"net/http"

	"baeriswyl.dev/tournament-maker/internal/middleware"
	"baeriswyl.dev/tournament-maker/internal/models"
	"baeriswyl.dev/tournament-maker/internal/repository"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type ParticipantHandler struct {
	participants *repository.ParticipantRepository
	tournaments  *repository.TournamentRepository
	users        *repository.UserRepository
}

func NewParticipantHandler(
	participants *repository.ParticipantRepository,
	tournaments *repository.TournamentRepository,
	users *repository.UserRepository,
) *ParticipantHandler {
	return &ParticipantHandler{participants: participants, tournaments: tournaments, users: users}
}

func (p *ParticipantHandler) GetParticipants(w http.ResponseWriter, r *http.Request) {
	requesterID, ok := middleware.GetIDFromContext(r.Context())
	if !ok {
		writeError(w, r, errAuthForbidden)
		return
	}

	id := chi.URLParam(r, "id")
	tournamentId, err := uuid.Parse(id)
	if err != nil {
		writeError(w, r, errRequestInvalidID)
		return
	}

	if !p.tournaments.IsRelated(tournamentId, requesterID) {
		writeError(w, r, errAuthForbidden)
		return
	}

	parts, err := p.tournaments.ListParticipants(tournamentId)
	if err != nil {
		writeError(w, r, errInternal, err)
		return
	}

	writeJSON(w, http.StatusOK, parts)
}

func (p *ParticipantHandler) AddParticipants(w http.ResponseWriter, r *http.Request) {
	requesterID, ok := middleware.GetIDFromContext(r.Context())
	if !ok {
		writeError(w, r, errAuthForbidden)
		return
	}

	id := chi.URLParam(r, "id")
	tournamentId, err := uuid.Parse(id)
	if err != nil {
		writeError(w, r, errRequestInvalidID)
		return
	}

	var body models.AddParticipantRequest
	if err := decode(r, &body); err != nil {
		writeError(w, r, errRequestInvalidBody)
		return
	}

	tournament, err := p.tournaments.GetByID(tournamentId)
	if err != nil {
		writeError(w, r, errTournamentNotFound)
		return
	}
	if tournament.State != "open" {
		writeError(w, r, errTournamentNotOpen)
		return
	}

	// owners can add users that the owner owns
	// other users can only add themselves
	if (p.tournaments.IsOwner(tournamentId, requesterID) && p.users.IsParent(body.ID, requesterID)) ||
		requesterID == body.ID {
		if err := p.participants.Add(tournamentId, body.ID); err != nil {
			writeError(w, r, errInternal, err)
			return
		}
		writeJSON(w, http.StatusCreated, body)
		return
	}

	writeError(w, r, errAuthForbidden)
}

func (p *ParticipantHandler) DeleteParticipant(w http.ResponseWriter, r *http.Request) {
	requesterID, ok := middleware.GetIDFromContext(r.Context())
	if !ok {
		writeError(w, r, errAuthForbidden)
		return
	}

	id := chi.URLParam(r, "id")
	tournamentId, err := uuid.Parse(id)
	if err != nil {
		writeError(w, r, errRequestInvalidID)
		return
	}

	var body models.DeleteParticipantRequest
	if err := decode(r, &body); err != nil {
		writeError(w, r, errRequestInvalidBody)
		return
	}

	if !p.tournaments.IsOwner(tournamentId, requesterID) && requesterID != body.ID {
		writeError(w, r, errAuthForbidden)
		return
	}

	if err = p.participants.Remove(tournamentId, body.ID); err != nil {
		writeError(w, r, errInternal, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}