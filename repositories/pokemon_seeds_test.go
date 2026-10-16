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

func TestLoadMegaPokemon(t *testing.T) {
	pokemons, err := loadMegaPokemon()
	if err != nil {
		t.Fatalf("load mega evolutions: %v", err)
	}
	if len(pokemons) != 15 {
		t.Fatalf("got %d mega evolutions, want 15", len(pokemons))
	}

	names := make(map[string]bool, len(pokemons))
	for _, pokemon := range pokemons {
		names[pokemon.Nombre] = true
	}
	for _, name := range []string{
		"Mega Venusaur",
		"Mega Charizard X",
		"Mega Charizard Y",
		"Mega Blastoise",
		"Mega Beedrill",
		"Mega Pidgeot",
		"Mega Alakazam",
		"Mega Slowbro",
		"Mega Gengar",
		"Mega Kangaskhan",
		"Mega Pinsir",
		"Mega Gyarados",
		"Mega Aerodactyl",
		"Mega Mewtwo X",
		"Mega Mewtwo Y",
	} {
		if !names[name] {
			t.Errorf("mega evolution %q not found", name)
		}
	}
}
