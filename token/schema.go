package token

type User struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone int    `json:"phone"`
	Role  string `json:"role"`
}
