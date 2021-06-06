package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/hellodoge/delivery-manager/dm"
	"github.com/hellodoge/delivery-manager/pkg/response"
	"net/http"
)

func (h *Handler) signUp(ctx *gin.Context) {
	var input dm.User

	if err := ctx.BindJSON(&input); err != nil {
		newErrorResponse(ctx, response.ErrorResponse{
			Internal:   err,
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
		})
		return
	}

	id, err := h.services.Authorization.CreateUser(input)
	if err != nil {
		newErrorResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

type signInInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *Handler) signIn(ctx *gin.Context) {
	var input signInInput

	if err := ctx.BindJSON(&input); err != nil {
		newErrorResponse(ctx, response.ErrorResponse{
			Internal:   err,
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
		})
		return
	}

	refreshToken, err := h.services.GenerateRefreshToken(input.Username, input.Password)
	if err != nil {
		newErrorResponse(ctx, err)
		return
	}
	token, err := h.services.GenerateToken(refreshToken)
	if err != nil {
		newErrorResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"refresh": refreshToken,
		"token":   token,
	})
}
