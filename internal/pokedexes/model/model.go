package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type BaseResourceReference struct {
	Name string `json:"name" bson:"name"`
	URL  string `json:"url" bson:"url"`
}

type BasePokedexesNames struct {
	Name     string `json:"name"`
	Language struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"language"`
}

type BasePokedexesDescription struct {
	Description string `json:"description"`
	Language    struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"language"`
}

type BasePokemonEntries struct {
	EnteryNumber   int                   `json:"entry_number"`
	PokemonSpecies BaseResourceReference `json:"pokemon_species"`
}

type PokedexesDetail struct {
	ID             int                        `json:"id"`
	IsMainSeries   bool                       `json:"is_main_series"`
	Region         BaseResourceReference      `json:"region"`
	VersionGroup   []BaseResourceReference    `json:"version_group"`
	Name           string                     `json:"name"`
	Names          []BasePokedexesNames       `json:"names"`
	Descriptions   []BasePokedexesDescription `json:"descriptions"`
	PokemonEntries []BasePokemonEntries       `json:"pokemon_entries"`
}

type ListPokedexes struct {
	Count    int                 `json:"count" bson:"count"`
	Next     *string             `json:"next" bson:"next"`
	Previous *string             `json:"previous" bson:"previous"`
	Results  []ListPokedexesItem `json:"results" bson:"results"`
}

type ListPokedexesItem struct {
	ID   int    `json:"id" bson:"id"`
	Name string `json:"name" bson:"name"`
	URL  string `json:"url" bson:"url"`
}

type PokedexesSaveDocument struct {
	ID             primitive.ObjectID         `bson:"_id,omitempty"`
	PokedexesID    int                        `json:"id"`
	IsMainSeries   bool                       `json:"is_main_series"`
	Region         BaseResourceReference      `json:"region"`
	VersionGroup   []BaseResourceReference    `json:"version_group"`
	Name           string                     `json:"name"`
	Names          []BasePokedexesNames       `json:"names"`
	Descriptions   []BasePokedexesDescription `json:"descriptions"`
	PokemonEntries []BasePokemonEntries       `json:"pokemon_entries"`
	LastSyncedAt   int64                      `json:"-" bson:"last_synced_at,omitempty"`
}

type PokedexesDocument struct {
	ID             primitive.ObjectID         `bson:"_id,omitempty"`
	PokedexesID    int                        `json:"id"`
	IsMainSeries   bool                       `json:"is_main_series"`
	Region         BaseResourceReference      `json:"region"`
	VersionGroup   []BaseResourceReference    `json:"version_group"`
	Name           string                     `json:"name"`
	Names          []BasePokedexesNames       `json:"names"`
	Descriptions   []BasePokedexesDescription `json:"descriptions"`
	PokemonEntries []BasePokemonEntries       `json:"pokemon_entries"`
	LastSyncedAt   int64                      `json:"-" bson:"last_synced_at,omitempty"`
}

type PokedexesListDocument struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
