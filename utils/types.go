package utils

type GenderDistributionResult struct {
	Female string `json:"female"` // Contoh: "50.0%"
	Male   string `json:"male"`   // Contoh: "50.0%"
}

type BaseResourceNavigation struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	URL  string `json:"url"`
}
