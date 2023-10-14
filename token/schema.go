package token

import (
	"github.com/google/uuid"
)

type User struct {
	Id    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Email string    `json:"email"`
	Phone int       `json:"phone"`
	Role  string    `json:"role"`
}

type TokenData struct {
	Id  string
	jwt string
}
