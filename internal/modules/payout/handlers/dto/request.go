package dto

type UpsertDestinationRequest struct {
	Type          string `json:"type"           validate:"required,oneof=sbp"`
	SBPPhone      string `json:"sbp_phone"      validate:"required,e164"`
	SBPBankID     string `json:"sbp_bank_id"    validate:"required,min=1,max=50"`
	RecipientName string `json:"recipient_name" validate:"required,min=2,max=200"`
}
