package users

import (
	"encoding/json"
	"errors"
	"net/http"

	"ufc-backend/internal/auth"
	"ufc-backend/internal/shared/http_response"
)

type Handler struct {
	service    *Service
	repository *Repository
}

func NewHandler(
	service *Service,
	repository *Repository,
) *Handler {
	return &Handler{
		service:    service,
		repository: repository,
	}
}

// Create user godoc
//
// @Summary Create user
// @Description Create new user
// @Tags users
// @Accept json
// @Produce json
//
// @Param body body CreateUserInput true "Create user payload"
//
// @Success 201 {object} User
// @Failure 400 {object} httpresponse.ErrorResponse
// @Failure 405 {object} httpresponse.ErrorResponse
//
// @Router /users [post]
func (h *Handler) Create(
	w http.ResponseWriter,
	r *http.Request,
) {

	if r.Method != http.MethodPost {
		httpresponse.Error(
			w,
			http.StatusMethodNotAllowed,
			"METHOD_NOT_ALLOWED",
			"method not allowed",
		)
		return
	}

	var body CreateUserInput

	err := json.NewDecoder(r.Body).Decode(&body)

	if err != nil {
		httpresponse.Error(
			w,
			http.StatusBadRequest,
			"INVALID_BODY",
			"invalid body",
		)
		return
	}

	user, err := h.service.Create(body)

	if err != nil {
		if errors.Is(err, ErrInvalidEmail) {
			httpresponse.Error(
				w,
				http.StatusBadRequest,
				"INVALID_EMAIL",
				"invalid email",
			)
			return
		}

		if errors.Is(err, ErrInvalidUsername) {
			httpresponse.Error(
				w,
				http.StatusBadRequest,
				"INVALID_USERNAME",
				"invalid username",
			)
			return
		}

		if errors.Is(err, ErrInvalidPassword) {
			httpresponse.Error(
				w,
				http.StatusBadRequest,
				"INVALID_PASSWORD",
				"invalid password",
			)
			return
		}

		if errors.Is(err, ErrEmailInUse) {
			httpresponse.Error(
				w,
				http.StatusConflict,
				"EMAIL_ALREADY_EXISTS",
				"email already exists",
			)
			return
		}

		if errors.Is(err, ErrUsernameInUse) {
			httpresponse.Error(
				w,
				http.StatusConflict,
				"USERNAME_ALREADY_EXISTS",
				"username already exists",
			)
			return
		}

		httpresponse.Error(
			w,
			http.StatusInternalServerError,
			"CREATE_USER_FAILED",
			"failed to create user",
		)
		return
	}

	httpresponse.Success(
		w,
		http.StatusCreated,
		user,
	)
}

// List users godoc
//
// @Summary List users
// @Description List all users
// @Tags users
// @Security BearerAuth
// @Produce json
//
// @Success 200 {array} User
// @Failure 401 {object} httpresponse.ErrorResponse
// @Failure 500 {object} httpresponse.ErrorResponse
//
// @Router /users/list [get]
func (h *Handler) List(
	w http.ResponseWriter,
	r *http.Request,
) {

	if r.Method != http.MethodGet {
		httpresponse.Error(
			w,
			http.StatusMethodNotAllowed,
			"METHOD_NOT_ALLOWED",
			"method not allowed",
		)
		return
	}

	users, err := h.repository.List()

	if err != nil {
		httpresponse.Error(
			w,
			http.StatusInternalServerError,
			"LIST_USERS_FAILED",
			"failed to list users",
		)
		return
	}

	httpresponse.Success(
		w,
		http.StatusOK,
		users,
	)
}

// Change password godoc
//
// @Summary Change password
// @Description Change authenticated user password
// @Tags users
// @Security BearerAuth
// @Accept json
// @Produce json
//
// @Param body body ChangePasswordInput true "Change password payload"
//
// @Success 200 {object} map[string]string
// @Failure 400 {object} httpresponse.ErrorResponse
// @Failure 401 {object} httpresponse.ErrorResponse
// @Failure 404 {object} httpresponse.ErrorResponse
//
// @Router /users/change-password [post]
func (h *Handler) ChangePassword(
	w http.ResponseWriter,
	r *http.Request,
) {

	if r.Method != http.MethodPost {
		httpresponse.Error(
			w,
			http.StatusMethodNotAllowed,
			"METHOD_NOT_ALLOWED",
			"method not allowed",
		)
		return
	}

	claims, ok := r.Context().Value(
		auth.UserContextKey,
	).(*auth.Claims)

	if !ok {
		httpresponse.Error(
			w,
			http.StatusUnauthorized,
			"UNAUTHORIZED",
			"unauthorized",
		)
		return
	}

	user, err := h.repository.FindInternalByEmail(
		claims.Email,
	)

	if err != nil {
		httpresponse.Error(
			w,
			http.StatusNotFound,
			"USER_NOT_FOUND",
			"user not found",
		)
		return
	}

	var body ChangePasswordInput

	err = json.NewDecoder(r.Body).Decode(&body)

	if err != nil {
		httpresponse.Error(
			w,
			http.StatusBadRequest,
			"INVALID_BODY",
			"invalid body",
		)
		return
	}

	err = h.service.ChangePassword(
		claims.UserID,
		body,
		user,
	)

	if err != nil {
		if errors.Is(err, ErrInvalidPassword) {
			httpresponse.Error(
				w,
				http.StatusBadRequest,
				"INVALID_PASSWORD",
				"invalid password",
			)
			return
		}

		if errors.Is(err, auth.ErrInvalidCredentials) {
			httpresponse.Error(
				w,
				http.StatusUnauthorized,
				"INVALID_CREDENTIALS",
				"invalid credentials",
			)
			return
		}

		httpresponse.Error(
			w,
			http.StatusInternalServerError,
			"CHANGE_PASSWORD_FAILED",
			"failed to change password",
		)
		return
	}

	httpresponse.Success(
		w,
		http.StatusOK,
		MessageResponse{
			Message: "password changed successfully",
		},
	)
}

// Change role godoc
//
// @Summary Change user role
// @Description Change role of a user
// @Tags users
// @Security BearerAuth
// @Accept json
// @Produce json
//
// @Param body body ChangeRoleInput true "Change role payload"
//
// @Success 200 {object} map[string]string
// @Failure 400 {object} httpresponse.ErrorResponse
// @Failure 401 {object} httpresponse.ErrorResponse
//
// @Router /users/change-role [patch]
func (h *Handler) ChangeRole(
	w http.ResponseWriter,
	r *http.Request,
) {

	if r.Method != http.MethodPatch {

		httpresponse.Error(
			w,
			http.StatusMethodNotAllowed,
			"METHOD_NOT_ALLOWED",
			"method not allowed",
		)

		return
	}

	var body ChangeRoleInput

	err := json.NewDecoder(
		r.Body,
	).Decode(&body)

	if err != nil {

		httpresponse.Error(
			w,
			http.StatusBadRequest,
			"INVALID_BODY",
			"invalid body",
		)

		return
	}

	err = h.service.ChangeRole(body)

	if err != nil {
		if errors.Is(err, ErrInvalidRole) {
			httpresponse.Error(
				w,
				http.StatusBadRequest,
				"INVALID_ROLE",
				"invalid role",
			)
			return
		}

		httpresponse.Error(
			w,
			http.StatusInternalServerError,
			"CHANGE_ROLE_FAILED",
			"failed to change role",
		)

		return
	}

	httpresponse.Success(
		w,
		http.StatusOK,
		MessageResponse{
			Message: "role updated successfully",
		},
	)
}
