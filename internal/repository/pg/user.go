package pg

import (
	"context"

	"github.com/clfdrive/server/domain"
	"github.com/jackc/pgx/v5"
)

type UserRepository struct {
	conn *pgx.Conn
}

func NewUserRepository(conn *pgx.Conn) *UserRepository {
	return &UserRepository{
		conn,
	}
}

func (repo UserRepository) Create(ctx context.Context, user *domain.User) error {
	panic("not implemented") // TODO: Implement
}

func (repo UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	panic("not implemented") // TODO: Implement
}

func (repo UserRepository) Update(ctx context.Context, userId int, user *domain.User) error {
	panic("not implemented") // TODO: Implement
}

