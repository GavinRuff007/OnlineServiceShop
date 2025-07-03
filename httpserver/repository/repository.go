package repository

import (
	"RestGoTest/httpserver/model"
	"database/sql"
)

type Product struct {
	model.Product
}

func GetProducts(db *sql.DB) ([]Product, error) {
	rows, err := db.Query("SELECT * FROM products")
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

func (p *Product) GetProduct(db *sql.DB) error {
	return db.QueryRow("SELECT productCode, name, price, status, inventory FROM products WHERE ID = ?", p.ID).
		Scan(&p.ProductCode, &p.Name, &p.Price, &p.Status, &p.Inventory)
}

func (p *Product) CreateProduct(db *sql.DB) error {
	res, err := db.Exec("INSERT INTO products(productCode, name, price, status, inventory) VALUES(?,?,?,?,?)", p.ProductCode, p.Name, p.Price, p.Status, p.Inventory)
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
func (p *Product) DeleteProduct(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM products WHERE id = ?", p.ID)
	return err
}

func DeleteAllProducts(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM products")
	return err
}

func (p *Product) UpdateProduct(db *sql.DB) error {
	_, err := db.Exec(`
        UPDATE products
        SET productCode = ?, name = ?, price = ?, status = ?, inventory = ?
        WHERE id = ?
    `, p.ProductCode, p.Name, p.Price, p.Status, p.Inventory, p.ID)
	return err
}
