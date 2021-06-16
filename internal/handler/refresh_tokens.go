package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/hellodoge/delivery-manager/pkg/response"
	"net/http"
)

type getRefreshTokensInput struct {
	Date string `json:"issued-after"`
}

func (h *Handler) getRefreshTokens(ctx *gin.Context) {
	var input getRefreshTokensInput
	if err := ctx.BindJSON(&input); err != nil {
		newErrorResponse(ctx, response.ErrorResponseParameters{
			Internal:   err,
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
		})
		return
	}
	userId, err := getUserId(ctx)
	if err != nil {
		newErrorResponse(ctx, err)
		return
	}
	tokens, err := h.services.GetUserRefreshTokens(userId, input.Date)
	if err != nil {
		newErrorResponse(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, tokens)
}
