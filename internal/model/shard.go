package model

// Shard is one data file: the chips of one manufacturer in one category.
type Shard struct {
	Manufacturer string  `json:"manufacturer"`
	Category     string  `json:"category"`
	Chips        []*Chip `json:"chips"`
}

// ShardRef points to a shard file from the index.
type ShardRef struct {
	Path     string `json:"path"`
	Category string `json:"category"`
	Count    int    `json:"count"`
	Rev      string `json:"rev"`
}
