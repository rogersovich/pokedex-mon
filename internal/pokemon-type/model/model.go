package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type ResourceReference struct {
	Name string `json:"name" bson:"name"`
	URL  string `json:"url" bson:"url"`
}

type TypeDamageRelation struct {
	DoubleDamgeFrom []ResourceReference `json:"double_damage_from" bson:"double_damage_from"`
	DoubleDamgeTo   []ResourceReference `json:"double_damage_to" bson:"double_damage_to"`
	HalfDamgeFrom   []ResourceReference `json:"half_damage_from" bson:"half_damage_from"`
	HalfDamgeTo     []ResourceReference `json:"half_damage_to" bson:"half_damage_to"`
	NoDamgeFrom     []ResourceReference `json:"no_damage_from" bson:"no_damage_from"`
	NoDamgeTo       []ResourceReference `json:"no_damage_to" bson:"no_damage_to"`
}

type PokemonTypeDocument struct {
	ID              primitive.ObjectID `bson:"_id,omitempty"`
	TypeID          int                `json:"id" bson:"id"`
	Name            string             `json:"name" bson:"name"`
	DamageRelation  TypeDamageRelation `json:"damage_relation" bson:"damage_relation"`
	MoveDamageClass ResourceReference  `json:"move_damage_class" bson:"move_damage_class"`
	LastSyncedAt    int64              `json:"-" bson:"last_synced_at,omitempty"`
}

type PokemonListTypeResponse struct {
	Count    int                   `json:"count" bson:"count"`
	Next     *string               `json:"next" bson:"next"`
	Previous *string               `json:"previous" bson:"previous"`
	Results  []PokemonTypeListItem `json:"results" bson:"results"`
}

type PokemonTypeDetailResponse struct {
	TypeID          int                `json:"id" bson:"id"`
	Name            string             `json:"name" bson:"name"`
	DamageRelation  TypeDamageRelation `json:"damage_relation" bson:"damage_relation"`
	MoveDamageClass ResourceReference  `json:"move_damage_class" bson:"move_damage_class"`
}

type PokemonTypeListItem struct {
	Name string `json:"name" bson:"name"`
	URL  string `json:"url" bson:"url"`
}
