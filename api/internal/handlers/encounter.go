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
	"gorm.io/gorm"
)

type EncounterHandler struct {
	encounters   *repository.EncounterRepository
	participants *repository.ParticipantRepository
	tournaments  *repository.TournamentRepository
	users        *repository.UserRepository
}

func NewEncounterHandler(
	encounters *repository.EncounterRepository,
	participants *repository.ParticipantRepository,
	tournaments *repository.TournamentRepository,
	users *repository.UserRepository,
) *EncounterHandler {
	return &EncounterHandler{encounters: encounters, participants: participants, tournaments: tournaments, users: users}
}

func (e *EncounterHandler) GetEncounters(w http.ResponseWriter, r *http.Request) {
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

	userId := r.URL.Query().Get("user_id")
	userUuid, err := uuid.Parse(userId)
	if userId != "" && err != nil {
		writeError(w, http.StatusBadRequest, "invalid user id")
		return
	}

	if !e.tournaments.IsRelated(tournamentId, requesterID) {
		writeError(w, http.StatusForbidden, "forbidden")
		return
	}

	var encounters []models.Encounter
	if userUuid == uuid.Nil {
		encounters, err = e.encounters.List(tournamentId, nil)
	} else {
		encounters, err = e.encounters.List(tournamentId, &userUuid)
	}

	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	writeJSON(w, http.StatusOK, encounters)
}

func (e *EncounterHandler) CreateEncounters(w http.ResponseWriter, r *http.Request) {
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

	if !e.tournaments.IsOwner(tournamentId, requesterID) {
		writeError(w, http.StatusForbidden, "forbidden")
		return
	}

	var encounters []models.CreateEncounterRequest
	if err := decode(r, &encounters); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}

	encounterModels := make([]models.Encounter, len(encounters))
	for i, encounter := range encounters {
		encounterModels[i] = models.Encounter{
			TournamentID: tournamentId,
			PartyA:       encounter.PartyA,
			PartyB:       encounter.PartyB,
			Round:        encounter.Round,
			ScoreA:       nil,
			ScoreB:       nil,
		}
	}

	if err := e.encounters.CreateBatch(encounterModels); err != nil {
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	writeJSON(w, http.StatusCreated, encounterModels)
}

func (e *EncounterHandler) UpdateEncounterScore(w http.ResponseWriter, r *http.Request) {
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

	var body models.UpdateEncounterRequest
	if err := decode(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}

	if body.PartyA != requesterID && body.PartyB != requesterID && !e.tournaments.IsOwner(tournamentId, requesterID) {
		writeError(w, http.StatusForbidden, "forbidden")
		return
	}

	encounter, err := e.encounters.Get(tournamentId, body.PartyA, body.PartyB, body.Round)
	if err != nil && err != gorm.ErrRecordNotFound {
		writeError(w, http.StatusNotFound, "not found")
		return
	} else if err != nil {
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	// owners can update the score without any other conditions.
	if e.tournaments.IsOwner(encounter.TournamentID, requesterID) {
		encounter.ScoreA = body.ScoreA
		encounter.ScoreB = body.ScoreB
		encounter.State = body.State
		encounter.VoterID = &requesterID

		if err := e.encounters.Update(encounter); err != nil {
			writeError(w, http.StatusInternalServerError, "internal server error")
			return
		}

		w.WriteHeader(http.StatusNoContent)
		return
	}

	// participants can either vote for a score, or confirm a voted score.
	if encounter.PartyA != requesterID && encounter.PartyB != requesterID {
		writeError(w, http.StatusForbidden, "forbidden")
		return
	}

	if encounter.State == models.EncounterStatePlanned {
		encounter.ScoreA = body.ScoreA
		encounter.ScoreB = body.ScoreB
		encounter.VoterID = &requesterID
		encounter.State = models.EncounterStateVoted
	} else if encounter.State == models.EncounterStateVoted {
		if *encounter.ScoreA != *body.ScoreA || *encounter.ScoreB != *body.ScoreB || *encounter.VoterID == requesterID {
			writeError(w, http.StatusForbidden, "forbidden")
			return
		}
		encounter.State = models.EncounterStateFinal
	} else {
		writeError(w, http.StatusBadRequest, "bad request")
		return
	}

	if err := e.encounters.Update(encounter); err != nil {
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (e *EncounterHandler) DeleteEncounterScore(w http.ResponseWriter, r *http.Request) {
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

	if !e.tournaments.IsOwner(tournamentId, requesterID) {
		writeError(w, http.StatusForbidden, "forbidden")
		return
	}

	var body models.DeleteEncounterRequest
	if err := decode(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if body.Delete {
		if err := e.encounters.DeleteEncounter(tournamentId, body.PartyA, body.PartyB, body.Round); err != nil {
			writeError(w, http.StatusInternalServerError, "internal server error")
			return
		}
	} else {
		if err := e.encounters.ClearScore(tournamentId, body.PartyA, body.PartyB, body.Round); err != nil {
			writeError(w, http.StatusInternalServerError, "internal server error")
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}
