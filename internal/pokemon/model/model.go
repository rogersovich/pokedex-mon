package model

import (
	evolution_model "pokedex/internal/evolution/model"
	"pokedex/utils"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// PokemonListItem represents an item in the initial list from PokeAPI
type PokemonListItem struct {
	Name    string        `json:"name" bson:"name"`
	URL     string        `json:"url" bson:"url"`
	Types   []PokemonType `json:"types" bson:"types"`
	Sprites Sprites       `json:"sprites" bson:"sprites"`
}

// PokemonListResponse represents the full response for a list of pokemons
type PokemonListResponse struct {
	Count    int               `json:"count" bson:"count"`
	Next     *string           `json:"next" bson:"next"`
	Previous *string           `json:"previous" bson:"previous"`
	Results  []PokemonListItem `json:"results" bson:"results"`
}

type Sprites struct {
	FrontDefault string                 `json:"front_default" bson:"front_default"`
	Other        map[string]interface{} `json:"other,omitempty" bson:"other,omitempty"`
}

// ResourceReference represents a generic name and URL reference
type ResourceReference struct {
	Name string `json:"name" bson:"name"`
	URL  string `json:"url" bson:"url"`
}

type PokemonType struct {
	Slot int               `json:"slot" bson:"slot"`
	Type ResourceReference `json:"type" bson:"type"`
}

type PokemonStat struct {
	BaseStat int               `json:"base_stat" bson:"base_stat"`
	Effort   int               `json:"effort" bson:"effort"`
	Stat     ResourceReference `json:"stat" bson:"stat"`
}

type PokemonAbility struct {
	Ability  ResourceReference `json:"ability" bson:"ability"`
	IsHidden bool              `json:"is_hidden" bson:"is_hidden"`
	Slot     int               `json:"slot" bson:"slot"`
}

type PokemonGameIndices struct {
	Version   ResourceReference `json:"version" bson:"version"`
	GameIndex int               `json:"game_index" bson:"game_index"`
}

type PokemonMoves struct {
	Move                ResourceReference `json:"move" bson:"move"`
	VersionGroupDetails []struct {
		LevelLearnedAt  int               `json:"level_learned_at" bson:"level_learned_at"`
		MoveLearnMethod ResourceReference `json:"move_learn_method" bson:"move_learn_method"`
		Order           int               `json:"order" bson:"order"`
		VersionGroup    ResourceReference `json:"version_group" bson:"version_group"`
	} `json:"version_group_details" bson:"version_group_details"`
}

type PokemonTraining struct {
	CaptureRate        int     `json:"capture_rate"`
	CaptureRatePercent float64 `json:"capture_rate_percent"`
	BaseExperience     int     `json:"base_experience"`
	BaseHappiness      int     `json:"base_happiness"`
	GrowthRate         string  `json:"growth_rate"`
}

type PokemonBreeding struct {
	EggGroups    []ResourceReference            `json:"egg_groups"`
	GenderRate   utils.GenderDistributionResult `json:"gender_rate"`
	HatchCounter int                            `json:"hatch_counter"`
	EggCycles    int                            `json:"egg_cycles"`
}

type PokemonStatFull struct {
	BaseStat int    `json:"base_stat" bson:"base_stat"`
	Effort   int    `json:"effort" bson:"effort"`
	StatName string `json:"stat_name" bson:"stat_name"`
	MaxStat  int    `json:"max_stat" bson:"max_stat"`
	MinStat  int    `json:"min_stat" bson:"min_stat"`
}

type GroupedMoveInfo struct {
	MoveName        string `json:"move_name"`
	MoveURL         string `json:"move_url"`
	LevelLearnedAt  int    `json:"level_learned_at"`
	MoveLearnMethod string `json:"move_learn_method"`
	Order           int    `json:"order"`
}

type MovesByLearnMethod struct {
	MethodName string            `json:"method_name"`
	Moves      []GroupedMoveInfo `json:"moves"`
}

type GroupedVersionMoves struct {
	GroupName     string               `json:"group_name"`
	MovesByMethod []MovesByLearnMethod `json:"moves_by_method"`
}

type PokemonDetail struct {
	ID                     int                  `json:"id" bson:"id"`
	Name                   string               `json:"name" bson:"name"`
	Height                 int                  `json:"height" bson:"height"`
	Weight                 int                  `json:"weight" bson:"weight"`
	BaseExperience         int                  `json:"base_experience" bson:"base_experience"`
	Sprites                Sprites              `json:"sprites" bson:"sprites"`
	Types                  []PokemonType        `json:"types" bson:"types"`
	Stats                  []PokemonStat        `json:"stats" bson:"stats"`
	Abilities              []PokemonAbility     `json:"abilities" bson:"abilities"`
	Forms                  []ResourceReference  `json:"forms" bson:"forms"`
	GameIndices            []PokemonGameIndices `json:"game_indices" bson:"game_indices"`
	HeldItems              []interface{}        `json:"held_items" bson:"held_items"`
	IsDefault              bool                 `json:"is_default" bson:"is_default"`
	LocationAreaEncounters string               `json:"location_area_encounters" bson:"location_area_encounters"`
	Moves                  []PokemonMoves       `json:"moves" bson:"moves"`
	Order                  int                  `json:"order" bson:"order"`
}

type PokemonOtherNames struct {
	Language string `json:"language"`
	Name     string `json:"name"`
}

type PokemonDetailResponse struct {
	ID             int                            `json:"id"`
	Name           string                         `json:"name"`
	Height         int                            `json:"height"`
	Weight         int                            `json:"weight"`
	BaseExperience int                            `json:"base_experience"`
	Sprites        Sprites                        `json:"sprites"`
	Types          []PokemonType                  `json:"types"`
	Stats          []PokemonStatFull              `json:"stats"`
	Abilities      []PokemonAbility               `json:"abilities"`
	Evolution      evolution_model.EvolutionChain `json:"evolution"`
	GroupedMoves   []GroupedVersionMoves          `json:"grouped_moves"`
	Order          int                            `json:"order"`
	Habitat        string                         `json:"habitat"`
	Thumbnail      string                         `json:"thumbnail"`
	Training       PokemonTraining                `json:"training"`
	Breeding       PokemonBreeding                `json:"breeding"`
	OtherNames     []PokemonOtherNames            `json:"other_names"`
}

// PokemonDocument is the structure to store in MongoDB
type PokemonDocument struct {
	ID                     primitive.ObjectID   `bson:"_id,omitempty"`
	PokemonID              int                  `bson:"id"` // Actual Pokemon ID from PokeAPI
	Name                   string               `bson:"name"`
	Height                 int                  `bson:"height"`
	Weight                 int                  `bson:"weight"`
	BaseExperience         int                  `bson:"base_experience"`
	Sprites                Sprites              `bson:"sprites"`
	Types                  []PokemonType        `bson:"types"`
	Stats                  []PokemonStat        `bson:"stats"`
	Abilities              []PokemonAbility     `bson:"abilities"`
	Forms                  []ResourceReference  `bson:"forms"`
	GameIndices            []PokemonGameIndices `bson:"game_indices"`
	HeldItems              []interface{}        `bson:"held_items"`
	IsDefault              bool                 `bson:"is_default"`
	LocationAreaEncounters string               `bson:"location_area_encounters"`
	Moves                  []PokemonMoves       `bson:"moves"`
	Order                  int                  `bson:"order"`
	LastSyncedAt           int64                `bson:"last_synced_at"`
}
