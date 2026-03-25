package models

type ForgotPasswordRequest struct {
	Phone       string `json:"phone" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=6"`
}

type AuthResponse struct {
	Message string `json:"message"`
}

type CheckPhoneRequest struct {
	Phone string `json:"phone" validate:"required"`
}
