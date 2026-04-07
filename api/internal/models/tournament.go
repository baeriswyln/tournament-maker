// Brief   :
// Author  : Nicolas Baeriswyl <nicolas@baeriswyl.dev>
// Created : 13.03.2026

package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ─── Hooks ────────────────────────────────────────────────────────────────────

func (t *Tournament) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}

type TournamentVisibility struct {
	Scoreboard bool `json:"scoreboard"`
	Scores     bool `json:"scores"`
}

type TournamentFormat struct {
	Name        string `json:"name"`
	NbGroups    string `json:"nb_groups"`
	NbEndRounds int    `json:"nb_endrounds"`
}

type Tournament struct {
	ID             uuid.UUID  `gorm:"type:uuid;primaryKey"              json:"id"`
	OwnerID        uuid.UUID  `gorm:"type:uuid;not null;index"          json:"owner_id"`
	Name           string     `gorm:"not null"                          json:"name"`
	StartDate      *time.Time `                                         json:"start_date,omitempty"`
	NbParticipants int        `                                         json:"nb_participants"`
	State          string     `gorm:"default:draft"                      json:"state"`
	// Visibility and Format are stored as JSON columns
	VisibilityJSON string               `gorm:"column:visibility_json"            json:"-"`
	FormatJSON     string               `gorm:"column:format_json"                json:"-"`
	Visibility     TournamentVisibility `gorm:"-"                                 json:"visibility"`
	Format         TournamentFormat     `gorm:"-"                                 json:"format"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relations
	Owner        User                    `gorm:"foreignKey:OwnerID"               json:"-"`
	Participants []TournamentParticipant `gorm:"foreignKey:TournamentID"          json:"-"`
	Encounters   []Encounter             `gorm:"foreignKey:TournamentID"          json:"-"`
}

type TournamentParticipant struct {
	TournamentID uuid.UUID `gorm:"type:uuid;primaryKey;index"  json:"tournament_id"`
	UserID       uuid.UUID `gorm:"type:uuid;primaryKey;index"  json:"user_id"`

	User       User       `gorm:"foreignKey:UserID"       json:"-"`
	Tournament Tournament `gorm:"foreignKey:TournamentID" json:"-"`
}

// ─── Request / Response DTOs ─────────────────────────────────────────────────

type TournamentRequest struct {
	Name       *string               `json:"name,omitempty"`
	StartDate  *time.Time            `json:"start_date,omitempty"`
	Visibility *TournamentVisibility `json:"visibility,omitempty"`
	Format     *TournamentFormat     `json:"format,omitempty"`
	State      *string               `json:"state,omitempty"`
}

type CreateTournamentResponse struct {
	ID uuid.UUID `json:"id"`
}

type TournamentResponseShort struct {
	Name           string           `json:"name"`
	StartDate      *time.Time       `json:"start_date,omitempty"`
	NbParticipants int              `json:"nb_participants"`
	Format         TournamentFormat `json:"format"`
	State          string           `json:"state"`
}

type AddParticipantRequest struct {
	ID uuid.UUID `json:"id"`
}

type DeleteParticipantRequest struct {
	ID uuid.UUID `json:"id"`
}

type ParticipantResponse struct {
	ID          uuid.UUID        `json:"id"`
	VisibleName *string          `json:"visible_name"`
	Encounters  EncounterSummary `json:"encounters"`
}

type TournamentsResponse struct {
	Owning        []Tournament `json:"owning"`
	Participating []Tournament `json:"participating"`
}

type ScoreboardEntry struct {
	ID          uuid.UUID `json:"id"`
	VisibleName *string   `json:"visible_name"`
	Played      int       `json:"played"`
	Won         int       `json:"won"`
	Draw        int       `json:"draw"`
	Lost        int       `json:"lost"`
	ScoreDiff   int       `json:"score_diff"`
	Points      int       `json:"points"`
}
