package main

import (
	"log"
	"net/http"

	"github.com/gorilla/handlers"

	"github.com/PranavKurundodi/pokedex/backend/handler"
	"github.com/PranavKurundodi/pokedex/backend/models"
	"github.com/gorilla/mux"
)

func main() {
	pokedex, err := models.ReadPokemonFromCSV("pokemon.csv")
	if err != nil {
		log.Fatalf("Failed to read pokemon data from CSV: %v", err)
	}

	handler.SetPokedex(pokedex)

	r := mux.NewRouter()

	r.HandleFunc("/pokemon", handler.GetPokemon).Methods("GET")
	r.HandleFunc("/pokemon/byName", handler.GetPokemonByName).Methods("GET")
	r.HandleFunc("/pokemon/model", handler.PokemonProbModel).Methods("GET")

	client, err := handler.NewMongoClient()
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}

	r.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		handler.RegisterUser(w, r, client)
	}).Methods("POST")

	r.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		handler.LoginUser(w, r, client)
	}).Methods("POST")

	pokedexRouter := r.PathPrefix("/pokedex").Subrouter()
	pokedexRouter.Use(handler.AuthMiddleware)
	pokedexRouter.HandleFunc("/add/{username}", func(w http.ResponseWriter, r *http.Request) {
		handler.AddToInventory(w, r, client)
	}).Methods("POST")
	pokedexRouter.HandleFunc("/remove/{username}", func(w http.ResponseWriter, r *http.Request) {
		handler.RemoveFromInventory(w, r, client)
	}).Methods("DELETE")
	pokedexRouter.HandleFunc("/display/{username}", func(w http.ResponseWriter, r *http.Request) {
		handler.GetUserInventory(w, r, client)
	}).Methods("GET")
	r.HandleFunc("/upload", handler.UploadHandler).Methods("POST")

	headers := handlers.AllowedHeaders([]string{"Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"})
	origins := handlers.AllowedOrigins([]string{"http://localhost:8000"})
	r.Use(handlers.CORS(headers, methods, origins))

	log.Println("Server started on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
