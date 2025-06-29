package model

import "go.mongodb.org/mongo-driver/bson/primitive"

// AbilityListItem
type AbilityListItem struct {
	Name string `json:"name" bson:"name"`
	URL  string `json:"url" bson:"url"`
}

// AbilityListResponse
type AbilityListResponse struct {
	Count    int               `json:"count" bson:"count"`
	Next     *string           `json:"next" bson:"next"`
	Previous *string           `json:"previous" bson:"previous"`
	Results  []AbilityListItem `json:"results" bson:"results"`
}

type EffectEntries struct {
	Effect      string `json:"effect" bson:"effect"`
	ShortEffect string `json:"short_effect" bson:"short_effect"`
	Language    struct {
		Name string `json:"name" bson:"name"`
		URL  string `json:"url" bson:"url"`
	} `json:"language" bson:"language"`
}

type FlavorTextEntries struct {
	FlavorText string `json:"flavor_text" bson:"flavor_text"`
	Language   struct {
		Name string `json:"name" bson:"name"`
		URL  string `json:"url" bson:"url"`
	} `json:"language" bson:"language"`
	VersionGroup struct {
		Name string `json:"name" bson:"name"`
		URL  string `json:"url" bson:"url"`
	} `json:"version_group" bson:"version_group"`
}

type Generation struct {
	Name string `json:"name" bson:"name"`
	URL  string `json:"url" bson:"url"`
}

type Pokemon struct {
	Name string `json:"name" bson:"name"`
	URL  string `json:"url" bson:"url"`
}

type PokemonEntry struct {
	IsHidden bool    `json:"is_hidden" bson:"is_hidden"`
	Pokemon  Pokemon `json:"pokemon" bson:"pokemon"`
	Slot     int     `json:"slot" bson:"slot"`
}

type NameEntry struct {
	Language struct {
		Name string `json:"name" bson:"name"`
		URL  string `json:"url" bson:"url"`
	} `json:"language" bson:"language"`
	Name string `json:"name" bson:"name"`
}

// AbilityDetai
type AbilityDetail struct {
	ID                int                 `json:"id" bson:"id"`
	EffectChanges     []interface{}       `json:"effect_changes" bson:"effect_changes"`
	EffectEntries     []EffectEntries     `json:"effect_entries" bson:"effect_entries"`
	FlavorTextEntries []FlavorTextEntries `json:"flavor_text_entries" bson:"flavor_text_entries"`
	Generation        Generation          `json:"generation" bson:"generation"`
	IsMainSeries      bool                `json:"is_main_series" bson:"is_main_series"`
	Name              string              `json:"name" bson:"name"`
	Names             []NameEntry         `json:"names" bson:"names"`
	Pokemon           []PokemonEntry      `json:"pokemon" bson:"pokemon"`
}

// AbilityDocument is the structure to store in MongoDB
type AbilityDocument struct {
	ID                primitive.ObjectID  `bson:"_id,omitempty"`
	AbilityID         int                 `bson:"id"`
	EffectChanges     []interface{}       `bson:"effect_changes"`
	EffectEntries     []EffectEntries     `bson:"effect_entries"`
	FlavorTextEntries []FlavorTextEntries `bson:"flavor_text_entries"`
	Generation        Generation          `bson:"generation"`
	IsMainSeries      bool                `bson:"is_main_series"`
	Name              string              `bson:"name"`
	Names             []NameEntry         `bson:"names"`
	Pokemon           []PokemonEntry      `bson:"pokemon"`
	LastSyncedAt      int64               `bson:"last_synced_at"`
}
