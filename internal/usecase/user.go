package usecase

import (
	"context"
	"fmt"

	"github.dena.jp/swet/go-sampleapi/internal/apierr"

	"github.com/jmoiron/sqlx"

	"github.dena.jp/swet/go-sampleapi/internal/model"
)

type userRepository interface {
	FindByEmail(ctx context.Context, queryer sqlx.QueryerContext, email string) (*model.User, error)
	Create(ctx context.Context, execer sqlx.ExecerContext, m *model.User) error
}

type User struct {
	userRepo userRepository
	db       *sqlx.DB
}

func NewUser(userRepo userRepository, db *sqlx.DB) *User {
	return &User{
		userRepo: userRepo,
		db:       db,
	}
}

func (u *User) Create(ctx context.Context, m *model.User) error {
	// TODO methodを実装する
	_, err := u.userRepo.FindByEmail(ctx, u.db, m.Email)
	if err != apierr.ErrUserNotExists && err != nil {
		return fmt.Errorf("user find with email %s failed: %w", m.Email, err)
	}
	if err == nil {
		return apierr.ErrEmailAlreadyExists
	}

	if err := u.userRepo.Create(ctx, u.db, m); err != nil {
		return fmt.Errorf("user create failed: with %w", err)
	}
	return nil
}
