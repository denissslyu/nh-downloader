package lrr_api

import (
	"net/http"
	"strings"

	"github.com/spf13/cast"

	"nh-downloader/internal/model"
	"nh-downloader/internal/service/lrr_serv"
	"nh-downloader/utils/logs"

	"github.com/gin-gonic/gin"
)

func GetItem(c *gin.Context) {
	id := c.Query("id")
	item, err := lrr_serv.GetItem(id)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
	c.JSON(http.StatusOK, item)
}

func Search(c *gin.Context) {
	query := c.Query("query")
	pageStr := c.Query("page")
	sort := c.Query("sort")
	order := c.Query("order")
	page := 1
	var err error
	if pageStr != "" {
		page, err = cast.ToIntE(pageStr)
		if err != nil {
			logs.Error("[lrr_api.Search] conv page to int failed:", err)
			page = 1
		}
	}
	option := &model.SimpleSearchOption{
		Filters: strings.Split(query, ","),
		Page:    page,
		Sort:    sort,
		Order:   order,
	}
	resp, err := lrr_serv.SimpleSearchItems(option)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
	c.JSON(http.StatusOK, resp)
}
