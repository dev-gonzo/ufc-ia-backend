package httpresponse

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Error ErrorBody `json:"error"`
}

type ErrorBody struct {
	Code    string `json:"code" example:"INVALID_BODY"`
	Message string `json:"message" example:"invalid body"`
}

type SuccessResponse struct {
	Data interface{} `json:"data"`
}

func JSON(
	w http.ResponseWriter,
	status int,
	payload interface{},
) {

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	w.WriteHeader(status)

	json.NewEncoder(w).Encode(payload)
}

func Success(
	w http.ResponseWriter,
	status int,
	data interface{},
) {
	JSON(
		w,
		status,
		SuccessResponse{
			Data: data,
		},
	)
}

func RecoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if recover() != nil {
				Error(
					w,
					http.StatusInternalServerError,
					"INTERNAL_ERROR",
					"internal server error",
				)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func Error(
	w http.ResponseWriter,
	status int,
	code string,
	message string,
) {

	JSON(
		w,
		status,
		ErrorResponse{
			Error: ErrorBody{
				Code:    code,
				Message: message,
			},
		},
	)
}
