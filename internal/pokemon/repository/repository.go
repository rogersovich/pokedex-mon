package repository

import (
	"context"
	"fmt"
	"sort"
	"time"

	"pokedex/database"
	pokemon_species_model "pokedex/internal/pokemon-species/model"
	pokemon_model "pokedex/internal/pokemon/model"
	"pokedex/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const pokemonCollectionName = "pokemons"

// PokemonRepository defines the interface for persisting and retrieving Pokemon data.
type PokemonRepository interface {
	SavePokemon(ctx context.Context, pokemon pokemon_model.PokemonDetail) error
	GetPokemonByID(ctx context.Context, id int) (pokemon_model.PokemonDetailResponse, error)
	GetPokemonByName(ctx context.Context, name string) (pokemon_model.PokemonDetailResponse, error)
	GetPokemonList(ctx context.Context, limit, offset int) ([]pokemon_model.PokemonDetail, int64, error)
	SearchPokemons(ctx context.Context, query string, limit, offset int) ([]pokemon_model.PokemonDetail, int64, error)
}

// MongoPokemonRepository implements the PokemonRepository interface for MongoDB.
type MongoPokemonRepository struct {
	collection        *mongo.Collection
	collectionSpecies *mongo.Collection
}

// NewMongoPokemonRepository creates a new MongoDB repository.
func NewMongoPokemonRepository() *MongoPokemonRepository {
	return &MongoPokemonRepository{
		collection:        database.MongoDatabase.Collection(pokemonCollectionName),
		collectionSpecies: database.MongoDatabase.Collection("pokemon-species"),
	}
}

// SavePokemon saves/updates a Pokemon in MongoDB.
func (r *MongoPokemonRepository) SavePokemon(ctx context.Context, pokemon pokemon_model.PokemonDetail) error {
	doc := pokemon_model.PokemonDocument{
		PokemonID:              pokemon.ID,
		Name:                   pokemon.Name,
		Height:                 pokemon.Height,
		Weight:                 pokemon.Weight,
		BaseExperience:         pokemon.BaseExperience,
		Sprites:                pokemon.Sprites,
		Types:                  pokemon.Types,
		Stats:                  pokemon.Stats,
		Abilities:              pokemon.Abilities,
		Forms:                  pokemon.Forms,
		GameIndices:            pokemon.GameIndices,
		HeldItems:              pokemon.HeldItems,
		IsDefault:              pokemon.IsDefault,
		LocationAreaEncounters: pokemon.LocationAreaEncounters,
		Moves:                  pokemon.Moves,
		Order:                  pokemon.Order,
		Species:                pokemon.Species,
		LastSyncedAt:           time.Now().Unix(),
	}

	filter := bson.M{"id": doc.PokemonID}
	update := bson.M{"$set": doc}
	opts := options.Update().SetUpsert(true)

	_, err := r.collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return fmt.Errorf("failed to save pokemon %s (ID: %d) to MongoDB: %w", pokemon.Name, pokemon.ID, err)
	}
	return nil
}

func (r *MongoPokemonRepository) GetPokemonByID(ctx context.Context, id int) (pokemon_model.PokemonDetailResponse, error) {
	var doc pokemon_model.PokemonDocument
	filter := bson.M{"id": id}
	err := r.collection.FindOne(ctx, filter).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return pokemon_model.PokemonDetailResponse{}, fmt.Errorf("pokemon not found: %d", id)
		}
		return pokemon_model.PokemonDetailResponse{}, fmt.Errorf("failed to retrieve pokemon by ID from DB: %w", err)
	}

	var speciesName = doc.Species.Name

	var docSpecies pokemon_species_model.PokemonSpeciesDocument
	filter = bson.M{"name": speciesName}
	err = r.collectionSpecies.FindOne(ctx, filter).Decode(&docSpecies)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return pokemon_model.PokemonDetailResponse{}, fmt.Errorf("pokemon species not found: %s", speciesName)
		}
		return pokemon_model.PokemonDetailResponse{}, fmt.Errorf("failed to retrieve pokemon species by id from DB: %w", err)
	}

	return r.toDetailResponse(doc, docSpecies), nil
}

