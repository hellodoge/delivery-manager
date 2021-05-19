package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) getDeliveries(ctx *gin.Context) {
	id, _ := ctx.Get(userIdContextKey)
	ctx.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

func (h *Handler) createDelivery(ctx *gin.Context) {

}

func (h *Handler) removeDelivery(ctx *gin.Context) {

}

func (h *Handler) deliveryInfo(ctx *gin.Context) {

}

func (h *Handler) updateDelivery(ctx *gin.Context) {

}