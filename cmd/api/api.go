package run

import (
	"fmt"

	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run this API service",
	Long:  `Run this API service`,
	Run:   run,
}

func run(cmd *cobra.Command, args []string) {
	fmt.Println("API is running...")
}