func (r *MongoPokemonRepository) GetPokemonByName(ctx context.Context, name string) (pokemon_model.PokemonDetailResponse, error) {
	var doc pokemon_model.PokemonDocument
	filter := bson.M{"name": name}
	err := r.collection.FindOne(ctx, filter).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return pokemon_model.PokemonDetailResponse{}, fmt.Errorf("pokemon not found: %s", name)
		}
		return pokemon_model.PokemonDetailResponse{}, fmt.Errorf("failed to retrieve pokemon by name from DB: %w", err)
	}

	var speciesName = doc.Species.Name

	var docSpecies pokemon_species_model.PokemonSpeciesDocument
	filter = bson.M{"name": speciesName}
	err = r.collectionSpecies.FindOne(ctx, filter).Decode(&docSpecies)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return pokemon_model.PokemonDetailResponse{}, fmt.Errorf("pokemon species not found: %s", name)
		}
		return pokemon_model.PokemonDetailResponse{}, fmt.Errorf("failed to retrieve pokemon species by name from DB: %w", err)
	}

	return r.toDetailResponse(doc, docSpecies), nil
}

func (r *MongoPokemonRepository) GetPokemonList(ctx context.Context, limit, offset int) ([]pokemon_model.PokemonDetail, int64, error) {
	totalCount, err := r.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count pokemons in DB: %w", err)
	}

	findOptions := options.Find()
	findOptions.SetLimit(int64(limit))
	findOptions.SetSkip(int64(offset))
	findOptions.SetSort(bson.D{{Key: "id", Value: 1}}) // Sort by actual Pokemon ID

	cursor, err := r.collection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to retrieve pokemon list from DB: %w", err)
	}
	defer cursor.Close(ctx)

	var pokemonDocs []pokemon_model.PokemonDocument
	if err = cursor.All(ctx, &pokemonDocs); err != nil {
		return nil, 0, fmt.Errorf("failed to decode pokemon list from DB: %w", err)
	}

	var pokemonDetails []pokemon_model.PokemonDetail
	for _, doc := range pokemonDocs {
		pokemonDetails = append(pokemonDetails, r.toDetail(doc))
	}

	return pokemonDetails, totalCount, nil
}

func (r *MongoPokemonRepository) SearchPokemons(ctx context.Context, query string, limit, offset int) ([]pokemon_model.PokemonDetail, int64, error) {
	// Buat filter regex untuk pencarian substring case-insensitive
	filter := bson.M{
		"name": bson.M{
			"$regex":   query,
			"$options": "i", // "i" for case-insensitive
		},
	}

	totalCount, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count search results in DB: %w", err)
	}

	findOptions := options.Find()
	findOptions.SetLimit(int64(limit))
	findOptions.SetSkip(int64(offset))
	findOptions.SetSort(bson.D{{Key: "id", Value: 1}}) // Sort by actual Pokemon ID

	cursor, err := r.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to search pokemons in DB: %w", err)
	}
	defer cursor.Close(ctx)

	var pokemonDocs []pokemon_model.PokemonDocument
	if err = cursor.All(ctx, &pokemonDocs); err != nil {
		return nil, 0, fmt.Errorf("failed to decode search results from DB: %w", err)
	}

	var pokemonDetails []pokemon_model.PokemonDetail
	for _, doc := range pokemonDocs {
		pokemonDetails = append(pokemonDetails, r.toDetail(doc))
	}

	return pokemonDetails, totalCount, nil
}

// toDetail converts a PokemonDocument to a pokemon_model.PokemonDetail.
func (r *MongoPokemonRepository) toDetail(doc pokemon_model.PokemonDocument) pokemon_model.PokemonDetail {
	return pokemon_model.PokemonDetail{
		ID:                     doc.PokemonID,
		Name:                   doc.Name,
		Height:                 doc.Height,
		Weight:                 doc.Weight,
		BaseExperience:         doc.BaseExperience,
		Sprites:                doc.Sprites,
		Types:                  doc.Types,
		Stats:                  doc.Stats,
		Abilities:              doc.Abilities,
		Forms:                  doc.Forms,
		GameIndices:            doc.GameIndices,
		HeldItems:              doc.HeldItems,
		IsDefault:              doc.IsDefault,
		LocationAreaEncounters: doc.LocationAreaEncounters,
		Moves:                  doc.Moves,
		Order:                  doc.Order,
		Species:                doc.Species,
	}
}

