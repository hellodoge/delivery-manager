package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/hellodoge/delivery-manager/dm"
	"github.com/hellodoge/delivery-manager/pkg/response"
	"net/http"
	"strconv"
)

func (h *Handler) createList(ctx *gin.Context) {
	var input dm.List
	if err := ctx.BindJSON(&input); err != nil {
		newErrorResponse(ctx, response.ErrorResponseParameters{
			Internal:   err,
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
		})
		return
	}

	id, err := getUserId(ctx)
	if err != nil {
		newErrorResponse(ctx, err)
		return
	}

	list, err := h.services.DMList.Create(id, input)
	if err != nil {
		newErrorResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, list)
}

func (h *Handler) getLists(ctx *gin.Context) {
	id, err := getUserId(ctx)
	if err != nil {
		newErrorResponse(ctx, err)
		return
	}

	lists, err := h.services.DMList.GetUserLists(id)
	if err != nil {
		newErrorResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, lists)
}

func (h *Handler) deleteList(ctx *gin.Context) {
	userId, err := getUserId(ctx)
	if err != nil {
		newErrorResponse(ctx, err)
		return
	}

	listId, err := strconv.Atoi(ctx.Param("list_id"))
	if err != nil {
		newErrorResponse(ctx, response.ErrorResponseParameters{
			Internal:   err,
			Message:    "Invalid list ID parameter",
			StatusCode: http.StatusForbidden,
		})
		return
	}

	err = h.services.Delete(userId, listId)
	if err != nil {
		newErrorResponse(ctx, err)
		return
	}

	ctx.Status(http.StatusOK)
}
