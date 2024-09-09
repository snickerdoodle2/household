package data

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"inzynierka/internal/data/validator"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Token struct {
	Plaintext string    `json:"token"`
	Hash      []byte    `json:"-"`
	UserID    uuid.UUID `json:"-"`
	Expiry    time.Time `json:"expiry"`
}

func generateToken(userID uuid.UUID, ttl time.Duration) (*Token, error) {
	token := &Token{
		UserID: userID,
		Expiry: time.Now().Add(ttl),
	}

	randomBytes := make([]byte, 16)

	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, err
	}

	token.Plaintext = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)

	hash := sha256.Sum256([]byte(token.Plaintext))
	token.Hash = hash[:]

	return token, nil
}

func ValidateTokenPlaintext(v *validator.Validator, tokenPlaintext string) {
	v.Check(tokenPlaintext != "", "token", "must be provided")
	v.Check(len(tokenPlaintext) == 26, "token", "must be 26 bytes long")
}

type TokenModel struct {
	DB *pgxpool.Pool
}

func (m TokenModel) New(userID uuid.UUID, ttl time.Duration) (*Token, error) {
	token, err := generateToken(userID, ttl)
	if err != nil {
		return nil, err
	}

	err = m.Insert(token)

	return token, err
}

func (m TokenModel) Delete(tokenPlaintext string) error {
	query := `
    DELETE FROM tokens
    WHERE hash = $1
    `

	hash := sha256.Sum256([]byte(tokenPlaintext))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := m.DB.Exec(ctx, query, hash[:])
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return ErrRecordNotFound
	}

	return nil
}

func (m TokenModel) Insert(token *Token) error {
	query := `
    INSERT INTO tokens (hash, user_id, expiry)
    VALUES ($1, $2, $3)
    `

	args := []any{token.Hash, token.UserID, token.Expiry}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := m.DB.Exec(ctx, query, args...)
	return err
}
