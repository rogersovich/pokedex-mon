package repository

import (
	"context"
	"fmt"
	"pokedex/database"
	"pokedex/internal/items/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const itemCollectionName = "items"

type ItemsRepository interface {
	SaveItem(ctx context.Context, ability model.ItemDetail) error
	GetItemByID(ctx context.Context, id int) (model.ItemDetail, error)
	GetItemList(ctx context.Context, limit, offset int, baseUrl string) ([]model.ListItemItem, int64, error)
}

type MongoItemsRepository struct {
	collection *mongo.Collection
}

func NewMongoItemsRepository() *MongoItemsRepository {
	return &MongoItemsRepository{
		collection: database.MongoDatabase.Collection(itemCollectionName),
	}
}

func (r *MongoItemsRepository) SaveItem(ctx context.Context, item model.ItemDetail) error {
	doc := model.ItemSaveDocument{
		ItemID:            item.ID,
		Name:              item.Name,
		Names:             item.Names,
		Cost:              item.Cost,
		Attributes:        item.Attributes,
		Category:          item.Category,
		BabyTriggerFor:    item.BabyTriggerFor,
		EffectEntries:     item.EffectEntries,
		FlavorTextEntries: item.FlavorTextEntries,
		FlingEffect:       item.FlingEffect,
		FlingPower:        item.FlingPower,
		GameIndicies:      item.GameIndicies,
		HeldByPokemon:     item.HeldByPokemon,
		Machine:           item.Machine,
		Sprites:           item.Sprites,
		LastSyncedAt:      time.Now().Unix(), // Menyimpan timestamp saat ini
	}

	filter := bson.M{"id": doc.ItemID} // Filter berdasarkan ID PokeAPI
	update := bson.M{"$set": doc}      // Menggunakan $set untuk memperbarui atau menyisipkan seluruh dokumen

	opts := options.Update().SetUpsert(true) // Opsi upsert: jika tidak ada, sisipkan; jika ada, perbarui.

	_, err := r.collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return fmt.Errorf("failed to save data %s (ID: %d) to MongoDB: %w", item.Name, item.ID, err)
	}
	return nil
}

func (r *MongoItemsRepository) GetItemByID(ctx context.Context, id int) (model.ItemDetail, error) {
	var doc model.ItemDetailDocument
	filter := bson.M{"id": id} // Mencari berdasarkan PokeAPI ID
	err := r.collection.FindOne(ctx, filter).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return model.ItemDetail{}, fmt.Errorf("items with ID %d not found", id)
		}
		return model.ItemDetail{}, fmt.Errorf("failed to retrieve items by ID from DB: %w", err)
	}
	return r.toDetail(doc), nil
}

func (r *MongoItemsRepository) GetItemList(ctx context.Context, limit, offset int, baseUrl string) ([]model.ListItemItem, int64, error) {
	totalCount, err := r.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count items in DB: %w", err)
	}

	findOptions := options.Find()
	findOptions.SetLimit(int64(limit))
	findOptions.SetSkip(int64(offset))
	findOptions.SetSort(bson.D{{Key: "id", Value: 1}}) // Sort by actual Pokemon ID

	cursor, err := r.collection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to retrieve items list from DB: %w", err)
	}
	defer cursor.Close(ctx)

	var docs []model.ItemListDocument
	if err = cursor.All(ctx, &docs); err != nil {
		return nil, 0, fmt.Errorf("failed to decode items list from DB: %w", err)
	}

	var listType []model.ListItemItem
	for _, doc := range docs {
		listType = append(listType, r.toPokedexList(doc, baseUrl))
	}

	return listType, totalCount, nil
}

func (r *MongoItemsRepository) toDetail(doc model.ItemDetailDocument) model.ItemDetail {
	return model.ItemDetail{
		ID:                doc.ItemID,
		Name:              doc.Name,
		Names:             doc.Names,
		Cost:              doc.Cost,
		Attributes:        doc.Attributes,
		Category:          doc.Category,
		BabyTriggerFor:    doc.BabyTriggerFor,
		EffectEntries:     doc.EffectEntries,
		FlavorTextEntries: doc.FlavorTextEntries,
		FlingEffect:       doc.FlingEffect,
		FlingPower:        doc.FlingPower,
		GameIndicies:      doc.GameIndicies,
		HeldByPokemon:     doc.HeldByPokemon,
		Machine:           doc.Machine,
		Sprites:           doc.Sprites,
	}
}

func (r *MongoItemsRepository) toPokedexList(doc model.ItemListDocument, baseUrl string) model.ListItemItem {
	return model.ListItemItem{
		ID:   doc.ID,
		Name: doc.Name,
		URL:  baseUrl + fmt.Sprintf("/%d", doc.ID),
	}
}
