package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

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
	pokedexRouter.HandleFunc("/display/{username}", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetUserInventory(w, r, client)
	}).Methods("GET")
	r.HandleFunc("/upload", uploadHandler).Methods("POST")

	log.Println("Server started on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
func uploadHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the multipart form data
	err := r.ParseMultipartForm(10 << 20) // 10 MB limit
	if err != nil {
		http.Error(w, fmt.Sprintf("Error: %s", err.Error()), http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("image")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error: %s", err.Error()), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Create a new file named "model/test.jpg"
	out, err := os.Create("C:\\Users\\prana\\OneDrive\\Documents\\college\\pokedex\\model\\test.jpg")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	defer out.Close()

	// Copy the uploaded file to the new file
	_, err = io.Copy(out, file)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("File uploaded successfully"))
}

