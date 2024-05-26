package routes

import (
	"net/http"
	"nh-downloader/internal/api"
	"nh-downloader/internal/api/lrr_api"
	"nh-downloader/internal/api/nh_api"

	"github.com/gin-gonic/gin"
)

var Router *gin.Engine

// Init initialize router
func Init() {
	Router = gin.Default()
	Router.GET("/", HomeHandler)
	Router.GET("/dump", api.Dump)

	nh := Router.Group("/nh")
	{
		nh.GET("/search", nh_api.Search)
		nh.GET("/item", nh_api.GetItem)
		nh.GET("/dl", nh_api.Download)
	}
	lrr := Router.Group("/lrr")
	{
		lrr.GET("/search", lrr_api.Search)
		lrr.GET("/item", lrr_api.GetItem)
	}
}

func HomeHandler(c *gin.Context) {
	c.String(http.StatusOK, "Welcome to the home page")
}
