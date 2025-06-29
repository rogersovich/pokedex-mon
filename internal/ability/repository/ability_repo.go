package repository

import (
	"context"
	"fmt"
	"pokedex/database"
	"pokedex/internal/ability/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const abilityCollectionName = "abilities"

// AbilityRepository defines the interface for persisting and retrieving Ability data.
type AbilityRepository interface {
	SaveAbility(ctx context.Context, ability model.AbilityDetail) error
	GetAbilityByID(ctx context.Context, id int) (model.AbilityDetail, error)
	GetAbilityByName(ctx context.Context, name string) (model.AbilityDetail, error)
}

// MongoAbilityRepository implements the AbilityRepository interface for MongoDB.
type MongoAbilityRepository struct {
	collection *mongo.Collection
}

// NewMongoAbilityRepository creates a new MongoDB repository for abilities.
func NewMongoAbilityRepository() *MongoAbilityRepository {
	return &MongoAbilityRepository{
		collection: database.MongoDatabase.Collection(abilityCollectionName),
	}
}

// SaveAbility saves an ability detail to MongoDB.
// It uses Upsert to either insert a new document or update an existing one based on 'id'.
func (r *MongoAbilityRepository) SaveAbility(ctx context.Context, ability model.AbilityDetail) error {
	doc := model.AbilityDocument{
		AbilityID:         ability.ID,
		EffectChanges:     ability.EffectChanges,
		EffectEntries:     ability.EffectEntries,
		FlavorTextEntries: ability.FlavorTextEntries,
		Generation:        ability.Generation,
		IsMainSeries:      ability.IsMainSeries,
		Name:              ability.Name,
		Names:             ability.Names,
		Pokemon:           ability.Pokemon,
		LastSyncedAt:      time.Now().Unix(), // Menyimpan timestamp saat ini
	}

	filter := bson.M{"id": doc.AbilityID} // Filter berdasarkan ID PokeAPI
	update := bson.M{"$set": doc}         // Menggunakan $set untuk memperbarui atau menyisipkan seluruh dokumen

	opts := options.Update().SetUpsert(true) // Opsi upsert: jika tidak ada, sisipkan; jika ada, perbarui.

	_, err := r.collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return fmt.Errorf("failed to save ability %s (ID: %d) to MongoDB: %w", ability.Name, ability.ID, err)
	}
	return nil
}

// GetAbilityByID retrieves an ability by its original PokeAPI ID from MongoDB.
func (r *MongoAbilityRepository) GetAbilityByID(ctx context.Context, id int) (model.AbilityDetail, error) {
	var doc model.AbilityDocument
	filter := bson.M{"id": id} // Mencari berdasarkan PokeAPI ID
	err := r.collection.FindOne(ctx, filter).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return model.AbilityDetail{}, fmt.Errorf("ability with ID %d not found", id)
		}
		return model.AbilityDetail{}, fmt.Errorf("failed to retrieve ability by ID from DB: %w", err)
	}
	return r.toDetail(doc), nil
}

// GetAbilityByName retrieves an ability by its name from MongoDB.
func (r *MongoAbilityRepository) GetAbilityByName(ctx context.Context, name string) (model.AbilityDetail, error) {
	var doc model.AbilityDocument
	filter := bson.M{"name": name} // Mencari berdasarkan nama Ability
	err := r.collection.FindOne(ctx, filter).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return model.AbilityDetail{}, fmt.Errorf("ability with name '%s' not found", name)
		}
		return model.AbilityDetail{}, fmt.Errorf("failed to retrieve ability by name from DB: %w", err)
	}
	return r.toDetail(doc), nil
}

func (r *MongoAbilityRepository) GetAbilityList(ctx context.Context, limit, offset int) ([]model.AbilityDetail, int64, error) {
	totalCount, err := r.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count abilitys in DB: %w", err)
	}

	findOptions := options.Find()
	findOptions.SetLimit(int64(limit))
	findOptions.SetSkip(int64(offset))
	findOptions.SetSort(bson.D{{Key: "id", Value: 1}}) // Sort by actual ability ID

	cursor, err := r.collection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to retrieve ability list from DB: %w", err)
	}
	defer cursor.Close(ctx)

	var abilityDocs []model.AbilityDocument
	if err = cursor.All(ctx, &abilityDocs); err != nil {
		return nil, 0, fmt.Errorf("failed to decode ability list from DB: %w", err)
	}

	var abilityDetails []model.AbilityDetail
	for _, doc := range abilityDocs {
		abilityDetails = append(abilityDetails, r.toDetail(doc))
	}

	return abilityDetails, totalCount, nil
}

// toDetail helper function converts an AbilityDocument to an AbilityDetail model.
// This is useful if your internal document structure differs slightly from the API model.
func (r *MongoAbilityRepository) toDetail(doc model.AbilityDocument) model.AbilityDetail {
	return model.AbilityDetail{
		ID:                doc.AbilityID,
		EffectChanges:     doc.EffectChanges,
		EffectEntries:     doc.EffectEntries,
		FlavorTextEntries: doc.FlavorTextEntries,
		Generation:        doc.Generation,
		IsMainSeries:      doc.IsMainSeries,
		Name:              doc.Name,
		Names:             doc.Names,
		Pokemon:           doc.Pokemon}
}
