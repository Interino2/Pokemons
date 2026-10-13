package repositories

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"strings"
)

//go:embed pokemon_gen1.json
var gen1PokemonData embed.FS

type seededPokemon struct {
	Number int     `json:"number"`
	Nombre string  `json:"name"`
	Tipo   string  `json:"type"`
	Altura float64 `json:"height"`
	Peso   float64 `json:"weight"`
	Imagen string  `json:"image"`
}

func loadGen1Pokemon() ([]seededPokemon, error) {
	data, err := gen1PokemonData.ReadFile("pokemon_gen1.json")
	if err != nil {
		return nil, fmt.Errorf("leyendo la Pokédex de Kanto: %w", err)
	}

	var pokemons []seededPokemon
	if err := json.Unmarshal(data, &pokemons); err != nil {
		return nil, fmt.Errorf("decodificando la Pokédex de Kanto: %w", err)
	}
	if len(pokemons) != 151 {
		return nil, fmt.Errorf("la Pokédex de Kanto contiene %d Pokémon; se esperaban 151", len(pokemons))
	}

	names := make(map[string]struct{}, len(pokemons))
	for index, pokemon := range pokemons {
		if pokemon.Number != index+1 || pokemon.Nombre == "" || pokemon.Tipo == "" ||
			pokemon.Altura <= 0 || pokemon.Peso <= 0 || pokemon.Imagen == "" {
			return nil, fmt.Errorf("datos no válidos para el Pokémon número %d", index+1)
		}
		name := strings.ToLower(pokemon.Nombre)
		if _, exists := names[name]; exists {
			return nil, fmt.Errorf("nombre duplicado en la Pokédex de Kanto: %q", pokemon.Nombre)
		}
		names[name] = struct{}{}
	}
	return pokemons, nil
}

func (r *PokemonRepository) SeedDefaults(ctx context.Context) error {
	pokemons, err := loadGen1Pokemon()
	if err != nil {
		return err
	}

	tx, err := r.DB.Begin(ctx)
	if err != nil {
		return fmt.Errorf("iniciando la carga de la Pokédex de Kanto: %w", err)
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, "SELECT pg_advisory_xact_lock(728194621)"); err != nil {
		return fmt.Errorf("bloqueando la carga de la Pokédex de Kanto: %w", err)
	}
	if _, err := tx.Exec(ctx, `ALTER TABLE pokemons ADD COLUMN IF NOT EXISTS pokedex_number INTEGER`); err != nil {
		return fmt.Errorf("actualizando el esquema de la Pokédex: %w", err)
	}

	names := make([]string, len(pokemons))
	for index, pokemon := range pokemons {
		names[index] = strings.ToLower(pokemon.Nombre)
	}

	if _, err := tx.Exec(
		ctx,
		`DELETE FROM pokemons
		 WHERE id NOT IN (
		   SELECT MIN(id)
		   FROM pokemons
		   WHERE LOWER(nombre) = ANY($1::text[])
		   GROUP BY LOWER(nombre)
		 )`,
		names,
	); err != nil {
		return fmt.Errorf("eliminando Pokémon ajenos a la primera generación: %w", err)
	}

	if _, err := tx.Exec(
		ctx,
		`CREATE UNIQUE INDEX IF NOT EXISTS pokemons_pokedex_number_unique
		 ON pokemons (pokedex_number)
		 WHERE pokedex_number IS NOT NULL`,
	); err != nil {
		return fmt.Errorf("creando el índice de la Pokédex: %w", err)
	}

	for _, pokemon := range pokemons {
		description := fmt.Sprintf("El Pokémon número %03d de la Pokédex de Kanto.", pokemon.Number)
		tag, err := tx.Exec(
			ctx,
			`UPDATE pokemons
			 SET tipo = $2, altura = $3, peso = $4, imagen = $5,
			     descripcion = $6, pokedex_number = $7
			 WHERE LOWER(nombre) = LOWER($1)`,
			pokemon.Nombre,
			pokemon.Tipo,
			pokemon.Altura,
			pokemon.Peso,
			pokemon.Imagen,
			description,
			pokemon.Number,
		)
		if err != nil {
			return fmt.Errorf("actualizando el Pokémon %q: %w", pokemon.Nombre, err)
		}
		if tag.RowsAffected() == 0 {
			_, err := tx.Exec(
				ctx,
				`INSERT INTO pokemons
				 (nombre, tipo, altura, peso, imagen, descripcion, pokedex_number)
				 VALUES ($1, $2, $3, $4, $5, $6, $7)`,
				pokemon.Nombre,
				pokemon.Tipo,
				pokemon.Altura,
				pokemon.Peso,
				pokemon.Imagen,
				description,
				pokemon.Number,
			)
			if err != nil {
				return fmt.Errorf("insertando el Pokémon %q: %w", pokemon.Nombre, err)
			}
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("guardando la Pokédex de Kanto: %w", err)
	}
	return nil
}
