package dto

type RegisterUserRequest struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
}

type RegisterUserResponse struct {
	SessionID string `json:"session_id"`
}

type LoginUserRequest struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
}

type LoginUserResponse struct {
	SessionID string `json:"session_id"`
}

type UserPayload struct {
	UserID   int    `json:"user_id"`
	UserName string `json:"user_name"`
}
