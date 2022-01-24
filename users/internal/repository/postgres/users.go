package postgres

import (
	"database/sql"
	"errors"
	"github.com/KirillMironov/rapu/users/domain"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"time"
)

type UsersRepository struct {
	db *sqlx.DB
}

func NewUsersRepository(db *sqlx.DB) *UsersRepository {
	return &UsersRepository{db: db}
}

func (u *UsersRepository) Create(user domain.User) (string, error) {
	var sqlStr = "INSERT INTO users (username, email, password, created_at) VALUES ($1, $2, $3, $4) RETURNING id"
	var userId string

	tx, err := u.db.Beginx()
	if err != nil {
		return "", err
	}

	err = tx.QueryRowx(sqlStr, user.Username, user.Email, user.Password, time.Now()).Scan(&userId)
	if err != nil {
		_ = tx.Rollback()
		if e, ok := err.(*pq.Error); ok && e.Code == "23505" {
			return "", domain.ErrUserAlreadyExists
		}
		return "", err
	}

	return userId, tx.Commit()
}

func (u *UsersRepository) GetByEmail(email string) (domain.User, error) {
	var sqlStr = "SELECT id, password FROM users WHERE email = $1"
	var user domain.User

	err := u.db.QueryRowx(sqlStr, email).Scan(&user.Id, &user.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.User{}, domain.ErrUserNotFound
		}
		return domain.User{}, err
	}

	return user, nil
}

func (u *UsersRepository) CheckExistence(userId string) (bool, error) {
	var sqlStr = "SELECT FROM users WHERE id = $1"

	err := u.db.QueryRowx(sqlStr, userId).Scan()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, domain.ErrUserNotFound
		}
		return false, err
	}

	return true, nil
}
