package usecase

import (
	"context"
	"errors"
	"testing"

	"github.dena.jp/swet/go-sampleapi/internal/apierr"

	"github.com/jmoiron/sqlx"
	"github.dena.jp/swet/go-sampleapi/internal/model"
)

type userRepositoryMock struct {
	findByEmailFn func(ctx context.Context, queryer sqlx.QueryerContext, email string) (*model.User, error)
	createFn      func(ctx context.Context, execer sqlx.ExecerContext, m *model.User) error
}

func (s *userRepositoryMock) FindByEmail(ctx context.Context, queryer sqlx.QueryerContext, email string) (*model.User, error) {
	return s.findByEmailFn(ctx, queryer, email)
}

func (s *userRepositoryMock) Create(ctx context.Context, execer sqlx.ExecerContext, m *model.User) error {
	return s.createFn(ctx, execer, m)
}

// Userの作成に成功するケース
func TestUser_Create(t *testing.T) {
	mock := &userRepositoryMock{
		findByEmailFn: func(ctx context.Context, queryer sqlx.QueryerContext, email string) (*model.User, error) {
			return nil, apierr.ErrUserNotExists
		},
		createFn: func(ctx context.Context, execer sqlx.ExecerContext, m *model.User) error {
			return nil
		},
	}
	userUseCase := NewUser(mock, nil)
	if err := userUseCase.Create(context.Background(), &model.User{
		FirstName:    "test_first_name",
		LastName:     "test_last_name",
		Email:        "test@example.com",
		PasswordHash: "aaa",
	}); err != nil {
		t.Fatal(err)
	}
}

// Repositoryからユーザが正しく返ってきたケース
func TestUser_Create_DuplicateEmail(t *testing.T) {
	mock := &userRepositoryMock{
		findByEmailFn: func(ctx context.Context, queryer sqlx.QueryerContext, email string) (*model.User, error) {
			return nil, apierr.ErrEmailAlreadyExists
		},
		createFn: func(ctx context.Context, execer sqlx.ExecerContext, m *model.User) error {
			return nil
		},
	}
	userUseCase := NewUser(mock, nil)
	if err := userUseCase.Create(context.Background(), &model.User{
		FirstName:    "test_first_name",
		LastName:     "test_last_name",
		Email:        "test@example.com",
		PasswordHash: "aaa",
	}); !errors.Is(err, apierr.ErrEmailAlreadyExists) {
		t.Errorf("error must be %v but %v", apierr.ErrEmailAlreadyExists, err)
	}
}

// RepositoryのCreateが失敗したケース
func TestUser_Create_Failed(t *testing.T) {
	mock := &userRepositoryMock{
		findByEmailFn: func(ctx context.Context, queryer sqlx.QueryerContext, email string) (*model.User, error) {
			return nil, apierr.ErrUserNotExists
		},
		createFn: func(ctx context.Context, execer sqlx.ExecerContext, m *model.User) error {
			return errors.New("unknown error")
		},
	}
	userUseCase := NewUser(mock, nil)
	if err := userUseCase.Create(context.Background(), &model.User{
		FirstName:    "test_first_name",
		LastName:     "test_last_name",
		Email:        "test@example.com",
		PasswordHash: "aaa",
	}); err == nil {
		t.Errorf("error should be not nil")
	}
}
