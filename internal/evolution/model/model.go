package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type ResourceReference struct {
	Name string `json:"name" bson:"name"`
	URL  string `json:"url" bson:"url"`
}

type EvolutionChain struct {
	ID              int               `json:"id" bson:"id"`
	BabyTriggerItem ResourceReference `json:"baby_trigger_item" bson:"baby_trigger_item"`
	Chain           ChainLink         `json:"chain" bson:"chain"`
}

type ChainLink struct {
	IsBaby           bool                     `json:"is_baby" bson:"is_baby"`
	Species          ResourceReference        `json:"species" bson:"species"`
	EvolutionDetails []EvolutionDetail        `json:"evolution_details" bson:"evolution_details"`
	EvolvesTo        []ChainLink              `json:"evolves_to" bson:"evolves_to"`
	EvolutionType    EvolutionPokemonResponse `json:"evolution_type" bson:"evolution_type"`
}

type EvolutionDetail struct {
	Item                  ResourceReference `json:"item" bson:"item"`
	Trigger               ResourceReference `json:"trigger" bson:"trigger"`
	Gender                int               `json:"gender" bson:"gender"`
	HeldItem              ResourceReference `json:"held_item" bson:"held_item"`
	KnownMove             ResourceReference `json:"known_move" bson:"known_move"`
	KnownMoveType         ResourceReference `json:"known_move_type" bson:"known_move_type"`
	Location              ResourceReference `json:"location" bson:"location"`
	MinLevel              int               `json:"min_level" bson:"min_level"`
	MinHappiness          int               `json:"min_happiness" bson:"min_happiness"`
	MinBeauty             int               `json:"min_beauty" bson:"min_beauty"`
	MinAffection          int               `json:"min_affection" bson:"min_affection"`
	PartySpecies          ResourceReference `json:"party_species" bson:"party_species"`
	PartyType             ResourceReference `json:"party_type" bson:"party_type"`
	RelativePhysicalStats int               `json:"relative_physical_stats" bson:"relative_physical_stats"`
	TimeOfDay             string            `json:"time_of_day" bson:"time_of_day"`
	TradeSpecies          ResourceReference `json:"trade_species" bson:"trade_species"`
	TurnUpsideDown        bool              `json:"turn_upside_down" bson:"turn_upside_down"`
}

type EvolutionListItem struct {
	URL string `json:"url" bson:"url"`
}

type EvolutionListResponse struct {
	Count    int                 `json:"count" bson:"count"`
	Next     *string             `json:"next" bson:"next"`
	Previous *string             `json:"previous" bson:"previous"`
	Results  []EvolutionListItem `json:"results" bson:"results"`
}

type EvolutionPokemonType struct {
	Slot int               `json:"slot" bson:"slot"`
	Type ResourceReference `json:"type" bson:"type"`
}

type EvolutionPokemonResponse struct {
	ID    int                    `json:"id" bson:"id"`
	Name  string                 `json:"name" bson:"name"`
	Types []EvolutionPokemonType `json:"types" bson:"types"`
}

type EvolutionChainDocument struct {
	ID              primitive.ObjectID `bson:"_id,omitempty"`
	EvolutionID     int                `json:"id" bson:"id"`
	BabyTriggerItem ResourceReference  `json:"baby_trigger_item" bson:"baby_trigger_item"`
	Chain           ChainLink          `json:"chain" bson:"chain"`
	LastSyncedAt    int64              `json:"-" bson:"last_synced_at,omitempty"`
}

type EvolutionPokemonTypeDocument struct {
	ID           primitive.ObjectID     `bson:"_id,omitempty"`
	PokemonID    int                    `json:"id" bson:"id"`
	Name         string                 `json:"name" bson:"name"`
	Types        []EvolutionPokemonType `json:"types" bson:"types"`
	LastSyncedAt int64                  `json:"-" bson:"last_synced_at,omitempty"`
}
