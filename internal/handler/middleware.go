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
		newErrorResponse(ctx, response.ErrorResponse{
			Message:    "Empty Authorization header",
			StatusCode: http.StatusUnauthorized,
		})
		return
	}
	header = strings.TrimPrefix(header, "Bearer ")

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