package run

import (
	"fmt"
	"log"

	"github.com/TimothyYe/bing-wallpaper/handler"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "run",
	Short: "Run this API service",
	Long:  `Run this API service`,
	Run:   run,
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST,HEAD,PATCH, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func run(cmd *cobra.Command, args []string) {
	router := gin.Default()

	router.Use(CORSMiddleware())
	router.GET("/", handler.RootHandler)

	if err := router.Run(":9000"); err != nil {
		log.Fatal(err)
	}

	fmt.Println("API is running...")
}
