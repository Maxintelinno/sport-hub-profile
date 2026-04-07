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

type CheckPhoneResponse struct {
	Message string `json:"message"`
	IsFound bool   `json:"isFound"`
}

type UpdatePasswordRequest struct {
	Phone       string `json:"phone" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=6"`
}

type UpdatePinRequest struct {
	Phone  string `json:"phone" validate:"required"`
	NewPin string `json:"new_pin" validate:"required,len=6"`
}

type ValidatePinRequest struct {
	Phone string `json:"phone" validate:"required"`
	Pin   string `json:"pin" validate:"required,len=6"`
}

type ValidatePinResponse struct {
	Message string `json:"message"`
	IsValid bool   `json:"is_valid"`
}
