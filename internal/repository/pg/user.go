package pg

import (
	"context"
	"fmt"
	"strings"

	"github.com/clfdrive/server/domain"
	"github.com/clfdrive/server/user"
	"github.com/jackc/pgx/v5"
)

const (
	insertUserQuery = "INSERT INTO users (email, password, verif_code) VALUES ($1, $2, $3) RETIRNING id"
	findUserByEmailQuery = "SELECT * FROM users u WHERE u.email = $1"
)

type UserRepository struct {
	conn *pgx.Conn
}

func NewUserRepository(conn *pgx.Conn) user.UserRepository {
	return &UserRepository{
		conn,
	}
}

func (repo UserRepository) Create(ctx context.Context, user *domain.User) (err error) {
	row := repo.conn.QueryRow(ctx, insertUserQuery, user.Email, user.Password, user.VerifCode)
	if err = row.Scan(&user.ID); err != nil {
		return err
	}

	return nil
}

func (repo UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	var user domain.User
	row := repo.conn.QueryRow(ctx, findUserByEmailQuery, email)
	if err :=  row.Scan(&user); err != nil {
		return domain.User{}, err
	}

	return user, nil
}

func (repo UserRepository) Update(ctx context.Context, userId int, user *domain.UpdateUserDTO) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if user.Verified != nil {
		setValues = append(setValues, fmt.Sprintf("text=$%d", argId))
		args = append(args, *user.Verified)
		argId++
	}

	if user.RefreshToken != nil {
		setValues = append(setValues, fmt.Sprintf("refresh_token=$%d", argId))
		args = append(args, *user.RefreshToken)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")
	query := fmt.Sprintf("UPDATE users SET %s WHERE id=$%d", setQuery, argId)

	args = append(args, userId)

	row := repo.conn.QueryRow(ctx, query, args...)
	if err := row.Scan(); err != nil {
		return err
	}

	return nil
}
