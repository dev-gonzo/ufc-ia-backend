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
