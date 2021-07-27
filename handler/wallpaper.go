package handler

import (
	bing_wallpaper "github.com/TimothyYe/bing-wallpaper"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// RootHandler handles default API requests
func RootHandler(c *gin.Context) {
	resolution := c.DefaultQuery("resolution", "1920")
	format := c.DefaultQuery("format", "json")
	index := c.DefaultQuery("index", "0")
	mkt := c.DefaultQuery("mkt", "zh-CN")

	// check index
	uIndex, err := strconv.ParseUint(index, 10, 64)
	if err != nil {
		// input index is invalid
		c.String(http.StatusInternalServerError, "the image index is invalid")
		return
	}

	// check format
	if format != "json" && format != "image" {
		c.String(http.StatusInternalServerError, "format parameter is invalid, should be json or image")
		return
	}

	// fetch bing information
	response, err := bing_wallpaper.Get(uint(uIndex), mkt, resolution)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	}

	// redirect to image URL directly
	if format == "image" && response.URL != "" {
		c.Redirect(http.StatusTemporaryRedirect, response.URL)
		return
	}

	// render response as JSON
	c.JSON(http.StatusOK, response)
}
