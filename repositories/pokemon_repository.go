package repositories

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"pokemon-app/models"
)

type PokemonRepository struct {
	DB *pgxpool.Pool
}

func NewPokemonRepository(db *pgxpool.Pool) *PokemonRepository {
	return &PokemonRepository{
		DB: db,
	}
}

func (r *PokemonRepository) GetAll(ctx context.Context) ([]models.Pokemon, error) {
	rows, err := r.DB.Query(
		ctx,
		`SELECT id, COALESCE(pokedex_number, id), nombre, tipo, altura, peso, imagen, descripcion, es_megaevolucion
		 FROM pokemons
		 ORDER BY pokedex_number NULLS LAST, es_megaevolucion, id`,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var pokemons []models.Pokemon

	for rows.Next() {
		var pokemon models.Pokemon

		err := rows.Scan(
			&pokemon.ID,
			&pokemon.PokedexNumber,
			&pokemon.Nombre,
			&pokemon.Tipo,
			&pokemon.Altura,
			&pokemon.Peso,
			&pokemon.Imagen,
			&pokemon.Descripcion,
			&pokemon.EsMegaEvolucion,
		)

		if err != nil {
			return nil, err
		}

		pokemons = append(pokemons, pokemon)
	}

	return pokemons, rows.Err()
}
func (r *PokemonRepository) GetByID(ctx context.Context, id int) (models.Pokemon, error) {

	var pokemon models.Pokemon

	err := r.DB.QueryRow(
		ctx,
		`SELECT id, COALESCE(pokedex_number, id), nombre, tipo, altura, peso, imagen, descripcion, es_megaevolucion
		 FROM pokemons
		 WHERE id = $1`,
		id,
	).Scan(
		&pokemon.ID,
		&pokemon.PokedexNumber,
		&pokemon.Nombre,
		&pokemon.Tipo,
		&pokemon.Altura,
		&pokemon.Peso,
		&pokemon.Imagen,
		&pokemon.Descripcion,
		&pokemon.EsMegaEvolucion,
	)

	return pokemon, err
}

func (r *PokemonRepository) Create(ctx context.Context, pokemon models.Pokemon) error {
	_, err := r.DB.Exec(
		ctx,
		`INSERT INTO pokemons (nombre, tipo, altura, peso, imagen, descripcion)
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		pokemon.Nombre,
		pokemon.Tipo,
		pokemon.Altura,
		pokemon.Peso,
		pokemon.Imagen,
		pokemon.Descripcion,
	)
	return err
}
