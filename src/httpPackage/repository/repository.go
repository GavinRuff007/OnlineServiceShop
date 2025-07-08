package repository

import (
	"RestGoTest/src/database"
	"RestGoTest/src/model"
	"context"
	"database/sql"
)

type Product struct {
	model.Product
}

var DB *sql.DB = database.InitDatabase()

func GetProducts(ctx context.Context) ([]Product, error) {
	rows, err := DB.QueryContext(ctx, "SELECT * FROM products")
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

func (p *Product) GetProduct(ctx context.Context) error {
	return DB.QueryRowContext(ctx, "SELECT productCode, name, price, status, inventory FROM products WHERE ID = ?", p.ID).
		Scan(&p.ProductCode, &p.Name, &p.Price, &p.Status, &p.Inventory)
}

func (p *Product) CreateProduct(ctx context.Context) error {
	res, err := DB.ExecContext(ctx, "INSERT INTO products(productCode, name, price, status, inventory) VALUES(?,?,?,?,?)",
		p.ProductCode, p.Name, p.Price, p.Status, p.Inventory)
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

func (p *Product) DeleteProduct(ctx context.Context) error {
	_, err := DB.ExecContext(ctx, "DELETE FROM products WHERE id = ?", p.ID)
	return err
}

func DeleteAllProducts(ctx context.Context) error {
	_, err := DB.ExecContext(ctx, "DELETE FROM products")
	return err
}

func (p *Product) UpdateProduct(ctx context.Context) error {
	_, err := DB.ExecContext(ctx, `
        UPDATE products
        SET productCode = ?, name = ?, price = ?, status = ?, inventory = ?
        WHERE id = ?
    `, p.ProductCode, p.Name, p.Price, p.Status, p.Inventory, p.ID)
	return err
}
