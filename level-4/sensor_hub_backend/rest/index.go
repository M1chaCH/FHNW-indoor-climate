package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterIndexRoutes(router gin.IRouter) {
	router.GET("", getIndex)
}

func getIndex(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}
