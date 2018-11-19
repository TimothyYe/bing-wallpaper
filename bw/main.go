package main

import (
	"fmt"
	"os"

	"github.com/TimothyYe/bing-wallpaper/cmd/run"

	"github.com/spf13/cobra"
)

var (
	argVerbose bool
	rootCmd    *cobra.Command
)

func init() {
	rootCmd = &cobra.Command{
		Use:   "bw",
		Short: "Bing wallpaper API",
		Long:  "Top level command for Bing wallpaper API service",
	}

	rootCmd.PersistentFlags().BoolVarP(&argVerbose, "verbose", "v", false, "verbose output")

	rootCmd.AddCommand(
		run.Cmd,
	)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
