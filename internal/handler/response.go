package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/hellodoge/delivery-manager/pkg/response"
	"github.com/sirupsen/logrus"
	"net/http"
)

func newErrorResponse(c *gin.Context, err error) {
	switch e := err.(type) {
	case response.ErrorResponse:
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
		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "Unknown Internal Server Error",
			"message": "Please, contact API service team",
		})
	}
}
