package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Version information - set during build
var Version = "dev"
var Commit = "none"
var Date = "unknown"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print claudes version information",
	Long:  `Print the version, commit hash, and build date of claudes.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Claudes - Claude CLI Profile Manager")
		fmt.Println("")
		fmt.Println("Version:", Version)
		fmt.Println("Commit:", Commit)
		fmt.Println("Built:", Date)
		fmt.Println("")
		fmt.Println("For updates, run:")
		fmt.Println("  curl -fsSL https://raw.githubusercontent.com/yasaricli/claudes/refs/heads/develop/install.sh | bash")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
