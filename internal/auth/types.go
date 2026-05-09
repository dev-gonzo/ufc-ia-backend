package auth

type User struct {
	ID           string
	Email        string
	Role         string
	PasswordHash string
}
