package routes

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/stretchr/testify/assert"
)

func TestPingEndpoint(t *testing.T) {
	// Setup
	router := SetupRouter(cache.New(5*time.Minute, 10*time.Minute))

	// Test request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "success")
}

func TestGetPokemonById(t *testing.T) {
	// Setup
	router := SetupRouter(cache.New(5*time.Minute, 10*time.Minute))

	tests := []struct {
		name       string
		identifier string
		wantStatus int
	}{
		{"valid pokemon id", "1", http.StatusOK},
		{"valid pokemon name", "bulbasaur", http.StatusOK},
		{"invalid pokemon", "invalid-pokemon", http.StatusNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/pokemon/"+tt.identifier, nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func TestGetAbilities(t *testing.T) {
	router := SetupRouter(cache.New(5*time.Minute, 10*time.Minute))

	t.Run("abilities for eevee", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/pokemon/eevee/abilities", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "run-away")
		assert.Contains(t, w.Body.String(), "adaptability")
	})
}

func TestGetEvolutionChain(t *testing.T) {
	router := SetupRouter(cache.New(5*time.Minute, 10*time.Minute))

	t.Run("evolution chain for eevee", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/pokemon/eevee/evolutions", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "vaporeon") // Eevee evolves to Vaporeon
		assert.Contains(t, w.Body.String(), "jolteon") // Eevee evolves to Jolteon
	})
}
