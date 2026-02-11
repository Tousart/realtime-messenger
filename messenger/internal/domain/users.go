package domain

type RegisterRequest struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
}

type LoginRequest struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
}

type User struct {
	UserID   int    `json:"user_id"`
	UserName string `json:"user_name"`
	Password string `json:"password"`
}
