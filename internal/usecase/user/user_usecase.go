package user

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/podoru/spinner-podoru/internal/domain/entity"
	"github.com/podoru/spinner-podoru/internal/domain/repository"
	"github.com/podoru/spinner-podoru/pkg/crypto"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrInvalidPassword   = errors.New("invalid current password")
	ErrSamePassword      = errors.New("new password must be different from current password")
)

type UseCase struct {
	userRepo repository.UserRepository
}

func NewUseCase(userRepo repository.UserRepository) *UseCase {
	return &UseCase{userRepo: userRepo}
}

func (uc *UseCase) GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	user, err := uc.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

func (uc *UseCase) Update(ctx context.Context, id uuid.UUID, input *entity.UserUpdate) (*entity.User, error) {
	user, err := uc.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	if input.Name != nil {
		user.Name = *input.Name
	}
	if input.AvatarURL != nil {
		user.AvatarURL = input.AvatarURL
	}
	user.UpdatedAt = time.Now()

	if err := uc.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (uc *UseCase) UpdatePassword(ctx context.Context, id uuid.UUID, input *entity.UserPasswordUpdate) error {
	user, err := uc.userRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrUserNotFound
	}

	if !crypto.CheckPassword(input.CurrentPassword, user.PasswordHash) {
		return ErrInvalidPassword
	}

	if input.CurrentPassword == input.NewPassword {
		return ErrSamePassword
	}

	newHash, err := crypto.HashPassword(input.NewPassword)
	if err != nil {
		return err
	}

	user.PasswordHash = newHash
	user.UpdatedAt = time.Now()

	return uc.userRepo.Update(ctx, user)
}
