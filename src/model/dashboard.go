package model

type DashboardStat struct {
	ID         string `json:"id"`
	Label      string `json:"label"`
	Value      string `json:"value"`
	Trend      string `json:"trend"`
	IsPositive bool   `json:"isPositive"`
}

type RevenueData struct {
	Day   string  `json:"day"`
	Value float64 `json:"value"`
}

type CategoryData struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

type DashboardSummary struct {
	Stats           []DashboardStat `json:"stats"`
	RevenueOverview []RevenueData   `json:"revenueOverview"`
	TopCategories   []CategoryData  `json:"topCategories"`
}
