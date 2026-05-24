package domain

type LegalInfoVerificationStatus string

const (
	LegalInfoVerificationPending  LegalInfoVerificationStatus = "pending"
	LegalInfoVerificationVerified LegalInfoVerificationStatus = "verified"
	LegalInfoVerificationFailed   LegalInfoVerificationStatus = "failed"
)

type LegalInfo struct {
	PartnerID           string
	Inn                 string
	Ogrn                string
	Kpp                 string
	LegalAddress        string
	LegalName           string
	VerificationStatus  LegalInfoVerificationStatus
	VerificationComment string
}
