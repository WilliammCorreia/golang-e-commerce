package api

import (
	"net/http"

	"e-commerce/internal/store"
)

type API struct {
	store *store.Store
}

func NewRouter(s *store.Store) http.Handler {
	a := &API{store: s}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", a.handleHealth)
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
