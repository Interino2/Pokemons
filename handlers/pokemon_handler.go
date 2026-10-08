package handlers

import (
	"bytes"
	"context"
	"errors"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/starfederation/datastar-go/datastar"

	"pokemon-app/models"
)

type PokemonStore interface {
	GetAll(ctx context.Context) ([]models.Pokemon, error)
	GetByID(ctx context.Context, id int) (models.Pokemon, error)
}

type PokemonHandler struct {
	Repository PokemonStore
	Template   *template.Template
}

func NewPokemonHandler(repository PokemonStore, tmpl *template.Template) *PokemonHandler {
	return &PokemonHandler{
		Repository: repository,
		Template:   tmpl,
	}
}

func (h *PokemonHandler) Home(w http.ResponseWriter, r *http.Request) {
	pokemons, err := h.Repository.GetAll(r.Context())
	if err != nil {
		log.Printf("cargando la colección Pokémon: %v", err)
		http.Error(w, "No se pudieron cargar los Pokémon.", http.StatusInternalServerError)
		return
	}

	data := struct {
		Pokemons []models.Pokemon
	}{
		Pokemons: pokemons,
	}

	var rendered bytes.Buffer
	if err := h.Template.ExecuteTemplate(&rendered, "index.html", data); err != nil {
		log.Printf("renderizando la página principal: %v", err)
		http.Error(w, "No se pudo mostrar la página.", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if _, err := rendered.WriteTo(w); err != nil {
		log.Printf("escribiendo la página principal: %v", err)
	}
}

func (h *PokemonHandler) GetPokemon(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	pokemon, err := h.Repository.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.NotFound(w, r)
			return
		}
		log.Printf("cargando Pokémon con ID %d: %v", id, err)
		http.Error(w, "No se pudo cargar el Pokémon.", http.StatusInternalServerError)
		return
	}

	var rendered bytes.Buffer
	if err := h.Template.ExecuteTemplate(&rendered, "pokemon-detail", pokemon); err != nil {
		log.Printf("renderizando Pokémon con ID %d: %v", id, err)
		http.Error(w, "No se pudo mostrar el Pokémon.", http.StatusInternalServerError)
		return
	}

	sse := datastar.NewSSE(w, r)
	if err := sse.PatchElements(rendered.String(), datastar.WithSelector("#pokemon-info")); err != nil {
		log.Printf("actualizando la ficha del Pokémon con ID %d: %v", id, err)
	}
}
