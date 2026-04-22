package cmd

import (
	"fmt"
	"os"

	"claudes/internal/config"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all available profiles",
	Long:  `List all .env files in the configuration directory.`,
	Run: func(cmd *cobra.Command, args []string) {
		profiles, err := config.ListProfiles()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		if len(profiles) == 0 {
			fmt.Println("No profiles found. Create a .env file in:", config.GetConfigDir())
			return
		}

		fmt.Println("Available profiles:")
		for _, p := range profiles {
			fmt.Println("  -", p)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
