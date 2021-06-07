package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/hellodoge/delivery-manager/pkg/response"
	"github.com/sirupsen/logrus"
	"net/http"
)

type ServiceErrorResponse struct {
	ErrorType string `json:"error"`
	Message   string `json:"message"`
}

func newErrorResponse(c *gin.Context, err error) {
	switch e := err.(type) {
	case response.ErrorResponseParameters:
		if e.IsInternal {
			logrus.Error(e)
		} else {
			logrus.Info(e)
		}
		var status int = http.StatusInternalServerError
		if e.StatusCode != 0 {
			status = e.StatusCode
		}
		c.AbortWithStatusJSON(status, e.Response())
	default:
		logrus.Error(e)
		c.AbortWithStatusJSON(http.StatusInternalServerError, ServiceErrorResponse{
			ErrorType: "Unknown Internal Server Error",
			Message:   "Please, contact API service team",
		})
	}
}
