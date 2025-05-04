package main

import (
	"time"

	"github.com/Ccw0925/pokedex-project/routes"
	"github.com/patrickmn/go-cache"
)

func main() {
	pokemonCache := cache.New(5*time.Minute, 10*time.Minute)

	r := routes.SetupRouter(pokemonCache)
	r.Run(":8080")
}
