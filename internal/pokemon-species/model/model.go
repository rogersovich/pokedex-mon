package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type PokemonSpeciesListItem struct {
	Name string `json:"name" bson:"name"`
	URL  string `json:"url" bson:"url"`
}

type PokemonSpeciesListResponse struct {
	Count    int                      `json:"count" bson:"count"`
	Next     *string                  `json:"next" bson:"next"`
	Previous *string                  `json:"previous" bson:"previous"`
	Results  []PokemonSpeciesListItem `json:"results" bson:"results"`
}

type PokemonSpeciesDetail struct {
	PokeAPIID            int                 `json:"id"`
	Name                 string              `json:"name"`
	BaseHappiness        int                 `json:"base_happiness"`
	CaptureRate          int                 `json:"capture_rate"`
	Color                ResourceReference   `json:"color"`
	EggGroups            []ResourceReference `json:"egg_groups"`
	EvolutionChain       ResourceReference   `json:"evolution_chain"`
	EvolvesFromSpecies   *ResourceReference  `json:"evolves_from_species,omitempty"` // Use pointer for null
	FlavorTextEntries    []interface{}       `json:"flavor_text_entries"`
	FormDescriptions     []interface{}       `json:"form_descriptions"`
	FormsSwitchable      bool                `json:"forms_switchable"`
	GenderRate           int                 `json:"gender_rate"`
	Genera               []interface{}       `json:"genera"`
	Generation           ResourceReference   `json:"generation"`
	GrowthRate           ResourceReference   `json:"growth_rate"`
	Habitat              ResourceReference   `json:"habitat"`
	HasGenderDifferences bool                `json:"has_gender_differences"`
	HatchCounter         int                 `json:"hatch_counter"`
	IsBaby               bool                `json:"is_baby"`
	IsLegendary          bool                `json:"is_legendary"`
	IsMythical           bool                `json:"is_mythical"`
	Names                []PokemonNames      `json:"names"`
	Order                int                 `json:"order"`
	PalParkEncounters    []PalParkEncounter  `json:"pal_park_encounters"`
	PokedexNumbers       []interface{}       `json:"pokedex_numbers"`
	Shape                ResourceReference   `json:"shape"`
	Varieties            []PokemonVariety    `json:"varieties"`
}

// --- POKEMON DOCUMENT MODEL ---
type PokemonSpeciesDocument struct {
	ID                   primitive.ObjectID  `bson:"_id,omitempty"` // MongoDB's primary key
	PokeAPIID            int                 `json:"id" bson:"pokeapi_id"`
	Name                 string              `json:"name" bson:"name"`
	BaseHappiness        int                 `json:"base_happiness" bson:"base_happiness"`
	CaptureRate          int                 `json:"capture_rate" bson:"capture_rate"`
	Color                ResourceReference   `json:"color" bson:"color"`
	EggGroups            []ResourceReference `json:"egg_groups" bson:"egg_groups"`
	EvolutionChain       ResourceReference   `json:"evolution_chain" bson:"evolution_chain"`
	EvolvesFromSpecies   *ResourceReference  `json:"evolves_from_species" bson:"evolves_from_species,omitempty"` // Use pointer for null
	FlavorTextEntries    []interface{}       `json:"flavor_text_entries" bson:"flavor_text_entries"`
	FormDescriptions     []interface{}       `json:"form_descriptions" bson:"form_descriptions"`
	FormsSwitchable      bool                `json:"forms_switchable" bson:"forms_switchable"`
	GenderRate           int                 `json:"gender_rate" bson:"gender_rate"`
	Genera               []interface{}       `json:"genera" bson:"genera"`
	Generation           ResourceReference   `json:"generation" bson:"generation"`
	GrowthRate           ResourceReference   `json:"growth_rate" bson:"growth_rate"`
	Habitat              ResourceReference   `json:"habitat" bson:"habitat"`
	HasGenderDifferences bool                `json:"has_gender_differences" bson:"has_gender_differences"`
	HatchCounter         int                 `json:"hatch_counter" bson:"hatch_counter"`
	IsBaby               bool                `json:"is_baby" bson:"is_baby"`
	IsLegendary          bool                `json:"is_legendary" bson:"is_legendary"`
	IsMythical           bool                `json:"is_mythical" bson:"is_mythical"`
	Names                []PokemonNames      `json:"names" bson:"names"`
	Order                int                 `json:"order" bson:"order"`
	PalParkEncounters    []PalParkEncounter  `json:"pal_park_encounters" bson:"pal_park_encounters"`
	PokedexNumbers       []interface{}       `json:"pokedex_numbers" bson:"pokedex_numbers"`
	Shape                ResourceReference   `json:"shape" bson:"shape"`
	Varieties            []PokemonVariety    `json:"varieties" bson:"varieties"`
	LastSyncedAt         int64               `json:"-" bson:"last_synced_at,omitempty"`
}

// ResourceReference represents a generic name and URL reference
type ResourceReference struct {
	Name string `json:"name" bson:"name"`
	URL  string `json:"url" bson:"url"`
}

// PalParkEncounter represents an encounter in Pal Park
type PalParkEncounter struct {
	Area      ResourceReference `json:"area" bson:"area"`
	BaseScore int               `json:"base_score" bson:"base_score"`
	Rate      int               `json:"rate" bson:"rate"`
}

// PokemonVariety represents a specific variety of a Pokemon
type PokemonVariety struct {
	IsDefault bool              `json:"is_default" bson:"is_default"`
	Pokemon   ResourceReference `json:"pokemon" bson:"pokemon"`
}

type PokemonNames struct {
	Name     string            `json:"name" bson:"name"`
	Language ResourceReference `json:"language" bson:"language"`
}
