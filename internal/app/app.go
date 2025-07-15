package app

import (
	"pokedex/config"
	"pokedex/database"
	"pokedex/internal/shared/pokeapi"

	ability_handler "pokedex/internal/ability/handler"
	ability_repo "pokedex/internal/ability/repository"
	ability_service "pokedex/internal/ability/service"
	evolution_handler "pokedex/internal/evolution/handler"
	evolution_repo "pokedex/internal/evolution/repository"
	evolution_service "pokedex/internal/evolution/service"
	generation_handler "pokedex/internal/generation/handler"
	generation_repo "pokedex/internal/generation/repository"
	generation_service "pokedex/internal/generation/service"
	item_handler "pokedex/internal/items/handler"
	item_repo "pokedex/internal/items/repository"
	item_service "pokedex/internal/items/service"
	pokedexes_handler "pokedex/internal/pokedexes/handler"
	pokedexes_repo "pokedex/internal/pokedexes/repository"
	pokedexes_service "pokedex/internal/pokedexes/service"
	pokemon_species_handler "pokedex/internal/pokemon-species/handler"
	pokemon_species_repo "pokedex/internal/pokemon-species/repository"
	pokemon_species_service "pokedex/internal/pokemon-species/service"
	pokemon_type_handler "pokedex/internal/pokemon-type/handler"
	pokemon_type_repo "pokedex/internal/pokemon-type/repository"
	pokemon_type_service "pokedex/internal/pokemon-type/service"
	pokemon_handler "pokedex/internal/pokemon/handler"
	pokemon_repo "pokedex/internal/pokemon/repository"
	pokemon_service "pokedex/internal/pokemon/service"
)

// Application holds all initialized handlers
type Application struct {
	PokemonHandler        *pokemon_handler.PokemonHandler
	AbilityHandler        *ability_handler.AbilityHandler
	PokemonSpeciesHandler *pokemon_species_handler.PokemonSpeciesHandler
	EvolutionHandler      *evolution_handler.EvolutionHandler
	PokemonTypeHandler    *pokemon_type_handler.PokemonTypeHandler
	PokedexesHandler      *pokedexes_handler.PokedexesHandler
	ItemHandler           *item_handler.ItemHandler
	GenerationHandler     *generation_handler.GenerationHandler
}

// NewApplication initializes and returns a new Application instance
func NewApplication(cfg *config.Config, pokeAPIClient *pokeapi.Client) *Application {
	// Initialize Repositories
	evolutionRepo := evolution_repo.NewMongoEvolutionRepository()
	pokedexesRepo := pokedexes_repo.NewMongoPokedexesRepository()
	pokemonRepo := pokemon_repo.NewMongoPokemonRepository()
	abilityRepo := ability_repo.NewMongoAbilityRepository()
	itemRepo := item_repo.NewMongoItemsRepository()
	generationRepo := generation_repo.NewMongoGenerationRepository()
	pokemonSpeciesRepo := pokemon_species_repo.NewMongoPokemonSpeciesRepository()
	pokemonTypeRepo := pokemon_type_repo.NewMongoPokemonTypeRepository()

	// Initialize Services
	evolutionService := evolution_service.NewEvolutionService(evolutionRepo, pokeAPIClient)
	pokedexesService := pokedexes_service.NewPokedexesService(pokedexesRepo, pokeAPIClient)
	pokemonService := pokemon_service.NewPokemonService(pokemonRepo, pokeAPIClient, evolutionService, pokedexesService)
	abilityService := ability_service.NewAbilityService(abilityRepo, pokeAPIClient)
	itemService := item_service.NewItemsService(itemRepo, pokeAPIClient)
	generationService := generation_service.NewGenerationService(generationRepo, pokeAPIClient)
	pokemonSpeciesService := pokemon_species_service.NewPokemonSpeciesService(pokemonSpeciesRepo, pokeAPIClient)
	pokemonTypeService := pokemon_type_service.NewPokemonTypeService(pokemonTypeRepo, pokeAPIClient)

	// Initialize Handlers
	evolutionHandler := evolution_handler.NewEvolutionHandler(evolutionService)
	pokedexesHandler := pokedexes_handler.NewPokedexesHandler(pokedexesService)
	pokemonHandler := pokemon_handler.NewPokemonHandler(pokemonService)
	abilityHandler := ability_handler.NewAbilityHandler(abilityService)
	itemHandler := item_handler.NewItemHandler(itemService)
	generationHandler := generation_handler.NewGenerationHandler(generationService)
	pokemonSpeciesHandler := pokemon_species_handler.NewPokemonSpeciesHandler(pokemonSpeciesService)
	pokemonTypeHandler := pokemon_type_handler.NewPokemonTypeHandler(pokemonTypeService)

	return &Application{
		PokemonHandler:        pokemonHandler,
		AbilityHandler:        abilityHandler,
		PokemonSpeciesHandler: pokemonSpeciesHandler,
		EvolutionHandler:      evolutionHandler,
		PokemonTypeHandler:    pokemonTypeHandler,
		PokedexesHandler:      pokedexesHandler,
		ItemHandler:           itemHandler,
		GenerationHandler:     generationHandler,
	}
}

// InitDatabaseAndPokeAPI initializes the database connection and PokeAPI client.
// It returns the PokeAPI client and a cleanup function to be deferred.
func InitDatabaseAndPokeAPI(cfg *config.Config) (*pokeapi.Client, func()) {
	database.ConnectDB(cfg)
	pokeAPIClient := pokeapi.NewClient(cfg)
	return pokeAPIClient, func() {
		database.DisconnectDB()
		pokeAPIClient.CloseClient()
	}
}
