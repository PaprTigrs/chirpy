package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/PaprTigrs/chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerAdminMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	html := fmt.Sprintf(`
<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>
`, cfg.fileserverHits.Load())

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html))
}

func (cfg *apiConfig) handlerAdminReset(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		respondWithError(w, 403, "Forbidden")
	}

	ctx := r.Context()

	err := cfg.db.DeleteAllUsers(ctx)
	if err != nil {
		respondWithError(w, 500, "Could not reset users")
	}

	cfg.fileserverHits.Store(0)

	respondWithJson(w, 200, map[string]string{"status": "ok"})
}

func readinessHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

type createChirpRequest struct {
	Body   string    `json:"body"`
	UserID uuid.UUID `json:"user_id"`
}

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	var req createChirpRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		respondWithError(w, 400, "Invalid JSON")
		return
	}

	if len(req.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long")
		return
	}

	cleaned := cleanChirp(req.Body)

	params := database.CreateChirpParams{
		Body:   cleaned,
		UserID: req.UserID,
	}

	ctx := r.Context()

	dbChirp, err := cfg.db.CreateChirp(ctx, params)
	if err != nil {
		respondWithError(w, 500, "Could not create chirp")
		return
	}

	chirp := Chirp{
		ID:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body:      dbChirp.Body,
		UserID:    dbChirp.UserID,
	}

	respondWithJson(w, 201, chirp)
}

func cleanChirp(body string) string {
	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}

	words := strings.Split(body, " ")
	for i, w := range words {
		lower := strings.ToLower(w)
		if _, exists := badWords[lower]; exists {
			words[i] = "****"
		}
	}

	return strings.Join(words, " ")
}

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	var req createUserRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		respondWithError(w, 400, "Invalid JSON")
	}

	ctx := r.Context()

	dbUser, err := cfg.db.CreateUser(ctx, req.Email)
	if err != nil {
		respondWithError(w, 500, "Could not create user")
	}

	user := User{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Email:     dbUser.Email,
	}

	respondWithJson(w, 201, user)
}
