package errors

import "errors"

var (
	ErrInvalidStatusTransition = errors.New("invalid status transition")
	ErrInvalidCredentials      = errors.New("invalid credentials")
	ErrEmailAlreadyInUse       = errors.New("email is already in use")

	ErrPartnerNotFound           = errors.New("partner not found")
	ErrPasswordUnchanged         = errors.New("password must be different from the current one")
	ErrLocationCategoryNotFound  = errors.New("location category not found")
	ErrLocationNotFound          = errors.New("location not found")
	ErrLocationNotOwnedByPartner = errors.New("location does not belong to partner")
	ErrInvalidCommissionRate     = errors.New("commission rate must be between 0 and 100")
	ErrPartnerStatusInvalid      = errors.New("partner status is invalid for this operation")
	ErrInvalidINN                = errors.New("invalid inn")
	ErrInvalidOGRN               = errors.New("invalid ogrn")
	ErrInvalidKPP                = errors.New("invalid kpp")
	ErrEmployeeNotFound          = errors.New("employee not found")
	ErrCannotDeleteOwner         = errors.New("owner account cannot be deleted")
	ErrCannotDeleteSelf          = errors.New("cannot delete your own account")
	ErrInvalidStatsPeriod        = errors.New("invalid stats period")
	ErrInvalidStatsDateRange     = errors.New("invalid stats date range")
	ErrInvalidStatsSort          = errors.New("invalid stats sort")
	ErrInvalidWorkingHours       = errors.New("invalid working hours")
)
