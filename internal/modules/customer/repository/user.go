package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/nlsnnn/berezhok/internal/adapters/postgresql/sqlc"
	"github.com/nlsnnn/berezhok/internal/modules/customer/domain"
	sharedDomain "github.com/nlsnnn/berezhok/internal/shared/domain"
)

type UserRepo struct {
	q *sqlc.Queries
}

func NewUserRepo(q *sqlc.Queries) *UserRepo {
	return &UserRepo{q: q}
}

func (r *UserRepo) FindOrCreateByPhone(ctx context.Context, phone string) (domain.User, error) {
	user, err := r.q.FindCustomerByPhone(ctx, phone)
	if err == nil {
		return userToDomain(user), nil
	}

	if err != pgx.ErrNoRows {
		return domain.User{}, err
	}

	id, err := r.q.CreateCustomer(ctx, phone)
	if err != nil {
		return domain.User{}, err
	}

	return domain.User{ID: id, Phone: sharedDomain.Phone{Number: phone}}, nil
}

func (r *UserRepo) FindByID(ctx context.Context, id string) (domain.User, error) {
	uid := uuid.MustParse(id)
	user, err := r.q.FindCustomerByID(ctx, uid)
	if err != nil {
		return domain.User{}, err
	}
	return userToDomain(user), nil
}

func userToDomain(u sqlc.User) domain.User {
	return domain.User{
		ID:    u.ID,
		Phone: sharedDomain.Phone{Number: u.Phone},
		Name:  u.Name.String,
	}
}
