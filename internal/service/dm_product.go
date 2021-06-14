package service

import (
	"github.com/hellodoge/delivery-manager/dm"
	"github.com/hellodoge/delivery-manager/internal/repository"
)

type DMProductService struct {
	repo repository.DMProduct
}

func NewDMProductService(repo repository.DMProduct) *DMProductService {
	return &DMProductService{
		repo: repo,
	}
}

func (s *DMProductService) Create(products []dm.Product) ([]dm.Product, error) {
	ids, err := s.repo.Create(products)
	if err != nil {
		return nil, err
	}
	for i := range products {
		products[i].Id = ids[i]
	}
	return products, nil
}

func (s *DMProductService) Search(query dm.ProductSearchQuery) ([]dm.Product, error) {
	return s.repo.Search(query)
}

func (s *DMProductService) Exists(productID int) (bool, error) {
	return s.repo.Exists(productID)
}
