package model

type Product struct {
	ID          int    `json:"id"`
	ProductCode string `json:"productCode"`
	Name        string `json:"name"`
	Price       int    `json:"price"`
	Status      string `json:"status"`
	Inventory   string `json:"inventory"`
}
