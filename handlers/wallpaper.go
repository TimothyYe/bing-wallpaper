package handlers

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	bingURL = `https://www.bing.com`
	bingAPI = `https://www.bing.com/HPImageArchive.aspx?format=xml&idx=0&n=1&mkt=en-US`
)

func RootHandler(c *gin.Context) {
	// size := c.DefaultQuery("size", "1920")
	// format := c.DefaultQuery("foramt", "json")

	resp, err := http.Get(bingURL)
	if err != nil {
		// handle error
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
	}

	fmt.Println(string(body))
	c.String(http.StatusOK, string(body))
}
