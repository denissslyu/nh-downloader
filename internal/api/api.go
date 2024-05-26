package api

import (
	"net/http"
	"nh-downloader/internal/service"
	"strings"

	"github.com/gin-gonic/gin"
)

func Dump(c *gin.Context) {
	query := c.Query("query")

	err := service.Dump(strings.Split(query, ","))
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
	c.JSON(http.StatusOK, "ok")
}
