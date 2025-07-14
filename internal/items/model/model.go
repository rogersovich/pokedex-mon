package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type BaseResourceReference struct {
	Name string `json:"name" bson:"name"`
	URL  string `json:"url" bson:"url"`
}

type BaseItemAttribute struct {
	Name string `json:"name" bson:"name"`
	URL  string `json:"url" bson:"url"`
}

type BaseItemCategory struct {
	Name string `json:"name" bson:"name"`
	URL  string `json:"url" bson:"url"`
}

type BaseItemEffectEntries struct {
	Effect      string                `json:"effect" bson:"effect"`
	Language    BaseResourceReference `json:"language" bson:"language"`
	ShortEffect string                `json:"short_effect" bson:"short_effect"`
}

type BaseItemFlavorText struct {
	Language     BaseResourceReference `json:"language" bson:"language"`
	Text         string                `json:"text" bson:"text"`
	VersionGroup BaseResourceReference `json:"version_group" bson:"version_group"`
}

type BaseItemFlingEffect struct {
	Name string `json:"name" bson:"name"`
	URL  string `json:"url" bson:"url"`
}

type BaseItemGameIndiex struct {
	GameIndex int                   `json:"game_index" bson:"game_index"`
	Version   BaseResourceReference `json:"version" bson:"version"`
}

type BaseItemNames struct {
	Name     string                `json:"name"`
	Language BaseResourceReference `json:"language"`
}

type BaseItemHeldByPokemon struct {
	Pokemon        BaseResourceReference `json:"pokemon" bson:"pokemon"`
	VersionDetails []struct {
		Version BaseResourceReference `json:"version" bson:"version"`
		Rarity  int                   `json:"rarity" bson:"rarity"`
	}
}

type BaseItemMachine struct {
	Machine      BaseResourceReference `json:"machine" bson:"machine"`
	VersionGroup BaseResourceReference `json:"version_group" bson:"version_group"`
}

type BaseItemSprites struct {
	Default string `json:"default" bson:"default"`
}

type BaseItemBabyTriggerFor struct {
	URL string `json:"url" bson:"url"`
}

type ItemDetail struct {
	ID                int                     `json:"id"`
	Name              string                  `json:"name"`
	Names             []BaseItemNames         `json:"names"`
	Cost              int                     `json:"cost"`
	Attributes        []BaseItemAttribute     `json:"attributes"`
	Category          BaseItemCategory        `json:"category"`
	BabyTriggerFor    BaseItemBabyTriggerFor  `json:"baby_trigger_for"`
	EffectEntries     []BaseItemEffectEntries `json:"effect_entries"`
	FlavorTextEntries []BaseItemFlavorText    `json:"flavor_text_entries"`
	FlingEffect       BaseItemFlingEffect     `json:"fling_effect"`
	FlingPower        int                     `json:"fling_power"`
	GameIndicies      []BaseItemGameIndiex    `json:"game_indicies"`
	HeldByPokemon     []BaseItemHeldByPokemon `json:"held_by_pokemon"`
	Machine           []BaseItemMachine       `json:"machine"`
	Sprites           BaseItemSprites         `json:"sprites"`
}

type ListItem struct {
	Count    int            `json:"count" bson:"count"`
	Next     *string        `json:"next" bson:"next"`
	Previous *string        `json:"previous" bson:"previous"`
	Results  []ListItemItem `json:"results" bson:"results"`
}

type ListItemItem struct {
	ID   int    `json:"id" bson:"id"`
	Name string `json:"name" bson:"name"`
	URL  string `json:"url" bson:"url"`
}

type ItemSaveDocument struct {
	ID                primitive.ObjectID      `bson:"_id,omitempty"`
	ItemID            int                     `json:"id"`
	Name              string                  `json:"name"`
	Names             []BaseItemNames         `json:"names"`
	Cost              int                     `json:"cost"`
	Attributes        []BaseItemAttribute     `json:"attributes"`
	Category          BaseItemCategory        `json:"category"`
	BabyTriggerFor    BaseItemBabyTriggerFor  `json:"baby_trigger_for"`
	EffectEntries     []BaseItemEffectEntries `json:"effect_entries"`
	FlavorTextEntries []BaseItemFlavorText    `json:"flavor_text_entries"`
	FlingEffect       BaseItemFlingEffect     `json:"fling_effect"`
	FlingPower        int                     `json:"fling_power"`
	GameIndicies      []BaseItemGameIndiex    `json:"game_indicies"`
	HeldByPokemon     []BaseItemHeldByPokemon `json:"held_by_pokemon"`
	Machine           []BaseItemMachine       `json:"machine"`
	Sprites           BaseItemSprites         `json:"sprites"`
	LastSyncedAt      int64                   `json:"-" bson:"last_synced_at,omitempty"`
}

type ItemDetailDocument struct {
	ID                primitive.ObjectID      `bson:"_id,omitempty"`
	ItemID            int                     `json:"id"`
	Name              string                  `json:"name"`
	Names             []BaseItemNames         `json:"names"`
	Cost              int                     `json:"cost"`
	Attributes        []BaseItemAttribute     `json:"attributes"`
	Category          BaseItemCategory        `json:"category"`
	BabyTriggerFor    BaseItemBabyTriggerFor  `json:"baby_trigger_for"`
	EffectEntries     []BaseItemEffectEntries `json:"effect_entries"`
	FlavorTextEntries []BaseItemFlavorText    `json:"flavor_text_entries"`
	FlingEffect       BaseItemFlingEffect     `json:"fling_effect"`
	FlingPower        int                     `json:"fling_power"`
	GameIndicies      []BaseItemGameIndiex    `json:"game_indicies"`
	HeldByPokemon     []BaseItemHeldByPokemon `json:"held_by_pokemon"`
	Machine           []BaseItemMachine       `json:"machine"`
	Sprites           BaseItemSprites         `json:"sprites"`
	LastSyncedAt      int64                   `json:"-" bson:"last_synced_at,omitempty"`
}

type ItemListDocument struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
