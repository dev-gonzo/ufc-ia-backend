package users

import (
	"encoding/json"
	"net/http"

	"ufc-backend/internal/auth"
	"ufc-backend/internal/shared/httpresponse"
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
		httpresponse.Error(
			w,
			http.StatusBadRequest,
			"CREATE_USER_FAILED",
			err.Error(),
		)
		return
	}

	httpresponse.JSON(
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

	httpresponse.JSON(
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
		httpresponse.Error(
			w,
			http.StatusBadRequest,
			"CHANGE_PASSWORD_FAILED",
			err.Error(),
		)
		return
	}

	httpresponse.JSON(
		w,
		http.StatusOK,
		map[string]string{
			"message": "password changed successfully",
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

		httpresponse.Error(
			w,
			http.StatusBadRequest,
			"CHANGE_ROLE_FAILED",
			err.Error(),
		)

		return
	}

	httpresponse.JSON(
		w,
		http.StatusOK,
		map[string]string{
			"message": "role updated successfully",
		},
	)
}