func (r *MongoPokemonRepository) toDetailResponse(
	doc pokemon_model.PokemonDocument,
	docSpecies pokemon_species_model.PokemonSpeciesDocument) pokemon_model.PokemonDetailResponse {

	thumbnailImg := utils.GetThumbnailPokemon(doc.PokemonID)

	eggGroups := make([]pokemon_model.ResourceReference, len(docSpecies.EggGroups))
	for i, eg := range docSpecies.EggGroups {
		eggGroups[i] = pokemon_model.ResourceReference{
			Name: eg.Name,
			URL:  eg.URL,
		}
	}

	calcStats := make([]pokemon_model.PokemonStatFull, len(doc.Stats))
	for i, stat := range doc.Stats {
		calcStats[i] = pokemon_model.PokemonStatFull{
			BaseStat: stat.BaseStat,
			Effort:   0,
			StatName: stat.Stat.Name,
			MinStat:  utils.CalcMinStat(stat.BaseStat, stat.Stat.Name),
			MaxStat:  utils.CalcMaxStat(stat.BaseStat, stat.Stat.Name),
		}
	}

	otherNames := make([]pokemon_model.PokemonOtherNames, len(docSpecies.Names))
	for i, name := range docSpecies.Names {
		otherNames[i] = pokemon_model.PokemonOtherNames{
			Name:     name.Name,
			Language: name.Language.Name,
		}
	}

	allowedMoveVersion := []string{"black-white", "red-blue"}
	allowedMoveMethod := []string{"egg", "level-up", "machine", "tutor"}

	pokedexNumbers := []pokemon_model.PokemonNumber{}
	for _, number := range docSpecies.PokedexNumbers {
		pokedexNumbers = append(pokedexNumbers, pokemon_model.PokemonNumber{
			EntryNumber: number.EntryNumber,
			Pokedex:     pokemon_model.ResourceReference(number.Pokedex),
		})
	}

	// Mendapatkan ID dari URL evolution chain, misal: https://pokeapi.co/api/v2/evolution-chain/67/
	evolutionURL := docSpecies.EvolutionChain.URL
	var evolutionID int
	fmt.Sscanf(evolutionURL, "https://pokeapi.co/api/v2/evolution-chain/%d/", &evolutionID)

	return pokemon_model.PokemonDetailResponse{
		ID:           doc.PokemonID,
		Name:         doc.Name,
		Height:       doc.Height,
		Weight:       doc.Weight,
		Thumbnail:    thumbnailImg,
		Order:        doc.Order,
		Habitat:      docSpecies.Habitat.Name,
		Abilities:    doc.Abilities,
		Types:        doc.Types,
		Stats:        calcStats,
		GroupedMoves: GroupMovesByVersion(doc.Moves, allowedMoveVersion, allowedMoveMethod),
		Sprites:      doc.Sprites,
		OtherNames:   otherNames,
		Training: pokemon_model.PokemonTraining{
			CaptureRate:        docSpecies.CaptureRate,
			CaptureRatePercent: utils.CalcCaptureRate(docSpecies.CaptureRate, 1.0, 1.0),
			BaseExperience:     doc.BaseExperience,
			BaseHappiness:      docSpecies.BaseHappiness,
			GrowthRate:         utils.ConvertGrowthRate(docSpecies.GrowthRate.Name),
		},
		Breeding: pokemon_model.PokemonBreeding{
			EggGroups:    eggGroups,
			GenderRate:   utils.CalcGenderDistribution(docSpecies.GenderRate),
			HatchCounter: docSpecies.HatchCounter,
			EggCycles:    utils.CalcEggCycles(docSpecies.HatchCounter),
		},
		IsBaby:         docSpecies.IsBaby,
		IsLegendary:    docSpecies.IsLegendary,
		IsMythical:     docSpecies.IsMythical,
		Color:          pokemon_model.ResourceReference(docSpecies.Color),
		Generation:     pokemon_model.ResourceReference(docSpecies.Generation),
		PokedexNumbers: pokedexNumbers,
		EvolutionID:    evolutionID,
	}
}

