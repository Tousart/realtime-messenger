package dto

type RegisterUserRequest struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
}

type RegisterUserResponse struct {
	SessionID    string `json:"session_id"`
	RedirectPath string `json:"redirect"`
}

type LoginUserRequest struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
}

type LoginUserResponse struct {
	SessionID    string `json:"session_id"`
	RedirectPath string `json:"redirect"`
}

type UserPayload struct {
	UserID   int    `json:"user_id"`
	UserName string `json:"user_name"`
}
