package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type ResourceReference struct {
	Name string `json:"name" bson:"name"`
	URL  string `json:"url" bson:"url"`
}

type TypeDamageRelations struct {
	DoubleDamgeFrom []ResourceReference `json:"double_damage_from" bson:"double_damage_from"`
	DoubleDamgeTo   []ResourceReference `json:"double_damage_to" bson:"double_damage_to"`
	HalfDamgeFrom   []ResourceReference `json:"half_damage_from" bson:"half_damage_from"`
	HalfDamgeTo     []ResourceReference `json:"half_damage_to" bson:"half_damage_to"`
	NoDamgeFrom     []ResourceReference `json:"no_damage_from" bson:"no_damage_from"`
	NoDamgeTo       []ResourceReference `json:"no_damage_to" bson:"no_damage_to"`
}

type PokemonTypeDocument struct {
	ID              primitive.ObjectID  `bson:"_id,omitempty"`
	TypeID          int                 `json:"id" bson:"id"`
	Name            string              `json:"name" bson:"name"`
	DamageRelations TypeDamageRelations `json:"damage_relations" bson:"damage_relations"`
	MoveDamageClass ResourceReference   `json:"move_damage_class" bson:"move_damage_class"`
	LastSyncedAt    int64               `json:"-" bson:"last_synced_at,omitempty"`
}

type PokemonListTypeResponse struct {
	Count    int                   `json:"count" bson:"count"`
	Next     *string               `json:"next" bson:"next"`
	Previous *string               `json:"previous" bson:"previous"`
	Results  []PokemonTypeListItem `json:"results" bson:"results"`
}

type PokemonTypeDetailResponse struct {
	TypeID          int                 `json:"id" bson:"id"`
	Name            string              `json:"name" bson:"name"`
	DamageRelations TypeDamageRelations `json:"damage_relations" bson:"damage_relations"`
	MoveDamageClass ResourceReference   `json:"move_damage_class" bson:"move_damage_class"`
}

type PokemonTypeListItem struct {
	TypeID int    `json:"id" bson:"id"`
	Name   string `json:"name" bson:"name"`
	URL    string `json:"url" bson:"url"`
}

type PokemonTypeListItemDocument struct {
	TypeID int    `json:"id" bson:"id"`
	Name   string `json:"name" bson:"name"`
	URL    string `json:"url" bson:"url"`
}

type PokemonWeaknessTypes struct {
	TypeID        int     `json:"type_id" bson:"type_id"`
	Name          string  `json:"name" bson:"name"`
	WeaknessPoint float64 `json:"weakness_point" bson:"weakness_point"`
}

type PokemonWeaknessResponse struct {
	PokeID      int                    `json:"id" bson:"id"`
	PokemonName string                 `json:"pokemon_name" bson:"pokemon_name"`
	Weakness    []PokemonWeaknessTypes `json:"weakness" bson:"weakness"`
}

type PokemonInfo struct {
	ID   int    `json:"id" bson:"id"`
	Name string `json:"name" bson:"name"`
}

type PokemonDamageRelations struct {
	PokeID          int                 `json:"pokemon_id"`
	PokemonName     string              `json:"pokemon_name"`
	DoubleDamgeFrom []ResourceReference `json:"double_damage_from"`
	HalfDamgeFrom   []ResourceReference `json:"half_damage_from" `
}
