package service

import (
	"RestGoTest/src/dto"
	"RestGoTest/src/model"
	"RestGoTest/src/repository"
	"context"
)

func AllProducts(ctx context.Context) ([]repository.Product, error) {
	return repository.GetProducts(ctx)
}

func FetchProduct(ctx context.Context, id int) (repository.Product, error) {
	var p repository.Product
	p.ID = id
	err := p.GetProduct(ctx)
	return p, err
}

func CreateProduct(ctx context.Context, p *repository.Product) (dto.CreateResponse, error) {
	err := p.CreateProduct(ctx)
	if err != nil {
		return dto.CreateResponse{}, err
	}
	return dto.CreateResponse{
		ID:          p.ID,
		ProductCode: p.ProductCode,
		Name:        p.Name,
	}, nil
}

func DeleteProduct(ctx context.Context, id int) error {
	p := repository.Product{Product: model.Product{ID: id}}
	return p.DeleteProduct(ctx)
}

func DeleteAllProducts(ctx context.Context) error {
	return repository.DeleteAllProducts(ctx)
}

func UpdateProduct(ctx context.Context, p *repository.Product) error {
	return p.UpdateProduct(ctx)
}
