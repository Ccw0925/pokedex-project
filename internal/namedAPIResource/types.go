package namedAPIResource

type NamedAPIResourceList struct {
	Count int `json:"count"`
	Next  string `json:"next"`
	Previous string `json:"previous"`
	Results []NamedAPIResource `json:"results"`
}

type NamedAPIResource struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}