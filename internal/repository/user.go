package repository

import (
	"context"
	"database/sql"

	"github.dena.jp/swet/go-sampleapi/internal/apierr"

	"github.com/jmoiron/sqlx"

	"github.dena.jp/swet/go-sampleapi/internal/model"
)

type User struct{}

func NewUser() *User {
	return &User{}
}

func (u *User) FindByEmail(ctx context.Context, queryer sqlx.QueryerContext, email string) (*model.User, error) {
	// TODO methodを実装する
	var m model.User
	if err := sqlx.GetContext(ctx, queryer, &m, "select * from users where email = ?", email); err == sql.ErrNoRows {
		return nil, apierr.ErrUserNotExists
	} else if err != nil {
		return nil, err
	}
	return &m, nil
}

func (u *User) Create(ctx context.Context, execer sqlx.ExecerContext, m *model.User) error {
	// TODO methodを実装する
	result, err := execer.ExecContext(ctx, "insert into users(first_name, last_name, email, password_hash) values(?, ?, ?, ?)", m.FirstName, m.LastName, m.Email, m.PasswordHash)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	m.ID = int(id)
	return nil
}
