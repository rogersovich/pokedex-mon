package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type ResourceReference struct {
	Name string `json:"name" bson:"name"`
	URL  string `json:"url" bson:"url"`
}

type BaseGenerationNames struct {
	Name     string `json:"name"`
	Language struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"language"`
}

type GenerationDetail struct {
	ID             int                   `json:"id"`
	Name           string                `json:"name"`
	Names          []BaseGenerationNames `json:"names"`
	Abilities      []ResourceReference   `json:"abilities"`
	MainRegion     ResourceReference     `json:"main_region"`
	Moves          []ResourceReference   `json:"moves"`
	PokemonSpecies []ResourceReference   `json:"pokemon_species"`
	Types          []ResourceReference   `json:"types"`
	VersionGroups  []ResourceReference   `json:"version_groups"`
}

type ListGeneration struct {
	Count    int                  `json:"count" bson:"count"`
	Next     *string              `json:"next" bson:"next"`
	Previous *string              `json:"previous" bson:"previous"`
	Results  []ListGenerationItem `json:"results" bson:"results"`
}

type ListGenerationItem struct {
	ID   int    `json:"id" bson:"id"`
	Name string `json:"name" bson:"name"`
	URL  string `json:"url" bson:"url"`
}

type GenerationSaveDocument struct {
	ID             primitive.ObjectID    `bson:"_id,omitempty"`
	GenerationID   int                   `json:"id" bson:"id"`
	Name           string                `json:"name" bson:"name"`
	Names          []BaseGenerationNames `json:"names" bson:"names"`
	Abilities      []ResourceReference   `json:"abilities" bson:"abilities"`
	MainRegion     ResourceReference     `json:"main_region" bson:"main_region"`
	Moves          []ResourceReference   `json:"moves" bson:"moves"`
	PokemonSpecies []ResourceReference   `json:"pokemon_species" bson:"pokemon_species"`
	Types          []ResourceReference   `json:"types" bson:"types"`
	VersionGroups  []ResourceReference   `json:"version_groups" bson:"version_groups"`
	LastSyncedAt   int64                 `json:"-" bson:"last_synced_at,omitempty"`
}

type GenerationDocument struct {
	ID             primitive.ObjectID    `bson:"_id,omitempty"`
	GenerationID   int                   `json:"id" bson:"id"`
	Name           string                `json:"name" bson:"name"`
	Names          []BaseGenerationNames `json:"names" bson:"names"`
	Abilities      []ResourceReference   `json:"abilities" bson:"abilities"`
	MainRegion     ResourceReference     `json:"main_region" bson:"main_region"`
	Moves          []ResourceReference   `json:"moves" bson:"moves"`
	PokemonSpecies []ResourceReference   `json:"pokemon_species" bson:"pokemon_species"`
	Types          []ResourceReference   `json:"types" bson:"types"`
	VersionGroups  []ResourceReference   `json:"version_groups" bson:"version_groups"`
	LastSyncedAt   int64                 `json:"-" bson:"last_synced_at,omitempty"`
}

type GenerationListDocument struct {
	ID   int    `json:"id" bson:"id"`
	Name string `json:"name" bson:"name"`
}
