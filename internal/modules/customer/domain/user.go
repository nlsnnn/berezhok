package domain

import (
	"time"

	"github.com/google/uuid"

	sharedDomain "github.com/nlsnnn/berezhok/internal/shared/domain"
)

type User struct {
	ID        uuid.UUID
	Phone     sharedDomain.Phone
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewUser(phone, name string) (User, error) {
	phoneObj, err := sharedDomain.NewPhone(phone)
	if err != nil {
		return User{}, err
	}

	return User{
		Phone: phoneObj,
		Name:  name,
	}, nil
}
