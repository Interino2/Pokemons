package handlers

import (
	"context"
	"errors"
	"html/template"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5"

	"pokemon-app/models"
)

type testPokemonStore struct {
	pokemons []models.Pokemon
	pokemon  models.Pokemon
	err      error
	created  models.Pokemon
}

func (s testPokemonStore) GetAll(_ context.Context) ([]models.Pokemon, error) {
	return s.pokemons, s.err
}

func (s testPokemonStore) GetByID(_ context.Context, _ int) (models.Pokemon, error) {
	return s.pokemon, s.err
}

func (s *testPokemonStore) Create(_ context.Context, pokemon models.Pokemon) error {
	s.created = pokemon
	return s.err
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
	handler := newTestHandler(t, &testPokemonStore{
		pokemons: []models.Pokemon{
			{ID: 25, PokedexNumber: 25, Nombre: "Pikachu", Tipo: "Eléctrico"},
		},
	})
	response := httptest.NewRecorder()

	handler.Home(response, httptest.NewRequest("GET", "/", nil))

	if response.Code != 200 {
		t.Fatalf("status = %d, want 200", response.Code)
	}
	body := response.Body.String()
	for _, want := range []string{"Pikachu", "Eléctrico", "N.º 025", "Pokémon en la colección", "Sortear mi equipo", "seis compañeros distintos"} {
		if !strings.Contains(body, want) {
			t.Errorf("response body does not contain %q", want)
		}
	}
	if got := response.Header().Get("Content-Type"); got != "text/html; charset=utf-8" {
		t.Errorf("Content-Type = %q, want text/html; charset=utf-8", got)
	}
}

func TestGetPokemonEscapesDatabaseContent(t *testing.T) {
	handler := newTestHandler(t, &testPokemonStore{
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
	handler := newTestHandler(t, &testPokemonStore{err: pgx.ErrNoRows})
	request := httptest.NewRequest("GET", "/pokemon/999", nil)
	request.SetPathValue("id", "999")
	response := httptest.NewRecorder()

	handler.GetPokemon(response, request)

	if response.Code != 404 {
		t.Errorf("status = %d, want 404", response.Code)
	}
}

func TestHomeReportsRepositoryFailure(t *testing.T) {
	handler := newTestHandler(t, &testPokemonStore{err: errors.New("database offline")})
	response := httptest.NewRecorder()

	handler.Home(response, httptest.NewRequest("GET", "/", nil))

	if response.Code != 500 {
		t.Errorf("status = %d, want 500", response.Code)
	}
	if !strings.Contains(response.Body.String(), "No se pudieron cargar") {
		t.Errorf("unexpected error response: %q", response.Body.String())
	}
}

func TestCreatePokemonPersistsAndRedirects(t *testing.T) {
	store := &testPokemonStore{}
	handler := newTestHandler(t, store)
	form := url.Values{
		"nombre":      {"Pikachu"},
		"tipo":        {"Eléctrico"},
		"altura":      {"0.4"},
		"peso":        {"6.0"},
		"imagen":      {"https://example.com/pikachu.png"},
		"descripcion": {"Un Pokémon eléctrico."},
	}
	request := httptest.NewRequest("POST", "/pokemon", strings.NewReader(form.Encode()))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	response := httptest.NewRecorder()

	handler.CreatePokemon(response, request)

	if response.Code != http.StatusSeeOther {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusSeeOther)
	}
	if response.Header().Get("Location") != "/" {
		t.Errorf("Location = %q, want /", response.Header().Get("Location"))
	}
	if store.created.Nombre != "Pikachu" || store.created.Altura != 0.4 || store.created.Peso != 6 {
		t.Errorf("unexpected Pokémon saved: %+v", store.created)
	}
}

func TestCreatePokemonRendersValidationError(t *testing.T) {
	store := &testPokemonStore{}
	handler := newTestHandler(t, store)
	form := url.Values{
		"nombre": {"Pikachu"},
		"tipo":   {"Eléctrico"},
		"altura": {"0"},
		"peso":   {"6"},
	}
	request := httptest.NewRequest("POST", "/pokemon", strings.NewReader(form.Encode()))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	response := httptest.NewRecorder()

	handler.CreatePokemon(response, request)

	if response.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", response.Code, http.StatusBadRequest)
	}
	if !strings.Contains(response.Body.String(), "La altura debe ser un número mayor que cero.") {
		t.Errorf("response does not contain the validation error: %q", response.Body.String())
	}
	if store.created.Nombre != "" {
		t.Errorf("invalid Pokémon was saved: %+v", store.created)
	}
}
