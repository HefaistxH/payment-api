package entity

type Credential struct {
	Id       string `json:"id"`
	UserId   string `json:"user_id"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Pin      string `json:"pin"`
	Role     string `json:"role"`
}
