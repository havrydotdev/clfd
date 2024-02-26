package user

import (
	"context"
	"errors"

	"github.com/clfdrive/server/domain"
	"github.com/clfdrive/server/internal/repository"
	"github.com/clfdrive/server/internal/rest"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	FindByEmail(ctx context.Context, email string) (domain.User, error)
	FindByID(ctx context.Context, userId int) (domain.User, error)
	Update(ctx context.Context, userId int, user *domain.UpdateUserDTO) error
	Delete(ctx context.Context, userId int) error
}

type Service struct {
	userRepo UserRepository
}

func NewService(userRepo UserRepository) rest.UserService {
	return &Service{
		userRepo,
	}
}

func (s *Service) Create(ctx context.Context, user *domain.User) error {
	userByEmail, err := s.userRepo.FindByEmail(ctx, user.Email)
	if err != nil {
		if err != repository.ErrNoRows {
			return err
		}
	} else {
		if userByEmail.Verified {
			return errors.New("email_exists")
		}

		err := s.userRepo.Delete(ctx, userByEmail.ID)
		if err != nil {
			return err
		}
	}

	if passLen := len(user.Password); passLen < 12 || passLen > 18 {
		return errors.New("invalid_password_len")
	}

	hashed, err := hashPassword(user.Password)
	if err != nil {
		return err
	}

	user.Password = hashed
	user.VerifCode = genVerifCode()

	err = sendEmail(user.Email, user.VerifCode)
	if err != nil {
		return err
	}

	return s.userRepo.Create(ctx, user)
}

func (s *Service) Verify(ctx context.Context, email string, verifCode string) error {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return err
	}

	if user.VerifCode != verifCode {
		return errors.New("incorrect_verif_code")
	}

	verified := true

	if err := s.userRepo.Update(ctx, user.ID, &domain.UpdateUserDTO{
		Verified: &verified,
	}); err != nil {
		return err
	}

	return nil
}

func (s *Service) SignIn(ctx context.Context, email string, password string) (string, string, error) {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return "", "", err
	}

	if !comparePasswordHash(password, user.Password) {
		return "", "", errors.New("incorrect_password")
	}

	return s.genTokenPairAndUpdate(ctx, user.ID)
}

func (s *Service) Refresh(ctx context.Context, oldRefreshToken string) (string, string, error) {
	userId, err := parseRefreshToken(oldRefreshToken)
	if err != nil {
		return "", "", err
	}

	user, err := s.userRepo.FindByID(ctx, userId)
	if err != nil {
		return "", "", err
	}

	if user.RefreshToken != oldRefreshToken {
		return "", "", errors.New("incorrect_refresh_token")
	}

	return s.genTokenPairAndUpdate(ctx, userId)
}

func (s *Service) FindByID(ctx context.Context, userId int) (domain.User, error) {
	return s.userRepo.FindByID(ctx, userId)
}

func (s *Service) genTokenPairAndUpdate(ctx context.Context, userId int) (string, string, error) {
	accessToken, refreshToken, err := GenerateTokenPair(userId)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, s.userRepo.Update(ctx, userId, &domain.UpdateUserDTO{
		RefreshToken: &refreshToken,
	})
}
