package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/KirillMironov/rapu/users/internal/domain"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"time"
)

type Users struct {
	db *sqlx.DB
}

func NewUsers(db *sqlx.DB) *Users {
	return &Users{db: db}
}

func (u *Users) Create(ctx context.Context, user domain.User) (string, error) {
	var sqlStr = "INSERT INTO users (username, email, password, created_at) VALUES ($1, $2, $3, $4) RETURNING id"
	var userId string

	tx, err := u.db.Beginx()
	if err != nil {
		return "", err
	}

	err = tx.QueryRowxContext(ctx, sqlStr, user.Username, user.Email, user.Password, time.Now()).Scan(&userId)
	if err != nil {
		_ = tx.Rollback()
		if e, ok := err.(*pq.Error); ok && e.Code == "23505" {
			return "", domain.ErrUserAlreadyExists
		}
		return "", err
	}

	return userId, tx.Commit()
}

func (u *Users) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	var sqlStr = "SELECT id, password FROM users WHERE email = $1"
	var user domain.User

	err := u.db.QueryRowxContext(ctx, sqlStr, email).Scan(&user.Id, &user.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.User{}, domain.ErrUserNotFound
		}
		return domain.User{}, err
	}

	return user, nil
}

func (u *Users) CheckExistence(ctx context.Context, userId string) (bool, error) {
	var sqlStr = "SELECT FROM users WHERE id = $1"

	err := u.db.QueryRowxContext(ctx, sqlStr, userId).Scan()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, domain.ErrUserNotFound
		}
		return false, err
	}

	return true, nil
}
