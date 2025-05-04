package routes

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"encoding/json"

	"github.com/Ccw0925/pokedex-project/internal/ability"
	"github.com/Ccw0925/pokedex-project/internal/namedAPIResource"
	"github.com/Ccw0925/pokedex-project/internal/pokemon"
	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func SetupRouter(pokemonCache *cache.Cache) *gin.Engine {
	r := gin.Default()

	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "success",
		})
	})

	// List all pokemons
	r.GET("/pokemon", func(c *gin.Context) {
		limitStr := c.DefaultQuery("limit", "20")
		offsetStr := c.DefaultQuery("offset", "0")

		// Convert limit and offset to integers
		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit parameter"})
			return
		}

		offset, err := strconv.Atoi(offsetStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid offset parameter"})
			return
		}

		p, err := fetchAllPokemon(limitStr, offsetStr, pokemonCache)

		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}

		pokemonList := make([]gin.H, len(p.Results))
		for i, pokemon := range p.Results {
			pokemonID := offset + i + 1
			pokemonList[i] = gin.H{
				"id":   pokemonID,
				"name": pokemon.Name,
			}
		}

		var next interface{}
		if p.Next != "" {
			next = fmt.Sprintf("/pokemon?limit=%s&offset=%d", limitStr, offset+limit)
		}

		var prev interface{}
		if p.Previous != "" {
			prev = fmt.Sprintf("/pokemon?limit=%s&offset=%d", limitStr, offset-limit)
		}

		c.JSON(http.StatusOK, gin.H{
			"count":    p.Count,
			"next":     next,
			"previous": prev,
			"pokemons": pokemonList,
		})
	})

	// Get a single pokemon by identifier
	r.GET("/pokemon/:identifier", func(c *gin.Context) {
		identifier := strings.ToLower(c.Param("identifier"))
		p, err := fetchPokemon(identifier, pokemonCache)

		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}

		titleName := cases.Title(language.English, cases.NoLower).String(p.Name)

		c.JSON(http.StatusOK, gin.H{
			"id":        p.ID,
			"name":      titleName,
			"height":    float64(p.Height) / 10,
			"weight":    float64(p.Weight) / 10,
			"types":     p.Types,
			"abilities": p.Abilities,
			"imageUrl":  p.Sprites.Other.OfficialArtwork.FrontDefault,
		})
	})

	// Get abilities for a single pokemon
	r.GET("/pokemon/:identifier/abilities", func(c *gin.Context) {
		identifier := strings.ToLower(c.Param("identifier"))
		p, err := fetchPokemon(identifier, pokemonCache)

		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}

		var abilities []ability.Ability
		for _, a := range p.Abilities {
			pokemonAbility, err := fetchAbility(a.Ability.Name, pokemonCache)
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{
					"error": err.Error(),
				})
				return
			}

			// Filter English effect entries
			var englishEffects []ability.EffectEntry
			for _, e := range pokemonAbility.EffectEntries {
				if e.Language.Name == "en" {
					englishEffects = append(englishEffects, e)
				}
			}

			pokemonAbility.EffectEntries = englishEffects
			abilities = append(abilities, *pokemonAbility)
		}

		c.JSON(http.StatusOK, gin.H{
			"abilities": abilities,
		})
	})

	// Get pokemon evolution chains
	r.GET("/pokemon/:identifier/evolutions", func(c *gin.Context) {
		identifier := strings.ToLower(c.Param("identifier"))
		s, err := fetchPokemonSpecies(identifier, pokemonCache)

		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}

		e, err := fetchPokemonEvolutionChain(s.EvolutionChain.Url, pokemonCache)

		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"evolutions": e.Chain,
		})
		
	})

	return r
}

func fetchAllPokemon(limit string, offset string, pokemonCache *cache.Cache) (*namedAPIResource.NamedAPIResourceList, error) {
	cacheKey := "list_pokemon_" + fmt.Sprintf("%s_%s", limit, offset)

	// Check cache first
	if cached, found := pokemonCache.Get(cacheKey); found {
		return cached.(*namedAPIResource.NamedAPIResourceList), nil
	}

	fmt.Println("No cache found. Fetching...")

	resp, err := http.Get("https://pokeapi.co/api/v2/pokemon?limit=" + limit + "&offset=" + offset)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch Pokémon list: %v", err)
	}
	defer resp.Body.Close()

	var p namedAPIResource.NamedAPIResourceList
	err = json.NewDecoder(resp.Body).Decode(&p)

	pokemonCache.Set(cacheKey, &p, cache.DefaultExpiration)

	return &p, err
}

func fetchPokemon(identifier string, pokemonCache *cache.Cache) (*pokemon.Pokemon, error) {
	// Check cache first
	if cached, found := pokemonCache.Get(identifier); found {
		return cached.(*pokemon.Pokemon), nil
	}

	fmt.Println("No cache found. Fetching...")

	resp, err := http.Get("https://pokeapi.co/api/v2/pokemon/" + identifier)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch Pokémon: %v", err)
	}
	defer resp.Body.Close()

	var p pokemon.Pokemon
	err = json.NewDecoder(resp.Body).Decode(&p)

	// Store in cache
	pokemonCache.Set(fmt.Sprintf("%d", p.ID), &p, cache.DefaultExpiration)
	pokemonCache.Set(p.Name, &p, cache.DefaultExpiration)

	return &p, err
}

func fetchAbility(identifier string, pokemonCache *cache.Cache) (*ability.Ability, error) {
	cacheKey := "ability_" + identifier

	// Check cache first
	if cached, found := pokemonCache.Get(cacheKey); found {
		return cached.(*ability.Ability), nil
	}

	fmt.Println("No cache found. Fetching...")

	resp, err := http.Get("https://pokeapi.co/api/v2/ability/" + identifier)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch ability: %v", err)
	}
	defer resp.Body.Close()

	var a ability.Ability
	err = json.NewDecoder(resp.Body).Decode(&a)

	// Store in cache
	pokemonCache.Set(cacheKey, &a, cache.DefaultExpiration)

	return &a, err
}

func fetchPokemonSpecies(identifier string, pokemonCache *cache.Cache) (*pokemon.PokemonSpecies, error) {
	cacheKey := "pokemon_species_" + identifier

	// Check cache first
	if cached, found := pokemonCache.Get(cacheKey); found {
		return cached.(*pokemon.PokemonSpecies), nil
	}

	fmt.Println("No cache found. Fetching...")

	resp, err := http.Get("https://pokeapi.co/api/v2/pokemon-species/" + identifier)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch Pokémon species: %v", err)
	}
	defer resp.Body.Close()

	var p pokemon.PokemonSpecies
	err = json.NewDecoder(resp.Body).Decode(&p)

	// Store in cache
	pokemonCache.Set(cacheKey, &p, cache.DefaultExpiration)

	return &p, err
}

func fetchPokemonEvolutionChain(url string, pokemonCache *cache.Cache) (*pokemon.EvolutionChain, error) {
	cacheKey := "evolution_chain_" + url

	// Check cache first
	if cached, found := pokemonCache.Get(cacheKey); found {
		return cached.(*pokemon.EvolutionChain), nil
	}

	fmt.Println("No cache found. Fetching...")

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch Pokémon evolution chain: %v", err)
	}
	defer resp.Body.Close()

	var e pokemon.EvolutionChain
	err = json.NewDecoder(resp.Body).Decode(&e)

	// Store in cache
	pokemonCache.Set(cacheKey, &e, cache.DefaultExpiration)

	return &e, err
}
