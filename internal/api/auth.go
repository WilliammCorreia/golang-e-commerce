package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	"e-commerce/internal/ids"
	"e-commerce/internal/password"
	"e-commerce/internal/store"
)

type registerRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type confirmRequest struct {
	Email       string `json:"email"`
	ConfirmCode string `json:"confirm_code"`
}

func (a *API) handleRegister(w http.ResponseWriter, r *http.Request) {
	req := registerRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "erreur lors de la récupération de l'email et du password")
		return
	}

	email := strings.ToLower(strings.TrimSpace(req.Email))
	if !strings.Contains(email, "@") {
		writeError(w, http.StatusBadRequest, "format de l'email invalide")
		return
	}

	if len(req.Password) < 12 {
		writeError(w, http.StatusBadRequest, "le mot de passe doit contenir au moins 12 caractères")
		return
	}

	passwordHashed, err := password.Hash(req.Password)
	if err != nil {
		log.Printf("register : %v", err)
		writeError(w, http.StatusInternalServerError, "une erreur est survenue")
		return
	}

	confirmCode := ids.Code()

	if err := a.store.CreateUser(email, passwordHashed, confirmCode); err != nil {
		if errors.Is(err, store.ErrEmailTaken) {
			writeError(w, http.StatusConflict, "email déjà utilisé")
			return
		}
		log.Printf("register : %v", err)
		writeError(w, http.StatusInternalServerError, "une erreur est survenue")
		return
	}

	if err := a.mailer.Send(email, "Confirmez votre inscription", "Ton code : "+confirmCode); err != nil {
		log.Printf("register : %v", err)
	}

	body := map[string]any{
		"message": "compte créé avec succès",
		"data":    map[string]string{"email": email},
	}
	writeJSON(w, http.StatusCreated, body)
}

func (a *API) handleConfirm(w http.ResponseWriter, r *http.Request) {
	req := confirmRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "erreur lors de la récupération de l'email et du code de confirmation")
		return
	}

	email := strings.ToLower(strings.TrimSpace(req.Email))
	if !strings.Contains(email, "@") {
		writeError(w, http.StatusBadRequest, "format de l'email invalide")
		return
	}

	confirmCode := strings.ToUpper(strings.TrimSpace(req.ConfirmCode))
	if len(confirmCode) == 0 {
		writeError(w, http.StatusBadRequest, "le code de confirmation est requis")
		return
	}

	if err := a.store.ConfirmUser(email, confirmCode); err != nil {
		if errors.Is(err, store.ErrInvalidConfirmCode) {
			writeError(w, http.StatusBadRequest, "code de confirmation invalide")
			return
		}
		log.Printf("confirm : %v", err)
		writeError(w, http.StatusInternalServerError, "une erreur est survenue")
		return
	}

	body := map[string]any{
		"message": "compte confirmé avec succès",
		"data":    map[string]string{"email": email},
	}
	writeJSON(w, http.StatusOK, body)
}
