package pg

import (
	"context"
	"fmt"
	"strings"

	"github.com/clfdrive/server/domain"
	"github.com/clfdrive/server/internal/repository"
	"github.com/clfdrive/server/user"
	"github.com/jackc/pgx/v5"
)

// TODO: refactor queries to use * instead of all fields
const (
	insertUserQuery      = "INSERT INTO users (email, password, verif_code) VALUES ($1, $2, $3) RETURNING id"
	findUserByEmailQuery = "SELECT id, email, password, refresh_token, verified, verif_code FROM users WHERE email = $1"
	findUserByIdQuery    = "SELECT id, email, password, refresh_token, verified, verif_code FROM users WHERE id = $1"
	deleteUserQuery      = "DELETE FROM users WHERE id = $1"
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
	var res domain.User
	row := repo.conn.QueryRow(ctx, findUserByEmailQuery, email)
	if err := row.Scan(
		&res.ID,
		&res.Email,
		&res.Password,
		&res.RefreshToken,
		&res.Verified,
		&res.VerifCode,
	); err != nil {
		if err == pgx.ErrNoRows {
			return domain.User{}, repository.ErrNoRows
		}

		return domain.User{}, err
	}

	return res, nil
}

func (repo UserRepository) FindByID(ctx context.Context, userId int) (domain.User, error) {
	var res domain.User
	row := repo.conn.QueryRow(ctx, findUserByIdQuery, userId)
	if err := row.Scan(
		&res.ID,
		&res.Email,
		&res.Password,
		&res.RefreshToken,
		&res.Verified,
		&res.VerifCode,
	); err != nil {
		if err == pgx.ErrNoRows {
			return domain.User{}, repository.ErrNoRows
		}

		return domain.User{}, err
	}

	return res, nil
}

func (repo UserRepository) Update(ctx context.Context, userId int, user *domain.UpdateUserDTO) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if user.Verified != nil {
		setValues = append(setValues, fmt.Sprintf("verified=$%d", argId))
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
	if err := row.Scan(); err != nil && err != pgx.ErrNoRows {
		return err
	}

	return nil
}

func (repo UserRepository) Delete(ctx context.Context, userId int) error {
	_, err := repo.conn.Exec(ctx, deleteUserQuery, userId)
	if err != nil {
		return err
	}

	return nil
}
