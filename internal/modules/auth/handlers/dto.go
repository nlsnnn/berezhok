package handlers

type LoginEmailPasswordRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=3"`
}

type LoginPhoneRequest struct {
	Phone string `json:"phone" validate:"required,e164"`
}

type LoginResponse struct {
	UserID string `json:"user_id"`
	Token  string `json:"token"`
}
