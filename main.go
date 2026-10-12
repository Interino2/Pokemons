package main

import (
	"html/template"
	"log"
	"net/http"
	"time"

	"pokemon-app/config"
	"pokemon-app/database"
	"pokemon-app/handlers"
	"pokemon-app/repositories"
)
//comentario
func main() {

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	db, err := database.Connect(cfg)

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	pokemonRepository := repositories.NewPokemonRepository(db)

	tmpl := template.Must(
		template.ParseFiles("templates/index.html"),
	)

	pokemonHandler := handlers.NewPokemonHandler(
		pokemonRepository,
		tmpl,
	)

	mux := http.NewServeMux()
	mux.Handle(
		"/static/",
		http.StripPrefix(
			"/static/",
			http.FileServer(http.Dir("static")),
		),
	)

	mux.HandleFunc("GET /{$}", pokemonHandler.Home)
	mux.HandleFunc("POST /pokemon", pokemonHandler.CreatePokemon)
	mux.HandleFunc("GET /pokemon/{id}", pokemonHandler.GetPokemon)

	log.Printf(
		"Servidor iniciado en http://localhost:%s",
		cfg.AppPort,
	)

	server := &http.Server{
		Addr:              ":" + cfg.AppPort,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}
	log.Fatal(server.ListenAndServe())
}
