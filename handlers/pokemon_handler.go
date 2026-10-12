package handlers

import (
	"bytes"
	"context"
	"errors"
	"html/template"
	"log"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/starfederation/datastar-go/datastar"

	"pokemon-app/models"
)

type PokemonStore interface {
	GetAll(ctx context.Context) ([]models.Pokemon, error)
	GetByID(ctx context.Context, id int) (models.Pokemon, error)
	Create(ctx context.Context, pokemon models.Pokemon) error
}

type pokemonForm struct {
	Nombre      string
	Tipo        string
	Altura      string
	Peso        string
	Imagen      string
	Descripcion string
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
	h.renderHome(w, r, http.StatusOK, pokemonForm{}, "")
}

func (h *PokemonHandler) CreatePokemon(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	if err := r.ParseForm(); err != nil {
		h.renderHome(w, r, http.StatusBadRequest, pokemonForm{}, "No se pudo leer el formulario. Inténtalo de nuevo.")
		return
	}

	form := pokemonForm{
		Nombre:      strings.TrimSpace(r.FormValue("nombre")),
		Tipo:        strings.TrimSpace(r.FormValue("tipo")),
		Altura:      strings.TrimSpace(r.FormValue("altura")),
		Peso:        strings.TrimSpace(r.FormValue("peso")),
		Imagen:      strings.TrimSpace(r.FormValue("imagen")),
		Descripcion: strings.TrimSpace(r.FormValue("descripcion")),
	}
	pokemon, err := pokemonFromForm(form)
	if err != nil {
		h.renderHome(w, r, http.StatusBadRequest, form, err.Error())
		return
	}

	if err := h.Repository.Create(r.Context(), pokemon); err != nil {
		log.Printf("guardando Pokémon %q: %v", pokemon.Nombre, err)
		h.renderHome(w, r, http.StatusInternalServerError, form, "No se pudo guardar el Pokémon. Inténtalo de nuevo.")
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *PokemonHandler) renderHome(w http.ResponseWriter, r *http.Request, status int, form pokemonForm, formError string) {
	pokemons, err := h.Repository.GetAll(r.Context())
	if err != nil {
		log.Printf("cargando la colección Pokémon: %v", err)
		http.Error(w, "No se pudieron cargar los Pokémon.", http.StatusInternalServerError)
		return
	}

	data := struct {
		Pokemons []models.Pokemon
		Form     pokemonForm
		Error    string
	}{
		Pokemons: pokemons,
		Form:     form,
		Error:    formError,
	}

	var rendered bytes.Buffer
	if err := h.Template.ExecuteTemplate(&rendered, "index.html", data); err != nil {
		log.Printf("renderizando la página principal: %v", err)
		http.Error(w, "No se pudo mostrar la página.", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	if _, err := rendered.WriteTo(w); err != nil {
		log.Printf("escribiendo la página principal: %v", err)
	}
}

func pokemonFromForm(form pokemonForm) (models.Pokemon, error) {
	if form.Nombre == "" || len([]rune(form.Nombre)) > 100 {
		return models.Pokemon{}, errors.New("El nombre es obligatorio y debe tener como máximo 100 caracteres.")
	}
	if form.Tipo == "" || len([]rune(form.Tipo)) > 50 {
		return models.Pokemon{}, errors.New("El tipo es obligatorio y debe tener como máximo 50 caracteres.")
	}

	altura, err := strconv.ParseFloat(form.Altura, 64)
	if err != nil || math.IsNaN(altura) || math.IsInf(altura, 0) || altura <= 0 {
		return models.Pokemon{}, errors.New("La altura debe ser un número mayor que cero.")
	}
	peso, err := strconv.ParseFloat(form.Peso, 64)
	if err != nil || math.IsNaN(peso) || math.IsInf(peso, 0) || peso <= 0 {
		return models.Pokemon{}, errors.New("El peso debe ser un número mayor que cero.")
	}
	if len(form.Descripcion) > 1000 {
		return models.Pokemon{}, errors.New("La descripción debe tener como máximo 1000 caracteres.")
	}
	if form.Imagen != "" {
		imageURL, err := url.ParseRequestURI(form.Imagen)
		if err != nil || (imageURL.Scheme != "http" && imageURL.Scheme != "https") || imageURL.Host == "" {
			return models.Pokemon{}, errors.New("La imagen debe ser una URL válida que empiece con http:// o https://.")
		}
	}

	return models.Pokemon{
		Nombre:      form.Nombre,
		Tipo:        form.Tipo,
		Altura:      altura,
		Peso:        peso,
		Imagen:      form.Imagen,
		Descripcion: form.Descripcion,
	}, nil
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
