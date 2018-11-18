package handler

import (
	"bing-wallpaper/model"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/beevik/etree"

	"github.com/gin-gonic/gin"
)

const (
	bingURL = `https://www.bing.com`
	bingAPI = `https://www.bing.com/HPImageArchive.aspx?format=xml&idx=%s&n=1&mkt=en-US`
)

func RootHandler(c *gin.Context) {
	// size := c.DefaultQuery("size", "1920")
	// format := c.DefaultQuery("foramt", "json")
	index := c.DefaultQuery("index", "0")

	_, err := strconv.Atoi(index)
	if err != nil {
		// input index is invalid
		c.String(http.StatusInternalServerError, "the image index is invalid")
	}

	resp, err := http.Get(fmt.Sprintf(bingAPI, index))
	if err != nil {
		c.String(http.StatusInternalServerError, "failed to request bing.com")
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.String(http.StatusInternalServerError, "failed to parse request body from bing.com")
	}

	response, err := parseResponse(body)
	if err != nil {
		c.String(http.StatusInternalServerError, "failed to parse request body from bing.com")
	}

	c.JSON(http.StatusOK, response)
}

func parseResponse(xmlInput []byte) (*model.Response, error) {
	doc := etree.NewDocument()
	if err := doc.ReadFromBytes(xmlInput); err != nil {
		return nil, err
	}

	// get image element
	imgElem := doc.SelectElement("images").SelectElement("image")

	resp := &model.Response{}
	resp.StartDate = imgElem.SelectElement("startdate").Text()
	resp.EndDate = imgElem.SelectElement("enddate").Text()
	resp.URL = fmt.Sprintf("%s_%s")
	return resp, nil
}
