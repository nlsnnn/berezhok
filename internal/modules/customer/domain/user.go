package domain

import sharedDomain "github.com/nlsnnn/berezhok/internal/shared/domain"

type User struct {
	ID    string
	Phone sharedDomain.Phone
	Name  string
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
