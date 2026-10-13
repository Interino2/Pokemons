package repositories

import "testing"

func TestLoadGen1Pokemon(t *testing.T) {
	pokemons, err := loadGen1Pokemon()
	if err != nil {
		t.Fatalf("load generation 1 Pokémon: %v", err)
	}
	if len(pokemons) != 151 {
		t.Fatalf("got %d Pokémon, want 151", len(pokemons))
	}

	for index, pokemon := range pokemons {
		if pokemon.Number != index+1 {
			t.Errorf("Pokémon at index %d has Pokédex number %d", index, pokemon.Number)
		}
	}

	for number, want := range map[int]string{
		1:   "Bulbasaur",
		25:  "Pikachu",
		29:  "Nidoran♀",
		32:  "Nidoran♂",
		122: "Mr. Mime",
		151: "Mew",
	} {
		if got := pokemons[number-1].Nombre; got != want {
			t.Errorf("Pokémon #%03d = %q, want %q", number, got, want)
		}
	}
}
