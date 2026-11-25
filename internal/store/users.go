package store

import (
	"errors"
	"fmt"

	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

// ErrEmailTaken est une erreur qui indique que l'email fourni est déjà utilisé par un autre utilisateur.
var ErrEmailTaken = errors.New("email déjà utilisé")

// CreateUser crée un nouvel utilisateur dans la base de données avec l'email, le mot de passe hashé et le code de confirmation fournis.
func (s *Store) CreateUser(email, passwordHash, confirmCode string) error {
	_, err := s.db.Exec(
		`INSERT INTO users (email, password_hash, confirm_code) VALUES (?, ?, ?)`,
		email, passwordHash, confirmCode,
	)

	var sqliteErr *sqlite.Error
	if errors.As(err, &sqliteErr) && sqliteErr.Code() == sqlite3.SQLITE_CONSTRAINT_UNIQUE {
		return ErrEmailTaken
	}
	if err != nil {
		return fmt.Errorf("création de l'utilisateur : %w", err)
	}
	return nil
}
