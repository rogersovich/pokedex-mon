package repository

import (
	"context"
	"fmt"
	"pokedex/database"
	"pokedex/internal/pokemon-type/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const typeCollectionName = "pokemon-types"
const pokemonCollectionName = "pokemons"

type PokemonTypeRepository interface {
	SavePokemonType(ctx context.Context, pokemon model.PokemonTypeDetailResponse) error
}

type MongoPokemonTypeRepository struct {
	collection        *mongo.Collection
	pokemonCollection *mongo.Collection
}

func NewMongoPokemonTypeRepository() *MongoPokemonTypeRepository {
	return &MongoPokemonTypeRepository{
		collection:        database.MongoDatabase.Collection(typeCollectionName),
		pokemonCollection: database.MongoDatabase.Collection(pokemonCollectionName),
	}
}

func (r *MongoPokemonTypeRepository) SavePokemonType(ctx context.Context, pokemon_type model.PokemonTypeDetailResponse) error {
	doc := model.PokemonTypeDocument{
		TypeID:          pokemon_type.TypeID,
		Name:            pokemon_type.Name,
		DamageRelation:  pokemon_type.DamageRelation,
		MoveDamageClass: pokemon_type.MoveDamageClass,
		LastSyncedAt:    time.Now().Unix(),
	}

	filter := bson.M{"id": doc.TypeID}
	update := bson.M{"$set": doc}
	opts := options.Update().SetUpsert(true)

	_, err := r.collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return fmt.Errorf("failed to save data to MongoDB: %w", err)
	}
	return nil
}
