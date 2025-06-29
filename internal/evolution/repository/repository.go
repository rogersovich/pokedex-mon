package repository

import (
	"context"
	"fmt"
	"pokedex/database"
	"pokedex/internal/evolution/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const evolutionCollectionName = "evolutions"
const pokemonCollectionName = "pokemons"

type EvolutionRepository interface {
	SaveEvolution(ctx context.Context, pokemon model.EvolutionChain) error
	GetEvolutionByID(ctx context.Context, id int) (model.EvolutionChain, error)
	GetEvolutionByName(ctx context.Context, name string) (model.EvolutionChain, error)
	GetEvolutionPokemonType(ctx context.Context, id int) (model.EvolutionPokemonResponse, error)
}

type MongoEvolutionRepository struct {
	collection        *mongo.Collection
	pokemonCollection *mongo.Collection
}

func NewMongoEvolutionRepository() *MongoEvolutionRepository {
	return &MongoEvolutionRepository{
		collection:        database.MongoDatabase.Collection(evolutionCollectionName),
		pokemonCollection: database.MongoDatabase.Collection(pokemonCollectionName),
	}
}

func (r *MongoEvolutionRepository) SaveEvolution(ctx context.Context, evolution model.EvolutionChain) error {
	doc := model.EvolutionChainDocument{
		EvolutionID:     evolution.ID,
		Chain:           evolution.Chain,
		BabyTriggerItem: evolution.BabyTriggerItem,
		LastSyncedAt:    time.Now().Unix(),
	}

	filter := bson.M{"id": doc.EvolutionID}
	update := bson.M{"$set": doc}
	opts := options.Update().SetUpsert(true)

	_, err := r.collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return fmt.Errorf("failed to save data to MongoDB: %w", err)
	}
	return nil
}

func (r *MongoEvolutionRepository) GetEvolutionByID(ctx context.Context, id int) (model.EvolutionChain, error) {
	var doc model.EvolutionChainDocument
	filter := bson.M{"id": id}
	err := r.collection.FindOne(ctx, filter).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return model.EvolutionChain{}, fmt.Errorf("data not found: %d", id)
		}
		return model.EvolutionChain{}, fmt.Errorf("failed to retrieve data by ID from DB: %w", err)
	}

	return r.toDetail(doc), nil
}

func (r *MongoEvolutionRepository) GetEvolutionByName(ctx context.Context, name string) (model.EvolutionChain, error) {
	var doc model.EvolutionChainDocument
	filter := bson.M{"name": name}
	err := r.collection.FindOne(ctx, filter).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return model.EvolutionChain{}, fmt.Errorf("data not found: %s", name)
		}
		return model.EvolutionChain{}, fmt.Errorf("failed to retrieve data by ID from DB: %w", err)
	}

	return r.toDetail(doc), nil
}

func (r *MongoEvolutionRepository) GetEvolutionPokemonType(ctx context.Context, pokemon_id int) (model.EvolutionPokemonResponse, error) {
	var doc model.EvolutionPokemonTypeDocument
	filter := bson.M{"id": pokemon_id}
	findOptions := options.FindOne()

	projection := bson.D{
		{Key: "id", Value: 1},
		{Key: "name", Value: 1},
		{Key: "types", Value: 1},
		// {Key: "_id", Value: 0}, // Uncomment to explicitly exclude _id
	}

	findOptions.SetProjection(projection)

	err := r.pokemonCollection.FindOne(ctx, filter, findOptions).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return model.EvolutionPokemonResponse{}, fmt.Errorf("data not found: %d", pokemon_id)
		}
		return model.EvolutionPokemonResponse{}, fmt.Errorf("failed to retrieve data by ID from DB: %w", err)
	}

	return r.toDetailPokemonType(doc), nil
}

func (r *MongoEvolutionRepository) toDetail(doc model.EvolutionChainDocument) model.EvolutionChain {
	evoDetail := model.EvolutionChain{
		ID:              doc.EvolutionID,
		BabyTriggerItem: doc.BabyTriggerItem,
		Chain:           doc.Chain,
	}

	return evoDetail
}

func (r *MongoEvolutionRepository) toDetailPokemonType(doc model.EvolutionPokemonTypeDocument) model.EvolutionPokemonResponse {
	return model.EvolutionPokemonResponse{
		ID:    doc.PokemonID,
		Name:  doc.Name,
		Types: doc.Types,
	}
}
