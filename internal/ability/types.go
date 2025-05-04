package ability

type Ability struct {
	Name          string        `json:"name"`
	EffectEntries []EffectEntry `json:"effect_entries"`
}

type EffectEntry struct {
	Effect   string `json:"effect"`
	Language struct {
		Name string `json:"name"`
	} `json:"language"`
}
