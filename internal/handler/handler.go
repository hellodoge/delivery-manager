package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/hellodoge/delivery-manager/internal/service"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	var router = gin.New()

	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.signUp)
		auth.POST("/sign-in", h.signIn)
		auth.POST("/refresh", h.refresh)
	}

	api := router.Group("/api", h.userIdentity)
	{
		deliveries := api.Group("/delivery")
		{
			deliveries.GET("/", h.getDeliveries)
			deliveries.POST("/create", h.createDelivery)
			deliveries.GET("/:delivery_id", h.deliveryInfo)
			deliveries.PUT("/:delivery_id", h.updateDelivery)
			deliveries.DELETE("/:delivery_id", h.removeDelivery)
		}

		lists := api.Group("/lists")
		{
			lists.GET("/", h.getLists)
			lists.POST("/create", h.createList)
			lists.DELETE("/:list_id", h.deleteList)

			products := lists.Group("/:list_id/products")
			{
				products.GET("/", h.getProducts)
				products.POST("/", h.addProducts)
				products.PUT("/", h.updateProducts)
				products.DELETE("/", h.removeProducts)
			}
		}

		products := api.Group("/products")
		{
			products.POST("/", h.createProducts)
			products.POST("/search", h.searchForProducts)
		}
	}

	return router
}
