package user

import (
	"context"
	"errors"

	"github.com/clfdrive/server/domain"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	FindByEmail(ctx context.Context, email string) (domain.User, error)
	Update(ctx context.Context, userId int, user *domain.UpdateUserDTO) error
}

type Service struct {
	userRepo UserRepository
}

func NewService(userRepo UserRepository) *Service {
	return &Service{
		userRepo,
	}
}

func (s *Service) Create(ctx context.Context, user *domain.User) error {
	hashed, err  := HashPassword(user.Password)
	if err != nil {
		return err
	}

	user.Password = hashed
	user.VerifCode = uuid.NewString()

	return s.userRepo.Create(ctx, user)
}

func (s *Service) Verify(ctx context.Context, email string, verifCode string) error {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return err
	}

	if user.VerifCode != verifCode {
		return errors.New("invalid verification code")
	}

	verified := true
	
	if err := s.userRepo.Update(ctx, user.ID, &domain.UpdateUserDTO{
		Verified: &verified,
	}); err != nil {
		return err
	}

	return nil
}

func (s *Service) SignIn(ctx context.Context, email string, password string) error {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return err
	}

	if !ComparePasswordHash(password, user.Password) {
		return errors.New("incorrect password")
	}

	return nil
}

func (s *Service) UpdateRefreshToken(ctx context.Context, userId int, refreshToken string) error {
	return s.userRepo.Update(ctx, userId, &domain.UpdateUserDTO{
		RefreshToken: &refreshToken,
	})
}

func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
    return string(bytes), err
}

func ComparePasswordHash(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}
