package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"claudes/internal/config"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add [profile]",
	Short: "Add a new profile interactively",
	Long:  `Add a new Claude configuration profile by entering credentials interactively.`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var profile string
		if len(args) > 0 {
			profile = args[0]
		}

		// Check if profile already exists
		if profile != "" && config.ProfileExists(profile) {
			fmt.Printf("Error: Profile '%s' already exists. Use 'claudes update %s' instead.\n", profile, profile)
			os.Exit(1)
		}

		// Interactive input
		reader := bufio.NewReader(os.Stdin)

		// Get profile name if not provided
		if profile == "" {
			for {
				fmt.Print("Enter profile name: ")
				input, _ := reader.ReadString('\n')
				profile = strings.TrimSpace(input)

				if profile == "" {
					fmt.Println("Profile name cannot be empty.")
					continue
				}

				if config.ProfileExists(profile) {
					fmt.Printf("Error: Profile '%s' already exists. Use 'claudes update %s' instead.\n", profile, profile)
					os.Exit(1)
				}

				break
			}
		}

		// Get auth token
		var authToken string
		for {
			fmt.Print("Enter ANTHROPIC_AUTH_TOKEN: ")
			input, _ := reader.ReadString('\n')
			authToken = strings.TrimSpace(input)

			if authToken == "" {
				fmt.Println("Auth token cannot be empty.")
				continue
			}
			break
		}

		// Get base URL
		var baseURL string
		for {
			fmt.Print("Enter ANTHROPIC_BASE_URL (e.g., https://api.anthropic.com/v1): ")
			input, _ := reader.ReadString('\n')
			baseURL = strings.TrimSpace(input)

			if baseURL == "" {
				fmt.Println("Base URL cannot be empty.")
				continue
			}
			break
		}

		// Get model
		var model string
		for {
			fmt.Print("Enter ANTHROPIC_MODEL (e.g., claude-3-opus-20240229): ")
			input, _ := reader.ReadString('\n')
			model = strings.TrimSpace(input)

			if model == "" {
				fmt.Println("Model cannot be empty.")
				continue
			}
			break
		}

		// Save profile
		profileConfig := config.ProfileConfig{
			AuthToken: authToken,
			BaseURL:   baseURL,
			Model:     model,
		}

		if err := config.SaveProfile(profile, profileConfig); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		fmt.Printf("Profile '%s' created successfully!\n", profile)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
