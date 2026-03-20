package auth

type User struct {
	UserId   int    `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
