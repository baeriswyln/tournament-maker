// Brief   :
// Author  : Nicolas Baeriswyl <nicolas@baeriswyl.dev>
// Created : 13.03.2026

package handlers

import (
	"net/http"

	"baeriswyl.dev/tournament-maker/internal/middleware"
	"baeriswyl.dev/tournament-maker/internal/models"
	"baeriswyl.dev/tournament-maker/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	users *repository.UserRepository
	auth  *middleware.Authenticator
}

func NewAuthHandler(users *repository.UserRepository, auth *middleware.Authenticator) *AuthHandler {
	return &AuthHandler{users: users, auth: auth}
}

func (a *AuthHandler) Authenticate(w http.ResponseWriter, r *http.Request) {
	var body models.LoginRequest
	if err := decode(r, &body); err != nil {
		writeError(w, r, errRequestInvalidBody)
		return
	}

	user, err := a.users.GetByUsername(body.Username)
	if err == nil && user != nil && user.Password != nil {
		err = bcrypt.CompareHashAndPassword([]byte(*user.Password), []byte(body.Password))
		if err == nil {
			pair, tokenErr := a.auth.IssueTokenPair(user.ID.String(), middleware.RefreshTokenWeek)
			if tokenErr == nil {
				resp := models.LoginResponse{
					RefreshToken: pair.RefreshToken,
					AccessToken:  pair.AccessToken,
					ExpiresAt:    pair.ExpiresAt,
					UserID:       user.ID,
				}
				writeJSON(w, http.StatusOK, resp)
				return
			}
		}
	}

	writeError(w, r, errAuthInvalidCredentials)
}

func (a *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var body models.RefreshRequest
	if err := decode(r, &body); err != nil {
		writeError(w, r, errRequestInvalidBody)
		return
	}

	pair, err := a.auth.Refresh(r.Context(), *body.RefreshToken)
	if err != nil {
		writeError(w, r, errAuthRefreshInvalid)
		return
	}

	writeJSON(w, http.StatusOK, pair)
}