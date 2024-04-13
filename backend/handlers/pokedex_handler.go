package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strings"

	"github.com/PranavKurundodi/pokedex/backend/models"
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
