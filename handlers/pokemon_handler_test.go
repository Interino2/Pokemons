package handlers

import (
	"context"
	"errors"
	"html/template"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5"

	"pokemon-app/models"
)

type testPokemonStore struct {
	pokemons []models.Pokemon
	pokemon  models.Pokemon
	err      error
}

func (s testPokemonStore) GetAll(_ context.Context) ([]models.Pokemon, error) {
	return s.pokemons, s.err
}

func (s testPokemonStore) GetByID(_ context.Context, _ int) (models.Pokemon, error) {
	return s.pokemon, s.err
}

func newTestHandler(t *testing.T, store PokemonStore) *PokemonHandler {
	t.Helper()

	tmpl, err := template.ParseFiles("../templates/index.html")
	if err != nil {
		t.Fatalf("parse templates: %v", err)
	}
	return NewPokemonHandler(store, tmpl)
}

func TestHomeRendersPokemonCollection(t *testing.T) {
	handler := newTestHandler(t, testPokemonStore{
		pokemons: []models.Pokemon{
			{ID: 25, Nombre: "Pikachu", Tipo: "Eléctrico"},
		},
	})
	response := httptest.NewRecorder()

	handler.Home(response, httptest.NewRequest("GET", "/", nil))

	if response.Code != 200 {
		t.Fatalf("status = %d, want 200", response.Code)
	}
	body := response.Body.String()
	for _, want := range []string{"Pikachu", "Eléctrico", "N.º 025", "Pokémon en la colección"} {
		if !strings.Contains(body, want) {
			t.Errorf("response body does not contain %q", want)
		}
	}
	if got := response.Header().Get("Content-Type"); got != "text/html; charset=utf-8" {
		t.Errorf("Content-Type = %q, want text/html; charset=utf-8", got)
	}
}

func TestGetPokemonEscapesDatabaseContent(t *testing.T) {
	handler := newTestHandler(t, testPokemonStore{
		pokemon: models.Pokemon{
			ID:          7,
			Nombre:      "K'; <script>alert(1)</script>",
			Tipo:        "Planta",
			Descripcion: "<img src=x onerror=alert(1)>",
			Imagen:      "https://example.com/bulbasaur.png",
		},
	})
	request := httptest.NewRequest("GET", "/pokemon/7", nil)
	request.SetPathValue("id", "7")
	response := httptest.NewRecorder()

	handler.GetPokemon(response, request)

	if response.Code != 200 {
		t.Fatalf("status = %d, want 200", response.Code)
	}
	body := response.Body.String()
	for _, unsafe := range []string{"<script>alert(1)</script>", "<img src=x onerror=alert(1)>"} {
		if strings.Contains(body, unsafe) {
			t.Errorf("response contains unescaped database content %q", unsafe)
		}
	}
	if !strings.Contains(body, `id="pokemon-info"`) {
		t.Error("response does not contain the Datastar target element")
	}
}

func TestGetPokemonReturnsNotFound(t *testing.T) {
	handler := newTestHandler(t, testPokemonStore{err: pgx.ErrNoRows})
	request := httptest.NewRequest("GET", "/pokemon/999", nil)
	request.SetPathValue("id", "999")
	response := httptest.NewRecorder()

	handler.GetPokemon(response, request)

	if response.Code != 404 {
		t.Errorf("status = %d, want 404", response.Code)
	}
}

func TestHomeReportsRepositoryFailure(t *testing.T) {
	handler := newTestHandler(t, testPokemonStore{err: errors.New("database offline")})
	response := httptest.NewRecorder()

	handler.Home(response, httptest.NewRequest("GET", "/", nil))

	if response.Code != 500 {
		t.Errorf("status = %d, want 500", response.Code)
	}
	if !strings.Contains(response.Body.String(), "No se pudieron cargar") {
		t.Errorf("unexpected error response: %q", response.Body.String())
	}
}
