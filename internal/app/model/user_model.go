package model

type AddUserRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	FullName string `json:"fullName" validate:"required"`
	Role     uint   `json:"role" validate:"required"`
}

type EditUserRequest struct {
	ID uint `json:"user_id" validate:"required"`
	AddUserRequest
}

type User struct {
	ID       uint   `json:"id"`
	ClientID uint   `json:"client_id"`
	Username string `json:"username"`
	FullName string `json:"fullName"`
	Role     uint   `json:"role"`
	IsAdmin  bool   `json:"is_admin"`
	IsLogin  bool   `json:"is_login"`
}

type GetUserResponse struct {
	HTTPResponse
	Data *[]User `json:"data,omitempty"`
}
