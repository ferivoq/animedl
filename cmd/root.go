package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	headers map[string]map[string]string
	rootCmd = &cobra.Command{
		Use:   "animedrive-dl",
		Short: "🎬 Download an anime link",
		Long:  `🎬 Download an anime by link and a whole series.`,
	}
)

func Execute(h map[string]map[string]string) {
	headers = h
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
