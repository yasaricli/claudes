package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"claudes/internal/config"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "claudes",
	Short: "Claude CLI configuration manager",
	Long:  `A CLI tool to manage multiple Claude configurations using environment files.`,
	Args:  cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			os.Exit(1)
		}

		profile := args[0]
		if err := runProfile(profile, args[1:]); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	},
}

func init() {
	// Initialize config directory if it doesn't exist
	config.InitConfigDir()
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func runProfile(profile string, remainingArgs []string) error {
	// Load environment variables from the profile
	envVars, err := config.LoadProfile(profile)
	if err != nil {
		return err
	}

	// Find the Claude CLI executable
	claudePath, err := exec.LookPath("claude")
	if err != nil {
		return fmt.Errorf("error: 'claude' command not found in PATH. Please install Claude CLI first.")
	}

	// Prepare the command with loaded environment variables
	c := exec.Command(claudePath, remainingArgs...)
	c.Env = append(os.Environ(), envVars...)
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr

	// Run the Claude CLI
	if err := c.Run(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			os.Exit(exitError.ExitCode())
		}
		return fmt.Errorf("error running claude: %w", err)
	}

	return nil
}
