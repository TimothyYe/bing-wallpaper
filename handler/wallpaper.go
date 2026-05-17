package handler

import (
	"math/rand/v2"
	"net/http"
	"strconv"
	"strings"

	bing_wallpaper "github.com/TimothyYe/bing-wallpaper"

	"github.com/gin-gonic/gin"
)

const (
	random = "random"
)

var (
	marketMap = map[int]string{
		0:  "zh-CN",
		1:  "en-US",
		2:  "ja-JP",
		3:  "en-AU",
		4:  "en-GB",
		5:  "de-DE",
		6:  "en-NZ",
		7:  "en-CA",
		8:  "en-IN",
		9:  "fr-FR",
		10: "fr-CA",
		11: "it-IT",
		12: "es-ES",
		13: "pt-BR",
		14: "en-ROW",
	}
)

func getRandomIndex() int {
	return rand.IntN(8)
}

func getRandomMarket() string {
	return marketMap[rand.IntN(len(marketMap))]
}

// RootHandler handles default API requests
func RootHandler(c *gin.Context) {
	resolution := c.DefaultQuery("resolution", "1920")
	format := c.DefaultQuery("format", "json")
	imageFormat := c.DefaultQuery("image_format", "jpg")
	index := c.DefaultQuery("index", "0")
	mkt := c.DefaultQuery("mkt", "zh-CN")

	// handle the random index
	if index == random {
		index = strconv.Itoa(getRandomIndex())
	}

	// handle the random market parameter
	if mkt == random {
		mkt = getRandomMarket()
	}

	// check index
	uIndex, err := strconv.ParseUint(index, 10, 64)
	if err != nil {
		c.String(http.StatusBadRequest, "the image index is invalid")
		return
	}

	// check format
	if format != "json" && format != "image" {
		c.String(http.StatusBadRequest, "format parameter is invalid, should be json or image")
		return
	}

	if imageFormat != "jpg" && imageFormat != "webp" {
		c.String(http.StatusBadRequest, "image_format parameter is invalid, should be jpg or webp")
		return
	}

	// fetch bing information
	response, err := bing_wallpaper.Get(c.Request.Context(), uint(uIndex), mkt, resolution)
	if err != nil {
		c.String(http.StatusBadGateway, err.Error())
		return
	}

	// check the image format
	if imageFormat == "webp" && strings.HasSuffix(response.URL, ".jpg") {
		response.URL = strings.TrimSuffix(response.URL, ".jpg") + ".webp"
	}

	// redirect to image URL directly
	if format == "image" && response.URL != "" {
		c.Redirect(http.StatusTemporaryRedirect, response.URL)
		return
	}

	// render response as JSON
	c.JSON(http.StatusOK, response)
}
