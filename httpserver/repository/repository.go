package repository

import (
	"RestGoTest/httpserver/database"
	"RestGoTest/httpserver/model"
	"database/sql"
)

type Product struct {
	model.Product
}

var DB *sql.DB = database.InitDatabase()

func GetProducts() ([]Product, error) {

	rows, err := DB.Query("SELECT * FROM products")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := []Product{}
	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.ID, &p.ProductCode, &p.Name, &p.Price, &p.Status, &p.Inventory); err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	return products, nil
}

func (p *Product) GetProduct() error {

	return DB.QueryRow("SELECT productCode, name, price, status, inventory FROM products WHERE ID = ?", p.ID).
		Scan(&p.ProductCode, &p.Name, &p.Price, &p.Status, &p.Inventory)
}

func (p *Product) CreateProduct() error {
	res, err := DB.Exec("INSERT INTO products(productCode, name, price, status, inventory) VALUES(?,?,?,?,?)", p.ProductCode, p.Name, p.Price, p.Status, p.Inventory)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	p.ID = int(id)
	return nil
}
func (p *Product) DeleteProduct() error {
	_, err := DB.Exec("DELETE FROM products WHERE id = ?", p.ID)
	return err
}

func DeleteAllProducts() error {
	_, err := DB.Exec("DELETE FROM products")
	return err
}

func (p *Product) UpdateProduct() error {
	_, err := DB.Exec(`
        UPDATE products
        SET productCode = ?, name = ?, price = ?, status = ?, inventory = ?
        WHERE id = ?
    `, p.ProductCode, p.Name, p.Price, p.Status, p.Inventory, p.ID)
	return err
}
