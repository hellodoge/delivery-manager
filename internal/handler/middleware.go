package handler

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hellodoge/delivery-manager/pkg/response"
	"net/http"
	"strings"
)

const (
	authorizationHeader = "Authorization"

	userIdContextKey = "userId"
)

func (h *Handler) userIdentity(ctx *gin.Context) {
	header := ctx.GetHeader(authorizationHeader)
	if header == "" {
		newErrorResponse(ctx, response.ErrorResponseParameters{
			Message:    "Empty Authorization header",
			StatusCode: http.StatusUnauthorized,
		})
		return
	}
	header = strings.TrimPrefix(header, "Bearer ")
	if header == "" {
		newErrorResponse(ctx, response.ErrorResponseParameters{
			Message:    "Invalid Bearer token",
			StatusCode: http.StatusUnauthorized,
		})
		return
	}
	if strings.ContainsRune(header, ' ') {
		newErrorResponse(ctx, response.ErrorResponseParameters{
			Message:    "Empty Bearer token",
			StatusCode: http.StatusUnauthorized,
		})
		return
	}

	userId, err := h.services.Authorization.ParseToken(header)
	if err != nil {
		newErrorResponse(ctx, err)
		return
	}

	ctx.Set(userIdContextKey, userId)
}

func getUserId(ctx *gin.Context) (int, error) {
	idContext, ok := ctx.Get(userIdContextKey)
	if !ok {
		return -1, errors.New("user id not found")
	}
	id, ok2 := idContext.(int)
	if !ok2 {
		return -1, fmt.Errorf("user id has other type, than int (%T)", idContext)
	}
	return id, nil
}
