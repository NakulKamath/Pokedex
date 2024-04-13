package models

import (
	"encoding/csv"
	"os"
)

type Pokemon struct {
	Name      string `json:"name"`
	Type1     string `json:"type1"`
	Type2     string `json:"type2"`
	Evolution string `json:"evolution"`
}
type User struct {
	Username  string    `json:"username" bson:"username"`
	Password  string    `json:"password" bson:"password"`
	Inventory []Pokemon `json:"inventory" bson:"inventory"` // Array of Pokémon IDs
}

func ReadPokemonFromCSV(filename string) ([]Pokemon, error) {
	var pokedex []Pokemon

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	reader := csv.NewReader(file)

	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	for _, record := range records {
		pokemon := Pokemon{
			Name:      record[0],
			Type1:     record[1],
			Type2:     record[2],
			Evolution: record[3],
		}
		pokedex = append(pokedex, pokemon)
	}

	return pokedex, nil
}
