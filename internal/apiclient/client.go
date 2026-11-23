package apiclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Client gère la communication avec l'API et conserve l'état de la session.
// Ses champs sont privés : on ne les modifie qu'à travers ses méthodes.
type Client struct {
	baseURL    string       // adresse de l'API, sans "/" final
	token      string       // jeton de session, "" tant qu'on n'est pas connecté
	httpClient *http.Client // client HTTP de la bibliothèque standard
}

// New crée un Client prêt à l'emploi pour l'API située à baseURL.
// Il faut toujours passer par New : un Client vide aurait un httpClient nil.
func New(baseURL string) *Client {
	return &Client{
		// On retire le "/" final pour éviter "...8080//orders" à la concaténation.
		baseURL: strings.TrimRight(baseURL, "/"),
		httpClient: &http.Client{
			Timeout: 10 * time.Second, // on n'attend jamais l'API indéfiniment
		},
	}
}

// SetToken enregistre le jeton de session après une connexion réussie.
// Récepteur pointeur : la méthode modifie le Client d'origine, pas une copie.
func (c *Client) SetToken(token string) {
	c.token = token
}

// IsAuthenticated indique si un jeton de session est actuellement enregistré.
func (c *Client) IsAuthenticated() bool {
	return c.token != ""
}

// do prépare et envoie une requête vers l'API, puis renvoie la réponse brute.
// body est encodé en JSON s'il n'est pas nil. L'appelant doit fermer resp.Body.
func (c *Client) do(method, path string, body any) (*http.Response, error) {
	// 1. Préparer le corps : rien par défaut
	var reader io.Reader // zero value : nil, donc "pas de corps" (cas d'un GET)
	if body != nil {
		jsonBody, err := json.Marshal(body) // struct Go → octets JSON
		if err != nil {
			return nil, fmt.Errorf("encodage JSON : %w", err)
		}
		reader = bytes.NewReader(jsonBody)
	}

	// 2. Créer la requête
	req, err := http.NewRequest(method, c.baseURL+path, reader)
	if err != nil {
		return nil, fmt.Errorf("création de la requête : %w", err)
	}

	// 3. Les en-têtes
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if c.IsAuthenticated() {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}

	// 4. Envoyer : c'est ici qu'on obtient une erreur si le serveur est éteint
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("envoi de la requête : %w", err)
	}
	return resp, nil
}

// APIError représente une réponse d'erreur de l'API (code HTTP hors 2xx).
// Les interfaces la récupèrent avec errors.As pour réagir selon le code,
// par exemple revenir à l'écran de connexion sur un 401.
type APIError struct {
	StatusCode int    // rempli à la main depuis la réponse HTTP, absent du JSON
	Message    string `json:"error"` // l'API répond {"error": "..."} (voir writeError)
}

// Error fait de APIError une erreur Go : elle satisfait l'interface error.
func (e *APIError) Error() string {
	return fmt.Sprintf("erreur %d : %s", e.StatusCode, e.Message)
}

// call envoie une requête à l'API avec do, puis lit la réponse :
//   - en cas de succès (2xx), le JSON reçu remplit out (laisser out à nil
//     quand il n'y a rien à lire, par exemple un 204 No Content) ;
//   - sinon, elle renvoie une *APIError avec le code HTTP et le message.
func (c *Client) call(method, path string, body, out any) error {
	resp, err := c.do(method, path, body)
	if err != nil {
		return err // resp vaut nil ici : on ne doit surtout pas y toucher
	}
	defer resp.Body.Close() // le corps est fermé à la sortie, quel que soit le chemin

	// Succès : on décode le JSON dans out, si l'appelant attend quelque chose.
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		if out != nil {
			if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
				return fmt.Errorf("décodage de la réponse : %w", err)
			}
		}
		return nil
	}

	// Erreur : on essaie de lire {"error": "..."}. Si le corps n'est pas du JSON
	// (ex. le 404 en texte brut du ServeMux), on ignore l'échec avec _ : le code
	// HTTP compte plus que le message.
	var apiErr APIError
	_ = json.NewDecoder(resp.Body).Decode(&apiErr)
	apiErr.StatusCode = resp.StatusCode
	if apiErr.Message == "" {
		apiErr.Message = http.StatusText(resp.StatusCode) // "Not Found", "Unauthorized"...
	}
	return &apiErr
}
