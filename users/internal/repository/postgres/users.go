package postgres

import (
	"github.com/KirillMironov/rapu/domain"
	"github.com/jmoiron/sqlx"
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

	err = tx.QueryRowx(sqlStr, user.Username, user.Email, user.Password, time.Now().UTC()).Scan(&userId)
	if err != nil {
		_ = tx.Rollback()
		return "", err
	}

	return userId, tx.Commit()
}

func (u *UsersRepository) GetByEmail(user domain.User) (string, string, error) {
	var sqlStr = "SELECT id, password FROM users WHERE email = $1"
	var userId, password string

	return userId, password, u.db.QueryRowx(sqlStr, user.Email).Scan(&userId, &password)
}
