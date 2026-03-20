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
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	users *repository.UserRepository
	auth  *middleware.Authenticator
}

func NewUserHandler(users *repository.UserRepository, auth *middleware.Authenticator) *UserHandler {
	return &UserHandler{users: users, auth: auth}
}

func (u *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	id, ok := middleware.GetIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusForbidden, "forbidden")
		return
	}

	users, err := u.users.ListByOwnerTournaments(id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, users)
}

// RegisterUser can be done by any non-registered user. Either providing username-password to get a login,
// or not in which case the user has limited permissions, and additionally gets a long-living refresh token.
func (u *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var userRequest models.CreateUserRequest
	if err := decode(r, &userRequest); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	if (userRequest.Password == nil) != (userRequest.Username == nil) {
		writeError(w, http.StatusBadRequest, "both username password are required")
		return
	}

	// maximum length is determined by bcrypt
	if userRequest.Password != nil && len(*userRequest.Password) < 8 && len([]byte(*userRequest.Password)) >= 72 {
		writeError(w, http.StatusBadRequest, "password must be at least 8 characters")
		return
	}

	user := models.User{VisibleName: &userRequest.VisibleName}
	if userRequest.Username != nil {
		hash, err := bcrypt.GenerateFromPassword([]byte(*userRequest.Password), bcrypt.DefaultCost)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "failed to hash password")
		}
		hashString := string(hash)
		user.Username = userRequest.Username
		user.Password = &hashString
	}

	if err := u.users.Create(&user); err != nil {
		if errors.Is(err, repository.ErrConflict) {
			writeError(w, http.StatusBadRequest, "username already taken")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	resp := models.CreateUserResponse{
		ID: user.ID,
	}

	if userRequest.Username == nil {
		pair, err := u.auth.IssueTokenPair(resp.ID.String(), middleware.RefreshTokenPermanent)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}

		resp.AccessToken = &pair.AccessToken
		resp.RefreshToken = &pair.RefreshToken
		resp.ExpiresAt = &pair.ExpiresAt
	}

	writeJSON(w, http.StatusCreated, resp)
}

// CreateUser generates a user without login and access.
func (u *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var userRequest models.CreateUserRequest
	if err := decode(r, &userRequest); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	requesterID, ok := middleware.GetIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusForbidden, "forbidden")
		return
	}

	var user *models.User
	if userRequest.ID != nil {
		// a request with a user ID will create a login for an existing user
		var err error
		user, err = u.users.GetByID(*userRequest.ID)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}

		if !user.IsActive || user.Username != nil || user.ID != requesterID {
			writeError(w, http.StatusForbidden, "forbidden")
			return
		}

		user.Username = userRequest.Username
		user.Password = userRequest.Password

		if err := u.users.Update(user); err != nil {
			if errors.Is(err, repository.ErrConflict) {
				writeError(w, http.StatusBadRequest, "username already taken")
				return
			}

			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
	} else {
		// a request without login will create a user that has no access.

		if userRequest.VisibleName == "" {
			writeError(w, http.StatusBadRequest, "name required")
			return
		}

		user = &models.User{
			VisibleName: &userRequest.VisibleName,
			OwnerID:     &requesterID,
		}

		if err := u.users.Create(user); err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	resp := models.CreateUserResponse{
		ID: user.ID,
	}

	writeJSON(w, http.StatusCreated, resp)
}

// GetUserSingle reads out user details of a single user. Can only be requested on the connected user.
func (u *UserHandler) GetUserSingle(w http.ResponseWriter, r *http.Request) {
	targetID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	requesterID, ok := middleware.GetIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusForbidden, "forbidden")
		return
	}

	if requesterID != targetID {
		writeError(w, http.StatusForbidden, "forbidden")
		return
	}

	user, err := u.users.GetByID(requesterID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
	}

	writeJSON(w, http.StatusOK, user)
}

// UpdateUser allows updating users visible name and the password.
// Visible name can be updated by the user themselfs, or by the owner.
// Password can only be updated by the user themselfs.
func (u *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	targetID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	requesterID, ok := middleware.GetIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusForbidden, "forbidden")
		return
	}

	var req models.UpdateUserRequest
	if err := decode(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	user, err := u.users.GetByID(targetID)
	if err != nil {
		writeError(w, http.StatusNotFound, "user not found")
		return
	}

	// can only modify themselfs, or an unregistered user that participated at minimum one tournament
	if requesterID != targetID && requesterID != *user.OwnerID {
		writeError(w, http.StatusForbidden, "forbidden")
		return
	}

	if req.VisibleName != nil {
		user.VisibleName = req.VisibleName
	}

	if req.Password != nil {
		if requesterID != targetID || user.Password == nil {
			writeError(w, http.StatusForbidden, "forbidden")
			return
		}

		if len(*req.Password) < 8 {
			writeError(w, http.StatusBadRequest, "password too short")
			return
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.DefaultCost)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "could not hash password")
			return
		}
		s := string(hash)
		user.Password = &s
	}

	if err := u.users.Update(user); err != nil {
		writeError(w, http.StatusInternalServerError, "could not update user")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// DeleteUser deletes a users connection data and visible name. This allows competitions to remain
// valid, but without the deleted user's data.
func (u *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	targetID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	id, ok := middleware.GetIDFromContext(r.Context())
	if !ok || targetID != id {
		writeError(w, http.StatusForbidden, "forbidden")
		return
	}

	repoErr := u.users.SoftDelete(id)
	if repoErr != nil {
		writeError(w, http.StatusInternalServerError, repoErr.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
