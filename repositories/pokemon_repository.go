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
		`SELECT id, nombre, tipo, altura, peso, imagen, descripcion
		 FROM pokemons
		 ORDER BY CASE nombre
    WHEN 'Bulbasaur' THEN 1
    WHEN 'Ivysaur' THEN 2
    WHEN 'Venusaur' THEN 3
    WHEN 'Charmander' THEN 4
    WHEN 'Charmeleon' THEN 5
    WHEN 'Charizard' THEN 6
    WHEN 'Squirtle' THEN 7
    WHEN 'Wartortle' THEN 8
    WHEN 'Blastoise' THEN 9
    WHEN 'Pikachu' THEN 10
    WHEN 'Raichu' THEN 11
    WHEN 'Gastly' THEN 12
    WHEN 'Haunter' THEN 13
    WHEN 'Gengar' THEN 14
    WHEN 'Eevee' THEN 15
    WHEN 'Vaporeon' THEN 16
    WHEN 'Jolteon' THEN 17
    WHEN 'Flareon' THEN 18
    WHEN 'Espeon' THEN 19
    WHEN 'Umbreon' THEN 20
    WHEN 'Leafeon' THEN 21
    WHEN 'Glaceon' THEN 22
    WHEN 'Sylveon' THEN 23
END`,
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
			&pokemon.Nombre,
			&pokemon.Tipo,
			&pokemon.Altura,
			&pokemon.Peso,
			&pokemon.Imagen,
			&pokemon.Descripcion,
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
		`SELECT id, nombre, tipo, altura, peso, imagen, descripcion
		 FROM pokemons
		 WHERE id = $1`,
		id,
	).Scan(
		&pokemon.ID,
		&pokemon.Nombre,
		&pokemon.Tipo,
		&pokemon.Altura,
		&pokemon.Peso,
		&pokemon.Imagen,
		&pokemon.Descripcion,
	)

	return pokemon, err
}
