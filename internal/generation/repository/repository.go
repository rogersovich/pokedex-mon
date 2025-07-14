package repository

import (
	"context"
	"fmt"
	"pokedex/database"
	"pokedex/internal/generation/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const generationCollectionName = "generations"

type GenerationRepository interface {
	SaveGenerations(ctx context.Context, ability model.GenerationDetail) error
	// GetGenerationByID(ctx context.Context, id int) (model.GenerationDetail, error)
	// GetGenerationList(ctx context.Context, limit, offset int, baseUrl string) ([]model.ListGenerationItem, int64, error)
}

type MongoGenerationRepository struct {
	collection *mongo.Collection
}

func NewMongoGenerationRepository() *MongoGenerationRepository {
	return &MongoGenerationRepository{
		collection: database.MongoDatabase.Collection(generationCollectionName),
	}
}

func (r *MongoGenerationRepository) SaveGenerations(ctx context.Context, generation model.GenerationDetail) error {
	doc := model.GenerationSaveDocument{
		GenerationID:   generation.ID,
		Name:           generation.Name,
		Names:          generation.Names,
		Abilities:      generation.Abilities,
		MainRegion:     generation.MainRegion,
		Moves:          generation.Moves,
		PokemonSpecies: generation.PokemonSpecies,
		Types:          generation.Types,
		VersionGroups:  generation.VersionGroups,
		LastSyncedAt:   time.Now().Unix(), // Menyimpan timestamp saat ini
	}

	filter := bson.M{"id": doc.ID} // Filter berdasarkan ID PokeAPI
	update := bson.M{"$set": doc}  // Menggunakan $set untuk memperbarui atau menyisipkan seluruh dokumen

	opts := options.Update().SetUpsert(true) // Opsi upsert: jika tidak ada, sisipkan; jika ada, perbarui.

	_, err := r.collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return fmt.Errorf("failed to save data %s (ID: %d) to MongoDB: %w", generation.Name, generation.ID, err)
	}
	return nil
}
