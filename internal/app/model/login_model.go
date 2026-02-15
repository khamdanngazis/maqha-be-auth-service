package model

// LoginRequest represents the structure of a login request.
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	HTTPResponse
	Data *struct {
		Token string `json:"token"`
	} `json:"data,omitempty"`
}
