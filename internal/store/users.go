package store

import (
	"errors"
	"fmt"

	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

// ErrEmailTaken est une erreur qui indique que l'email fourni est déjà utilisé par un autre utilisateur.
var ErrEmailTaken = errors.New("email déjà utilisé")

// ErrInvalidConfirmCode est une erreur qui indique que le code de confirmation fourni est invalide.
var ErrInvalidConfirmCode = errors.New("code de confirmation invalide")

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

// ConfirmUser confirme un utilisateur en vérifiant l'email et le code de confirmation fournis. Si la confirmation est réussie, l'utilisateur est marqué comme confirmé dans la base de données.
func (s *Store) ConfirmUser(email, confirmCode string) error {
	res, err := s.db.Exec(
		`UPDATE users SET confirmed = 1, confirm_code = NULL WHERE email = ? AND confirm_code = ?`,
		email, confirmCode,
	)
	if err != nil {
		return fmt.Errorf("confirmation de l'utilisateur : %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("récupération du nombre de lignes affectées : %w", err)
	}
	if rowsAffected == 0 {
		return ErrInvalidConfirmCode
	}

	return nil
}
