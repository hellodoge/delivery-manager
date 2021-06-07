package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hellodoge/delivery-manager/dm"
	"github.com/hellodoge/delivery-manager/pkg/response"
	"net/http"
)

func (h *Handler) createProducts(ctx *gin.Context) {
	var input []dm.Product
	if err := ctx.BindJSON(&input); err != nil {
		newErrorResponse(ctx, response.ErrorResponseParameters{
			Internal:   err,
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
		})
		return
	}
	products, err := h.services.DMProduct.Create(input)
	if err != nil {
		newErrorResponse(ctx, err)
	}
	ctx.JSON(http.StatusOK, products)
}

func (h *Handler) searchForProducts(ctx *gin.Context) {
	var input dm.ProductSearchQuery
	if err := ctx.BindJSON(&input); err != nil {
		newErrorResponse(ctx, response.ErrorResponseParameters{
			Internal:   err,
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
		})
		return
	}
	result, err := h.services.DMProduct.Search(input)
	if err != nil {
		newErrorResponse(ctx, response.ErrorResponseParameters{
			Internal:   err,
			IsInternal: true,
		})
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (h *Handler) getProducts(ctx *gin.Context) {

}

type AddProductsInput struct {
	ListID   int               `json:"list_id" binding:"required"`
	Products []dm.ProductIndex `json:"products"`
}

func (h *Handler) addProducts(ctx *gin.Context) {
	userId, err := getUserId(ctx)
	if err != nil {
		newErrorResponse(ctx, err)
		return
	}

	var input AddProductsInput
	if err := ctx.BindJSON(&input); err != nil {
		newErrorResponse(ctx, response.ErrorResponseParameters{
			Internal:   err,
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
		})
		return
	}

	for _, product := range input.Products {
		exist, err := h.services.DMProduct.Exists(product.Id)
		if err != nil {
			newErrorResponse(ctx, err)
			return
		}
		if !exist {
			newErrorResponse(ctx, response.ErrorResponseParameters{
				Message:    fmt.Sprintf("Product %d does not exist", product.Id),
				StatusCode: http.StatusBadRequest,
			})
		}
	}

	err = h.services.DMList.AddProduct(userId, input.ListID, input.Products)
	if err != nil {
		newErrorResponse(ctx, err)
		return
	}

	ctx.Status(http.StatusOK)
}

func (h *Handler) removeProducts(ctx *gin.Context) {

}

func (h *Handler) updateProducts(ctx *gin.Context) {

}
