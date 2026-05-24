package handlers

type LoginEmailPasswordRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=3"`
}

type SendCodeRequest struct {
	Phone string `json:"phone" validate:"required,e164"`
}

type LoginPhoneRequest struct {
	Phone string `json:"phone" validate:"required,e164"`
	Code  string `json:"code" validate:"required,len=6"`
}

type LoginResponse struct {
	UserID string `json:"user_id"`
	Token  string `json:"token"`
}

type LoginPartnerResponse struct {
	UserID     string  `json:"user_id"`
	EmployeeID string  `json:"employee_id"`
	PartnerID  string  `json:"partner_id"`
	LocationID *string `json:"location_id,omitempty"`
	Role       string  `json:"role"`
	Token      string  `json:"token"`
	MustChange bool    `json:"must_change_password"`
}

type LoginAdminResponse struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Name   string `json:"name"`
	Role   string `json:"role"`
	Token  string `json:"token"`
}
