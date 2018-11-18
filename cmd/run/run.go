package run

import (
	"bing-wallpaper/handler"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "run",
	Short: "Run this API service",
	Long:  `Run this API service`,
	Run:   run,
}

func run(cmd *cobra.Command, args []string) {
	router := gin.Default()

	router.GET("/", handler.RootHandler)

	router.Run(":7000")
	fmt.Println("API is running...")
}
