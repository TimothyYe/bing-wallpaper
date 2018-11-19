package handler

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/TimothyYe/bing-wallpaper/model"

	"github.com/beevik/etree"

	"github.com/gin-gonic/gin"
)

const (
	bingURL = `https://www.bing.com`
	bingAPI = `https://www.bing.com/HPImageArchive.aspx?format=xml&idx=%s&n=1&mkt=%s`
)

var (
	resolutionMap map[string]string
	markets       map[string]bool
)

func init() {
	resolutionMap = map[string]string{}
	resolutionMap["1366"] = "1366x768.jpg"
	resolutionMap["1920"] = "1920x1080.jpg"

	markets = map[string]bool{
		"en-US": true,
		"zh-CN": true,
		"ja-JP": true,
		"en-AU": true,
		"en-UK": true,
		"de-DE": true,
		"en-NZ": true,
		"en-CA": true,
	}
}

// RootHandler handles default API requests
func RootHandler(c *gin.Context) {
	resolution := c.DefaultQuery("resolution", "1920")
	format := c.DefaultQuery("format", "json")
	index := c.DefaultQuery("index", "0")
	mkt := c.DefaultQuery("mkt", "zh-CN")

	// check index
	_, err := strconv.Atoi(index)
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

	// check resolution
	if resolution != "1920" && resolution != "1366" {
		c.String(http.StatusInternalServerError, "resolution parameter is invalid, should be 1920 or 1366")
		return
	}

	// check mkt
	if !markets[mkt] {
		c.String(http.StatusInternalServerError, "mkt parameter is invalid")
		return
	}

	resp, err := http.Get(fmt.Sprintf(bingAPI, index, mkt))
	if err != nil {
		c.String(http.StatusInternalServerError, "failed to request bing.com")
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.String(http.StatusInternalServerError, "failed to parse request body from bing.com")
		return
	}

	response, err := parseResponse(body, resolution)
	if err != nil {
		c.String(http.StatusInternalServerError, "failed to parse request body from bing.com")
		return
	}

	if format == "image" {
		// redirect to image URL directly
		c.Redirect(http.StatusTemporaryRedirect, response.URL)
		return
	}

	// render response as JSON
	c.JSON(http.StatusOK, response)
}

func parseResponse(xmlInput []byte, resolution string) (*model.Response, error) {
	doc := etree.NewDocument()
	if err := doc.ReadFromBytes(xmlInput); err != nil {
		return nil, err
	}

	// get image element
	imgElem := doc.SelectElement("images").SelectElement("image")

	resp := &model.Response{}
	resp.StartDate = imgElem.SelectElement("startdate").Text()
	resp.EndDate = imgElem.SelectElement("enddate").Text()
	resp.URL = fmt.Sprintf("%s%s_%s", bingURL, imgElem.SelectElement("urlBase").Text(), resolutionMap[resolution])
	resp.Copyright = imgElem.SelectElement("copyright").Text()
	resp.CopyrightLink = imgElem.SelectElement("copyrightlink").Text()

	return resp, nil
}