func isVersionAllowed(versionName string, allowedVersions []string) bool {
	if len(allowedVersions) == 0 { // Jika daftar filter kosong, izinkan semua versi
		return true
	}
	for _, allowed := range allowedVersions {
		if versionName == allowed {
			return true // Ditemukan di daftar yang diizinkan
		}
	}
	return false // Tidak ditemukan di daftar yang diizinkan
}

func isMethodAllowed(methodName string, allowedMethods []string) bool {
	if len(allowedMethods) == 0 {
		return true
	}
	for _, allowed := range allowedMethods {
		if methodName == allowed {
			return true
		}
	}
	return false
}

func GroupMovesByVersion(pokemonMoves []pokemon_model.PokemonMoves, allowedVersions []string, allowedMethods []string) []pokemon_model.GroupedVersionMoves {
	// Langkah 1: Kumpulkan data ke dalam map sementara berdasarkan `versionGroupName`
	// Ini menyimpan semua gerakan (belum dikelompokkan oleh metode) untuk setiap versi.
	tempGroupedByVersionMap := make(map[string][]pokemon_model.GroupedMoveInfo)

	for _, moveData := range pokemonMoves {
		moveName := moveData.Move.Name
		moveURL := moveData.Move.URL

		for _, detail := range moveData.VersionGroupDetails {
			versionGroupName := detail.VersionGroup.Name

			if !isVersionAllowed(versionGroupName, allowedVersions) {
				continue
			}

			methodName := detail.MoveLearnMethod.Name
			// Filter berdasarkan metode belajar
			if !isMethodAllowed(methodName, allowedMethods) {
				continue
			}

			info := pokemon_model.GroupedMoveInfo{
				MoveName:        moveName,
				MoveURL:         moveURL,
				LevelLearnedAt:  detail.LevelLearnedAt,
				MoveLearnMethod: methodName, // Gunakan methodName yang sudah difilter
				Order:           detail.Order,
			}
			tempGroupedByVersionMap[versionGroupName] = append(tempGroupedByVersionMap[versionGroupName], info)
		}
	}

	// Langkah 2: Ubah map sementara menjadi slice GroupedVersionMoves,
	// sambil melakukan pengelompokan internal berdasarkan `move_learn_method`.
	var finalResult []pokemon_model.GroupedVersionMoves
	for groupName, movesInVersion := range tempGroupedByVersionMap {
		// Buat map sementara untuk mengelompokkan gerakan berdasarkan metode belajar
		tempGroupedByMethodMap := make(map[string][]pokemon_model.GroupedMoveInfo)
		for _, moveInfo := range movesInVersion {
			tempGroupedByMethodMap[moveInfo.MoveLearnMethod] = append(tempGroupedByMethodMap[moveInfo.MoveLearnMethod], moveInfo)
		}

		// Konversi map metode belajar ke slice MovesByLearnMethod
		var movesByMethodSlice []pokemon_model.MovesByLearnMethod
		for methodName, moves := range tempGroupedByMethodMap {
			// Opsional: Urutkan gerakan dalam setiap metode berdasarkan move_name atau level_learned_at
			sort.Slice(moves, func(i, j int) bool {
				// Urutkan berdasarkan level, jika sama, urutkan berdasarkan nama
				if moves[i].LevelLearnedAt != moves[j].LevelLearnedAt {
					return moves[i].LevelLearnedAt < moves[j].LevelLearnedAt
				}
				return moves[i].MoveName < moves[j].MoveName
			})
			movesByMethodSlice = append(movesByMethodSlice, pokemon_model.MovesByLearnMethod{
				MethodName: methodName,
				Moves:      moves,
			})
		}

		// Opsional: Urutkan metode belajar (misal: "level-up" sebelum "machine")
		sort.Slice(movesByMethodSlice, func(i, j int) bool {
			return movesByMethodSlice[i].MethodName < movesByMethodSlice[j].MethodName
		})

		finalResult = append(finalResult, pokemon_model.GroupedVersionMoves{
			GroupName:     groupName,
			MovesByMethod: movesByMethodSlice, // Gunakan slice yang sudah dikelompokkan berdasarkan metode
		})
	}

	// Opsional: Urutkan hasil akhir berdasarkan GroupName agar output konsisten
	sort.Slice(finalResult, func(i, j int) bool {
		return finalResult[i].GroupName < finalResult[j].GroupName
	})

	return finalResult
}
