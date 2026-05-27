package model

type CategoryLookup struct {
	ID   string `json:"id" bson:"_id"`
	Name string `json:"name" bson:"name"`
}

type OrderFilters struct {
	Status []string `json:"status"`
	Type   []string `json:"type"`
}

type FilterConfig struct {
	Orders OrderFilters `json:"orders"`
}

type SortConfig struct {
	Categories []string `json:"categories"`
	Menu       []string `json:"menu"`
	Orders     []string `json:"orders"`
}

type EnumConfig struct {
	Filters FilterConfig `json:"filters"`
	Sorts   SortConfig   `json:"sorts"`
}
