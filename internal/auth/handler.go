package auth

import (
	"encoding/json"
	"net/http"

	"ufc-backend/internal/shared/httpresponse"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

type LoginInput struct {
	Email    string `json:"email" example:"admin@email.com"`
	Password string `json:"password" example:"123456"`
}

type LoginResponse struct {
	Token string `json:"token" example:"jwt-token"`
}

// Login godoc
//
// @Summary Login
// @Description Authenticate user and return JWT token
// @Tags auth
// @Accept json
// @Produce json
//
// @Param body body LoginInput true "Login payload"
//
// @Success 200 {object} LoginResponse
// @Failure 400 {object} httpresponse.ErrorResponse
// @Failure 401 {object} httpresponse.ErrorResponse
// @Failure 405 {object} httpresponse.ErrorResponse
//
// @Router /login [post]
func (h *Handler) Login(
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

	var body LoginInput

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

	token, err := h.service.Login(
		body.Email,
		body.Password,
	)

	if err != nil {
		httpresponse.Error(
			w,
			http.StatusUnauthorized,
			"INVALID_CREDENTIALS",
			"invalid credentials",
		)
		return
	}

	httpresponse.JSON(
		w,
		http.StatusOK,
		LoginResponse{
			Token: token,
		},
	)
}
