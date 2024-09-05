package data

import (
	"context"
	"errors"
	"inzynierka/internal/data/validator"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrDuplicateUsername = errors.New("duplicate username")
)

type password struct {
	plaintext *string
	hash      []byte
}

func (p *password) Set(plaintext string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintext), 12)
	if err != nil {
		return err
	}

	p.plaintext = &plaintext
	p.hash = hash

	return nil
}

func (p *password) Matches(plaintext string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plaintext))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}

type User struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Name      string    `json:"name"`
	Password  password  `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	Version   int       `json:"-"`
}

func ValidatePasswordPlain(v *validator.Validator, password string) {
	v.Check(password != "", "password", "must be provided")
	v.Check(utf8.RuneCountInString(password) >= 8, "password", "must be at least 8 characters long")
	v.Check(utf8.RuneCountInString(password) <= 64, "password", "must not be more than 64 characters long")
}

func ValidateUsername(v *validator.Validator, username string) {
	v.Check(username != "", "username", "must be provided")
	v.Check(utf8.RuneCountInString(username) >= 4, "username", "must be at least 8 characters long")
	v.Check(utf8.RuneCountInString(username) <= 32, "username", "must not be more than 32 characters long")
}

func ValidateUser(v *validator.Validator, user *User) {
	v.Check(user.Name != "", "name", "must be provided")
	v.Check(utf8.RuneCountInString(user.Name) <= 256, "name", "must not be more than 32 characters long")

	ValidateUsername(v, user.Username)

	if user.Password.plaintext != nil {
		ValidatePasswordPlain(v, *user.Password.plaintext)
	}

	if user.Password.hash == nil {
		panic("missing password hash")
	}
}

type UserModel struct {
	DB *pgxpool.Pool
}

func (m UserModel) Insert(user *User) error {
	query := `
    INSERT INTO users (id, username, display_name, password_hash)
    VALUES ($1, $2, $3, $4)
    RETURNING created_at, version
    `
	id, err := uuid.NewRandom()
	if err != nil {
		return err
	}

	user.ID = id

	args := []any{id, user.Username, user.Name, user.Password.hash}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = m.DB.QueryRow(ctx, query, args...).Scan(&user.CreatedAt, &user.Version)
	if err != nil {
		switch {
		case strings.HasPrefix(err.Error(), "ERROR: duplicate key value violates unique constraint \"users_username_key\""):
			return ErrDuplicateUsername
		default:
			return err
		}
	}

	return nil
}

func (m UserModel) GetByUsername(username string) (*User, error) {
	query := `
    SELECT id, username, display_name, password_hash, created_at, version
    FROM users
    WHERE username = $1
    `

	var user User

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := m.DB.QueryRow(ctx, query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Name,
		&user.Password.hash,
		&user.CreatedAt,
		&user.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}

func (m UserModel) Update(user *User) error {
	query := `
    UPDATE users
    SET display_name = $1, password_hash = $2, version = version + 1
    WHERE id = $3
    RETURNING version
    `

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return m.DB.QueryRow(ctx, query, user.Name, user.Password.hash, user.ID).Scan(&user.Version)
}

func (m UserModel) DeleteByUsername(username string) error {
	query := `
    DELETE FROM users
    WHERE username = $1
    `

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := m.DB.Exec(ctx, query, username)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrRecordNotFound
	}

	return nil
}
