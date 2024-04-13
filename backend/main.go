package main

import (
	"log"
	"net/http"

	"github.com/PranavKurundodi/pokedex/backend/handlers"
	"github.com/PranavKurundodi/pokedex/backend/models"
	"github.com/gorilla/mux"
)

func main() {
	pokedex, err := models.ReadPokemonFromCSV("pokemon.csv")
	if err != nil {
		log.Fatalf("Failed to read pokemon data from CSV: %v", err)
	}

	handlers.SetPokedex(pokedex)

	r := mux.NewRouter()

	r.HandleFunc("/pokemon", handlers.GetPokemon).Methods("GET")
	r.HandleFunc("/pokemon/byName", handlers.GetPokemonByName).Methods("GET")
	r.HandleFunc("/pokemon/model", handlers.PokemonProbModel).Methods("GET")

	log.Println("Server started on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
