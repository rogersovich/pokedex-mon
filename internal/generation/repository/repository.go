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
	GetGenerationByID(ctx context.Context, id int) (model.GenerationDetail, error)
	GetGenerationList(ctx context.Context, limit, offset int, baseUrl string) ([]model.ListGenerationItem, int64, error)
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

func (r *MongoGenerationRepository) GetGenerationByID(ctx context.Context, id int) (model.GenerationDetail, error) {
	var doc model.GenerationDocument
	filter := bson.M{"id": id} // Mencari berdasarkan PokeAPI ID
	err := r.collection.FindOne(ctx, filter).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return model.GenerationDetail{}, fmt.Errorf("generation with ID %d not found", id)
		}
		return model.GenerationDetail{}, fmt.Errorf("failed to retrieve Generation by ID from DB: %w", err)
	}
	return r.toDetail(doc), nil
}

func (r *MongoGenerationRepository) GetGenerationList(ctx context.Context, limit, offset int, baseUrl string) ([]model.ListGenerationItem, int64, error) {
	totalCount, err := r.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count Generation in DB: %w", err)
	}

	findOptions := options.Find()
	findOptions.SetLimit(int64(limit))
	findOptions.SetSkip(int64(offset))
	findOptions.SetSort(bson.D{{Key: "id", Value: 1}}) // Sort by actual Pokemon ID

	cursor, err := r.collection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to retrieve Generation list from DB: %w", err)
	}
	defer cursor.Close(ctx)

	var docs []model.GenerationListDocument
	if err = cursor.All(ctx, &docs); err != nil {
		return nil, 0, fmt.Errorf("failed to decode Generation list from DB: %w", err)
	}

	var listType []model.ListGenerationItem
	for _, doc := range docs {
		listType = append(listType, r.toListItem(doc, baseUrl))
	}

	return listType, totalCount, nil
}

func (r *MongoGenerationRepository) toDetail(doc model.GenerationDocument) model.GenerationDetail {
	return model.GenerationDetail{
		ID:             doc.GenerationID,
		Name:           doc.Name,
		Names:          doc.Names,
		Abilities:      doc.Abilities,
		MainRegion:     doc.MainRegion,
		Moves:          doc.Moves,
		PokemonSpecies: doc.PokemonSpecies,
		Types:          doc.Types,
		VersionGroups:  doc.VersionGroups,
	}
}

func (r *MongoGenerationRepository) toListItem(doc model.GenerationListDocument, baseUrl string) model.ListGenerationItem {
	return model.ListGenerationItem{
		ID:   doc.ID,
		Name: doc.Name,
		URL:  baseUrl + fmt.Sprintf("/%d", doc.ID),
	}
}
