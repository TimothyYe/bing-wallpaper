package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RootHandler(c *gin.Context) {
	quality := c.DefaultQuery("quality", "high")
	c.String(http.StatusOK, fmt.Sprintf("Quality: %s", quality))
}
