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
	GetPokemonTypeByID(ctx context.Context, id int) (model.PokemonTypeDetailResponse, error)
	GetPokemonTypeByName(ctx context.Context, name string) (model.PokemonTypeDetailResponse, error)
	GetPokemonTypeList(ctx context.Context, limit, offset int, baseUrl string) ([]model.PokemonTypeListItem, int64, error)
	GetWeaknessPokemonTypes(ctx context.Context, pokemonID int, pokemonTypes []string) ([]model.PokemonWeaknessTypes, error)
	GetPokemonByID(ctx context.Context, pokemonID int) (model.PokemonInfo, error)
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
		DamageRelations: pokemon_type.DamageRelations,
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

func (r *MongoPokemonTypeRepository) GetPokemonTypeByID(ctx context.Context, id int) (model.PokemonTypeDetailResponse, error) {
	var doc model.PokemonTypeDocument
	filter := bson.M{"id": id}
	err := r.collection.FindOne(ctx, filter).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return model.PokemonTypeDetailResponse{}, fmt.Errorf("data not found: %d", id)
		}
		return model.PokemonTypeDetailResponse{}, fmt.Errorf("failed to retrieve data: %w", err)
	}

	return r.toDetail(doc), nil
}

func (r *MongoPokemonTypeRepository) GetPokemonTypeByName(ctx context.Context, name string) (model.PokemonTypeDetailResponse, error) {
	var doc model.PokemonTypeDocument
	filter := bson.M{"name": name}
	err := r.collection.FindOne(ctx, filter).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return model.PokemonTypeDetailResponse{}, fmt.Errorf("data not found: %s", name)
		}
		return model.PokemonTypeDetailResponse{}, fmt.Errorf("failed to retrieve data: %w", err)
	}

	return r.toDetail(doc), nil
}

func (r *MongoPokemonTypeRepository) GetPokemonTypeList(ctx context.Context, limit, offset int, baseUrl string) ([]model.PokemonTypeListItem, int64, error) {
	totalCount, err := r.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count types in DB: %w", err)
	}

	findOptions := options.Find()
	findOptions.SetLimit(int64(limit))
	findOptions.SetSkip(int64(offset))
	findOptions.SetSort(bson.D{{Key: "id", Value: 1}}) // Sort by actual Pokemon ID

	cursor, err := r.collection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to retrieve types list from DB: %w", err)
	}
	defer cursor.Close(ctx)

	var docs []model.PokemonTypeListItemDocument
	if err = cursor.All(ctx, &docs); err != nil {
		return nil, 0, fmt.Errorf("failed to decode types list from DB: %w", err)
	}

	var listType []model.PokemonTypeListItem
	for _, doc := range docs {
		listType = append(listType, r.toDetailList(doc, baseUrl))
	}

	return listType, totalCount, nil
}

func (r *MongoPokemonTypeRepository) GetWeaknessPokemonTypes(ctx context.Context, pokemonID int, pokemonTypes []string) ([]model.PokemonWeaknessTypes, error) {
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: "id", Value: 1}})

	cursor, err := r.collection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve types list from DB: %w", err)
	}
	defer cursor.Close(ctx)

	var docs []model.PokemonTypeListItemDocument
	if err = cursor.All(ctx, &docs); err != nil {
		return nil, fmt.Errorf("failed to decode types list from DB: %w", err)
	}

	var listType []model.PokemonWeaknessTypes

	excludedTypes := map[string]bool{
		"stellar": true,
		"unknown": true,
		"shadow":  true,
	}

	for _, doc := range docs {
		// Kecualikan type dengan nama "stellar", "unknown", dan "shadow"
		if !excludedTypes[doc.Name] {
			listType = append(listType, r.toWeaknessTypes(doc))
		}
	}

	return listType, nil
}

func (r *MongoPokemonTypeRepository) GetPokemonByID(ctx context.Context, pokemonID int) (model.PokemonInfo, error) {
	var doc model.PokemonInfo

	filter := bson.M{"id": pokemonID}

	err := r.pokemonCollection.FindOne(ctx, filter).Decode(&doc)
	if err != nil {
		return model.PokemonInfo{}, err
	}

	return doc, nil
}

func (r *MongoPokemonTypeRepository) toWeaknessTypes(doc model.PokemonTypeListItemDocument) model.PokemonWeaknessTypes {

	res := model.PokemonWeaknessTypes{
		Name:          doc.Name,
		WeaknessPoint: 0,
	}

	return res
}

func (r *MongoPokemonTypeRepository) toDetailList(doc model.PokemonTypeListItemDocument, baseUrl string) model.PokemonTypeListItem {
	url := baseUrl + fmt.Sprintf("/%d", doc.TypeID)

	res := model.PokemonTypeListItem{
		TypeID: doc.TypeID,
		Name:   doc.Name,
		URL:    url,
	}

	return res
}

func (r *MongoPokemonTypeRepository) toDetail(doc model.PokemonTypeDocument) model.PokemonTypeDetailResponse {
	res := model.PokemonTypeDetailResponse{
		TypeID:          doc.TypeID,
		Name:            doc.Name,
		DamageRelations: doc.DamageRelations,
		MoveDamageClass: doc.MoveDamageClass,
	}

	return res
}
