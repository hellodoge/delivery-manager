package service

import (
	"github.com/hellodoge/delivery-manager/dm"
	"github.com/hellodoge/delivery-manager/internal/repository"
	"github.com/hellodoge/delivery-manager/pkg/response"
	"github.com/hellodoge/delivery-manager/pkg/util"
	"net/http"
)

type DMListService struct {
	repo repository.DMList
}

func NewDMListService(repo repository.DMList) *DMListService {
	return &DMListService{
		repo: repo,
	}
}

func (s *DMListService) Create(userID int, list dm.List) (dm.List, error) {
	listId, err := s.repo.Create(userID, list)
	if err != nil {
		return dm.List{}, err
	}
	list.Id = listId
	return list, nil
}

func (s *DMListService) GetUserLists(userID int) ([]dm.List, error) {
	lists, err := s.repo.GetUserLists(userID)
	if err != nil {
		return nil, err
	}
	return lists, nil
}

func (s *DMListService) Delete(userID, listID int) error {
	if err := s.ErrorIfNotOwner(userID, listID); err != nil {
		return err
	}
	return s.repo.Delete(listID)
}

func (s *DMListService) AddProduct(userID, listID int, index []dm.ProductIndex) error {
	if err := s.ErrorIfNotOwner(userID, listID); err != nil {
		return err
	}
	for _, position := range index {
		err := s.repo.AddProduct(listID, position.Id, position.Count)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *DMListService) ErrorIfNotOwner(userID, listID int) error {
	owners, err := s.repo.GetOwners(listID)
	if err != nil {
		return err
	}
	if !util.ContainsInt(userID, owners) {
		return response.ErrorResponseParameters{
			Message:    "You are not owner of requested list",
			StatusCode: http.StatusForbidden,
		}
	}
	return nil
}
