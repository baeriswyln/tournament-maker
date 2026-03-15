// Brief   :
// Author  : Nicolas Baeriswyl <nicolas@baeriswyl.dev>
// Created : 13.03.2026

package models

import (
	"time"

	"github.com/google/uuid"
)

type EncounterState string

const (
	EncounterStatePlanned EncounterState = "planned"
	EncounterStateVoted   EncounterState = "voted"
	EncounterStateFinal   EncounterState = "final"
)

type Encounter struct {
	TournamentID uuid.UUID      `gorm:"type:uuid;primaryKey;index"  json:"tournament_id"`
	PartyA       uuid.UUID      `gorm:"type:uuid;primaryKey"        json:"partyA"`
	PartyB       uuid.UUID      `gorm:"type:uuid;primaryKey"        json:"partyB"`
	Round        string         `gorm:"primaryKey"                  json:"round"`
	ScoreA       *int           `                                   json:"scoreA,omitempty"`
	ScoreB       *int           `                                   json:"scoreB,omitempty"`
	State        EncounterState `gorm:"default:planned"             json:"state"`
	VoterID      *uuid.UUID     `gorm:"type:uuid;primaryKey"        json:"-"`
	CreatedAt    time.Time      `                                   json:"created_at"`
	UpdatedAt    time.Time      `                                   json:"updated_at"`

	Voter      *User      `gorm:"foreignKey:VoterID"      json:"-"`
	Tournament Tournament `gorm:"foreignKey:TournamentID" json:"-"`
}

// ─── Request / Response DTOs ─────────────────────────────────────────────────

type EncounterSummary struct {
	Total    int `json:"total"`
	Finished int `json:"finished"`
	Score    int `json:"score"`
}

type CreateEncounterRequest struct {
	PartyA uuid.UUID `json:"partyA"`
	PartyB uuid.UUID `json:"partyB"`
	Round  string    `json:"round"`
}

type UpdateEncounterRequest struct {
	PartyA uuid.UUID      `json:"partyA"`
	PartyB uuid.UUID      `json:"partyB"`
	Round  string         `json:"round"`
	ScoreA *int           `json:"scoreA"`
	ScoreB *int           `json:"scoreB"`
	State  EncounterState `json:"state"`
}

type DeleteEncounterRequest struct {
	PartyA uuid.UUID `json:"partyA"`
	PartyB uuid.UUID `json:"partyB"`
	Round  string    `json:"round"`
}
