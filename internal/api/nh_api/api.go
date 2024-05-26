package nh_api

import (
	"net/http"
	"nh-downloader/internal/model"
	"nh-downloader/utils/logs"
	"strings"

	"github.com/spf13/cast"

	"nh-downloader/internal/service/nh_serv"

	"github.com/gin-gonic/gin"
)

func GetItem(c *gin.Context) {
	id := c.Query("id")

	item, err := nh_serv.GetItem(id)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
	c.JSON(http.StatusOK, item)
}

func Search(c *gin.Context) {
	query := c.Query("query")
	pageStr := c.Query("page")
	page := 1
	var err error
	if pageStr != "" {
		page, err = cast.ToIntE(pageStr)
		if err != nil {
			logs.Error("[nh_api.Search] conv page to int failed:", err)
			page = 1
		}
	}
	option := &model.SimpleSearchOption{
		Filters: strings.Split(query, ","),
		Page:    page,
	}
	resp, err := nh_serv.SimpleSearchItems(option)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
	c.JSON(http.StatusOK, resp)
}

func Download(c *gin.Context) {
	id := c.Query("id")

	err := nh_serv.DownloadByItemId(id)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
	c.String(http.StatusOK, "ok")
}
