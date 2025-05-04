package pokemon

type Pokemon struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Height int    `json:"height"`
	Weight int    `json:"weight"`
	Types  []struct {
		Type struct {
			Name string `json:"name"`
		} `json:"type"`
	} `json:"types"`
	Abilities []struct {
		Ability struct {
			Name string `json:"name"`
		} `json:"ability"`
	} `json:"abilities"`
	Sprites struct {
		Other struct {
			OfficialArtwork struct {
				FrontDefault string `json:"front_default"`
			} `json:"official-artwork"`
		} `json:"other"`
	} `json:"sprites"`
}

type PokemonSpecies struct {
	Id int `json:"id"`
	Name string `json:"name"`
	EvolutionChain struct {
		Url string `json:"url"`
	} `json:"evolution_chain"`
}

type EvolutionChain struct {
	Chain *EvolutionDetails `json:"chain"`
}

type EvolutionDetails struct {
	Species struct {
		Name string `json:"name"`
	} `json:"species"`
	EvolvesTo []EvolutionDetails `json:"evolves_to"`
}
