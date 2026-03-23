package domain

import "github.com/nlsnnn/berezhok/internal/shared/errors"

type Phone struct {
	Number string
}

func NewPhone(number string) (Phone, error) {
	if len(number) > 15 || len(number) < 7 {
		return Phone{}, errors.ErrInvalidPhoneNumber
	}

	return Phone{Number: number}, nil
}
