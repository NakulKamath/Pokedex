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

	client, err := handlers.NewMongoClient()
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}

	r.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		handlers.RegisterUser(w, r, client)
	}).Methods("POST")

	r.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		handlers.LoginUser(w, r, client)
	}).Methods("POST")

	pokedexRouter := r.PathPrefix("/pokedex").Subrouter()
	pokedexRouter.Use(handlers.AuthMiddleware)
	pokedexRouter.HandleFunc("/add/{username}", func(w http.ResponseWriter, r *http.Request) {
		handlers.AddToInventory(w, r, client)
	}).Methods("POST")
	pokedexRouter.HandleFunc("/remove/{username}", func(w http.ResponseWriter, r *http.Request) {
		handlers.RemoveFromInventory(w, r, client)
	}).Methods("DELETE")

	log.Println("Server started on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
