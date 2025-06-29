package repository

import (
	"context"
	"fmt"
	"time"

	"pokedex/database"
	"pokedex/internal/pokemon-species/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const pokemonSpeciesCollectionName = "pokemon-species"

type PokemonSpeciesRepository interface {
	SavePokemonSpecies(ctx context.Context, pokemon model.PokemonSpeciesDetail) error
	GetPokemonSpeciesByID(ctx context.Context, id int) (model.PokemonSpeciesDetail, error)
	GetPokemonSpeciesByName(ctx context.Context, name string) (model.PokemonSpeciesDetail, error)
}

type MongoPokemonSpeciesRepository struct {
	collection *mongo.Collection
}

func NewMongoPokemonSpeciesRepository() *MongoPokemonSpeciesRepository {
	return &MongoPokemonSpeciesRepository{
		collection: database.MongoDatabase.Collection(pokemonSpeciesCollectionName),
	}
}

func (r *MongoPokemonSpeciesRepository) SavePokemonSpecies(ctx context.Context, pokemon model.PokemonSpeciesDetail) error {
	doc := model.PokemonSpeciesDocument{
		PokeAPIID:            pokemon.PokeAPIID,
		Name:                 pokemon.Name,
		BaseHappiness:        pokemon.BaseHappiness,
		CaptureRate:          pokemon.CaptureRate,
		Color:                pokemon.Color,
		EggGroups:            pokemon.EggGroups,
		EvolutionChain:       pokemon.EvolutionChain,
		EvolvesFromSpecies:   pokemon.EvolvesFromSpecies,
		FlavorTextEntries:    pokemon.FlavorTextEntries,
		FormDescriptions:     pokemon.FormDescriptions,
		FormsSwitchable:      pokemon.FormsSwitchable,
		GenderRate:           pokemon.GenderRate,
		Genera:               pokemon.Genera,
		Generation:           pokemon.Generation,
		GrowthRate:           pokemon.GrowthRate,
		Habitat:              pokemon.Habitat,
		HasGenderDifferences: pokemon.HasGenderDifferences,
		HatchCounter:         pokemon.HatchCounter,
		IsBaby:               pokemon.IsBaby,
		IsLegendary:          pokemon.IsLegendary,
		IsMythical:           pokemon.IsMythical,
		Names:                pokemon.Names,
		Order:                pokemon.Order,
		PalParkEncounters:    pokemon.PalParkEncounters,
		PokedexNumbers:       pokemon.PokedexNumbers,
		Shape:                pokemon.Shape,
		Varieties:            pokemon.Varieties,
		LastSyncedAt:         time.Now().Unix(),
	}

	filter := bson.M{"id": doc.PokeAPIID}
	update := bson.M{"$set": doc}
	opts := options.Update().SetUpsert(true)

	_, err := r.collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return fmt.Errorf("failed to save data %s (ID: %d) to MongoDB: %w", pokemon.Name, pokemon.PokeAPIID, err)
	}
	return nil
}

func (r *MongoPokemonSpeciesRepository) GetPokemonSpeciesByID(ctx context.Context, id int) (model.PokemonSpeciesDetail, error) {
	var doc model.PokemonSpeciesDocument
	filter := bson.M{"id": id}
	err := r.collection.FindOne(ctx, filter).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return model.PokemonSpeciesDetail{}, fmt.Errorf("pokemon not found: %d", id)
		}
		return model.PokemonSpeciesDetail{}, fmt.Errorf("failed to retrieve pokemon by ID from DB: %w", err)
	}

	return r.toDetail(doc), nil
}

func (r *MongoPokemonSpeciesRepository) GetPokemonSpeciesByName(ctx context.Context, name string) (model.PokemonSpeciesDetail, error) {
	var doc model.PokemonSpeciesDocument
	filter := bson.M{"name": name}
	err := r.collection.FindOne(ctx, filter).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return model.PokemonSpeciesDetail{}, fmt.Errorf("pokemon not found: %s", name)
		}
		return model.PokemonSpeciesDetail{}, fmt.Errorf("failed to retrieve pokemon by name from DB: %w", err)
	}

	return r.toDetail(doc), nil
}

func (r *MongoPokemonSpeciesRepository) toDetail(doc model.PokemonSpeciesDocument) model.PokemonSpeciesDetail {
	return model.PokemonSpeciesDetail{
		PokeAPIID:            doc.PokeAPIID,
		Name:                 doc.Name,
		BaseHappiness:        doc.BaseHappiness,
		CaptureRate:          doc.CaptureRate,
		Color:                doc.Color,
		EggGroups:            doc.EggGroups,
		EvolutionChain:       doc.EvolutionChain,
		EvolvesFromSpecies:   doc.EvolvesFromSpecies,
		FlavorTextEntries:    doc.FlavorTextEntries,
		FormDescriptions:     doc.FormDescriptions,
		FormsSwitchable:      doc.FormsSwitchable,
		GenderRate:           doc.GenderRate,
		Genera:               doc.Genera,
		Generation:           doc.Generation,
		GrowthRate:           doc.GrowthRate,
		Habitat:              doc.Habitat,
		HasGenderDifferences: doc.HasGenderDifferences,
		HatchCounter:         doc.HatchCounter,
		IsBaby:               doc.IsBaby,
		IsLegendary:          doc.IsLegendary,
		IsMythical:           doc.IsMythical,
		Names:                doc.Names,
		Order:                doc.Order,
		PalParkEncounters:    doc.PalParkEncounters,
		PokedexNumbers:       doc.PokedexNumbers,
		Shape:                doc.Shape,
		Varieties:            doc.Varieties,
	}
}
