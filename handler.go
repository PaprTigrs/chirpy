package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type chirpRequest struct {
	Body string `json:"body"`
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondWithJson(w, code, map[string]string{"error": msg})
}

func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"Internal server error"}`))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}

func cleanChirp(body string) string {
	badWords := []string{"kerfuffle", "sharbert", "fornax"}

	words := strings.Split(body, " ")

	for i, word := range words {
		lower := strings.ToLower(word)

		for _, bad := range badWords {
			if lower == bad {
				words[i] = "****"
			}
		}
	}

	return strings.Join(words, " ")
}

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	hits := cfg.fileserverHits.Load()
	html := fmt.Sprintf(`
	<html>
		<body>
			<h1>Welcome, Chirpy Admin</h1>
			<p>Chirpy has been visited %d times!</p>
		</body>
	</html>
	`, hits)

	w.Write([]byte(html))
}

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))

	cfg.fileserverHits.Store(0)
}

func (cfg *apiConfig) handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	var req chirpRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	if len(req.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	cleaned := cleanChirp(req.Body)

	respondWithJson(w, http.StatusOK, map[string]string{
		"cleaned_body": cleaned,
	})
}
