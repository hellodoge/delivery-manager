package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/hellodoge/delivery-manager/dm"
	"github.com/hellodoge/delivery-manager/pkg/response"
	"net/http"
)

// @Summary Sign Up
// @Tags auth
// @Description create account
// @ID create-account
// @Accept json
// @Produce json
// @Param account-info body dm.User true "account info"
// @Success 200 {integer} integer 1
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} ServiceErrorResponse
// @Router /auth/sign-up [post]
func (h *Handler) signUp(ctx *gin.Context) {
	var input dm.User

	if err := ctx.BindJSON(&input); err != nil {
		newErrorResponse(ctx, response.ErrorResponseParameters{
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
		newErrorResponse(ctx, response.ErrorResponseParameters{
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

type refreshInput struct {
	RefreshToken string `json:"refresh"`
}

func (h *Handler) refresh(ctx *gin.Context) {
	var input refreshInput
	if err := ctx.BindJSON(&input); err != nil {
		newErrorResponse(ctx, response.ErrorResponseParameters{
			Internal:   err,
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
		})
		return
	}

	token, err := h.services.GenerateToken(input.RefreshToken)
	if err != nil {
		newErrorResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
	})
}
