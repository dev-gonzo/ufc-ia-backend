package users

import (
	"strings"

	"ufc-backend/internal/auth"
)

func isValidRole(
	role string,
) bool {

	switch role {

	case auth.RoleAdmin:
		return true

	case auth.RoleManager:
		return true

	case auth.RoleUser:
		return true

	default:
		return false
	}
}

func isValidEmail(
	email string,
) bool {

	email = strings.TrimSpace(email)

	return email != ""
}

func isValidUsername(
	username string,
) bool {

	username = strings.TrimSpace(username)

	return len(username) >= 3
}

func isValidPassword(
	password string,
) bool {

	return len(password) >= 6
}
