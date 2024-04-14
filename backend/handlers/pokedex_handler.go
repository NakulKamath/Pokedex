package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strings"

	"github.com/PranavKurundodi/pokedex/backend/models"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var pokedex []models.Pokemon

func SetPokedex(pd []models.Pokemon) {
	pokedex = pd
}

func GetPokemon(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(pokedex)
}

func GetPokemonByName(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	name := queryParams.Get("name")

	var foundPokemon *models.Pokemon

	for _, pokemon := range pokedex {
		if pokemon.Name == name {
			foundPokemon = &models.Pokemon{
				Name:      pokemon.Name,
				Type1:     pokemon.Type1,
				Type2:     pokemon.Type2,
				Evolution: pokemon.Evolution,
			}

			break
		}
	}

	if foundPokemon == nil {
		http.NotFound(w, r)
		return
	}

	json.NewEncoder(w).Encode(foundPokemon)
}
func revstr(a string) (r string) {
	runes := []rune(a)

	// Reverse the slice
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}

	// Convert the slice of runes back to a string
	r = string(runes)
	return r
}

func PokemonProbModel(w http.ResponseWriter, r *http.Request) {
	cmd := exec.Command("python", "C:\\Users\\prana\\OneDrive\\Documents\\college\\pokedex\\model\\test.py")

	// Create a buffer to store the output
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Run the command and wait for it to finish
	err := cmd.Run()
	if err != nil {
		log.Fatalf("Failed to run Python script: %v\nStderr: %s", err, stderr.String())
	}

	// Get the output of the Python script
	pythonOutput := stdout.String()

	// Convert the slice of runes back to a string
	reversedStr := revstr(pythonOutput)

	index := strings.Index(reversedStr, "pets")
	if index == -1 {
		// "pets" not found
		fmt.Println("String 'pets' not found")
		return
	}

	// Extract the substring before the index of "pets"
	substr := reversedStr[:index]

	toppokemon := revstr(substr)
	toppokemon = strings.ToLower(toppokemon)
	toppokemon = strings.TrimSpace(strings.ReplaceAll(toppokemon, "\n", ""))

	// Print the substring
	var toppokemonStats models.Pokemon
	for _, pokemon := range pokedex {
		if toppokemon == pokemon.Name {
			toppokemonStats = pokemon
		}
	}
	json.NewEncoder(w).Encode(toppokemonStats)

}

var ctx = context.TODO()

func NewMongoClient() (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func RegisterUser(w http.ResponseWriter, r *http.Request, client *mongo.Client) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check if user already exists
	collection := client.Database("pokedex").Collection("users")
	var existingUser models.User
	err = collection.FindOne(ctx, bson.M{"username": user.Username}).Decode(&existingUser)
	if err == nil {
		http.Error(w, "Username already exists", http.StatusBadRequest)
		return
	}

	// If inventory is nil, initialize it as an empty array
	if user.Inventory == nil {
		user.Inventory = []models.Pokemon{}
	}

	// Insert new user with inventory
	_, err = collection.InsertOne(ctx, user)
	if err != nil {
		http.Error(w, "Failed to register user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func LoginUser(w http.ResponseWriter, r *http.Request, client *mongo.Client) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check if user exists and password matches
	collection := client.Database("pokedex").Collection("users")
	var existingUser models.User
	err = collection.FindOne(ctx, bson.M{"username": user.Username}).Decode(&existingUser)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	if existingUser.Password != user.Password {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Generate and set JWT token
	// Here you can use libraries like jwt-go to generate JWT tokens

	w.WriteHeader(http.StatusOK)
}

func AddToInventory(w http.ResponseWriter, r *http.Request, client *mongo.Client) {
	vars := mux.Vars(r)
	username := vars["username"] // Get username from URL path
	var pokemon models.Pokemon

	err := json.NewDecoder(r.Body).Decode(&pokemon)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Access the "users" collection
	collection := client.Database("pokedex").Collection("users")

	// Update the user's inventory
	_, err = collection.UpdateOne(context.Background(), bson.M{"username": username}, bson.M{"$addToSet": bson.M{"inventory": pokemon}})
	if err != nil {
		http.Error(w, "Failed to add Pokemon to inventory", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func RemoveFromInventory(w http.ResponseWriter, r *http.Request, client *mongo.Client) {
	vars := mux.Vars(r)
	username := vars["username"] // Get username from URL path
	var pokemon models.Pokemon

	err := json.NewDecoder(r.Body).Decode(&pokemon)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Access the "users" collection
	collection := client.Database("pokedex").Collection("users")

	// Update the user's inventory
	_, err = collection.UpdateOne(context.Background(), bson.M{"username": username}, bson.M{"$pull": bson.M{"inventory": pokemon}})
	if err != nil {
		http.Error(w, "Failed to remove Pokemon from inventory", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check authentication token
		// Code to verify JWT token goes here

		// If token is valid, proceed to next middleware
		next.ServeHTTP(w, r)
	})
}
