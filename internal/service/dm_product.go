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
	var output = make([]dm.Product, 0, len(products))
	for _, product := range products {
		id, err := s.repo.Create(product)
		if err != nil {
			return nil, err
		}
		product.Id = id
		output = append(output, product)
	}
	return output, nil
}

func (s *DMProductService) Search(query dm.ProductSearchQuery) ([]dm.Product, error) {
	return s.repo.Search(query)
}

func (s *DMProductService) Exists(productID int) (bool, error) {
	return s.repo.Exists(productID)
}
