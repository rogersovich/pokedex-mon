package repository

import (
	"context"
	"fmt"
	"pokedex/database"
	"pokedex/internal/pokedexes/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const pokedexesCollectionName = "pokedexes"

type PokedexesRepository interface {
	SavePokedexes(ctx context.Context, ability model.PokedexesDetail) error
	GetPokedexByID(ctx context.Context, id int) (model.PokedexesDetail, error)
	GetPokedexList(ctx context.Context, limit, offset int, baseUrl string) ([]model.ListPokedexesItem, int64, error)
}

type MongoPokedexesRepository struct {
	collection *mongo.Collection
}

func NewMongoPokedexesRepository() *MongoPokedexesRepository {
	return &MongoPokedexesRepository{
		collection: database.MongoDatabase.Collection(pokedexesCollectionName),
	}
}

func (r *MongoPokedexesRepository) SavePokedexes(ctx context.Context, pokedexes model.PokedexesDetail) error {
	doc := model.PokedexesSaveDocument{
		PokedexesID:    pokedexes.ID,
		IsMainSeries:   pokedexes.IsMainSeries,
		Region:         pokedexes.Region,
		VersionGroup:   pokedexes.VersionGroup,
		Name:           pokedexes.Name,
		Names:          pokedexes.Names,
		Descriptions:   pokedexes.Descriptions,
		PokemonEntries: pokedexes.PokemonEntries,
		LastSyncedAt:   time.Now().Unix(), // Menyimpan timestamp saat ini
	}

	filter := bson.M{"id": doc.PokedexesID} // Filter berdasarkan ID PokeAPI
	update := bson.M{"$set": doc}           // Menggunakan $set untuk memperbarui atau menyisipkan seluruh dokumen

	opts := options.Update().SetUpsert(true) // Opsi upsert: jika tidak ada, sisipkan; jika ada, perbarui.

	_, err := r.collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return fmt.Errorf("failed to save data %s (ID: %d) to MongoDB: %w", pokedexes.Name, pokedexes.ID, err)
	}
	return nil
}

func (r *MongoPokedexesRepository) GetPokedexByID(ctx context.Context, id int) (model.PokedexesDetail, error) {
	var doc model.PokedexesDocument
	filter := bson.M{"id": id} // Mencari berdasarkan PokeAPI ID
	err := r.collection.FindOne(ctx, filter).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return model.PokedexesDetail{}, fmt.Errorf("pokedex with ID %d not found", id)
		}
		return model.PokedexesDetail{}, fmt.Errorf("failed to retrieve pokedex by ID from DB: %w", err)
	}
	return r.toDetail(doc), nil
}

func (r *MongoPokedexesRepository) GetPokedexList(ctx context.Context, limit, offset int, baseUrl string) ([]model.ListPokedexesItem, int64, error) {
	totalCount, err := r.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count pokedex in DB: %w", err)
	}

	findOptions := options.Find()
	findOptions.SetLimit(int64(limit))
	findOptions.SetSkip(int64(offset))
	findOptions.SetSort(bson.D{{Key: "id", Value: 1}}) // Sort by actual Pokemon ID

	cursor, err := r.collection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to retrieve pokedex list from DB: %w", err)
	}
	defer cursor.Close(ctx)

	var docs []model.PokedexesListDocument
	if err = cursor.All(ctx, &docs); err != nil {
		return nil, 0, fmt.Errorf("failed to decode pokedex list from DB: %w", err)
	}

	var listType []model.ListPokedexesItem
	for _, doc := range docs {
		listType = append(listType, r.toPokedexList(doc, baseUrl))
	}

	return listType, totalCount, nil
}

func (r *MongoPokedexesRepository) toDetail(doc model.PokedexesDocument) model.PokedexesDetail {
	return model.PokedexesDetail{
		ID:             doc.PokedexesID,
		IsMainSeries:   doc.IsMainSeries,
		Region:         doc.Region,
		VersionGroup:   doc.VersionGroup,
		Name:           doc.Name,
		Names:          doc.Names,
		Descriptions:   doc.Descriptions,
		PokemonEntries: doc.PokemonEntries,
	}
}

func (r *MongoPokedexesRepository) toPokedexList(doc model.PokedexesListDocument, baseUrl string) model.ListPokedexesItem {
	return model.ListPokedexesItem{
		ID:   doc.ID,
		Name: doc.Name,
		URL:  baseUrl + fmt.Sprintf("/%d", doc.ID),
	}
}
