package api

import (
	"net/http"

	"e-commerce/internal/store"
)

type API struct {
	store  *store.Store
	mailer Mailer
}

type Mailer interface {
	Send(to, subject, body string) error
}

func NewRouter(s *store.Store, m Mailer) http.Handler {
	a := &API{store: s, mailer: m}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", a.handleHealth)
	mux.HandleFunc("POST /auth/register", a.handleRegister)
	return mux
}

func (a *API) handleHealth(w http.ResponseWriter, r *http.Request) {
	body := map[string]string{"status": "ok"}
	if err := a.store.Ping(); err != nil {
		writeError(w, http.StatusServiceUnavailable, "la connexion avec la base de données a échoué")
		return
	}
	writeJSON(w, http.StatusOK, body)
}
