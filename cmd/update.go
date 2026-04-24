package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"claudes/internal/config"

	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update [profile]",
	Short: "Update an existing profile interactively",
	Long:  `Update an existing Claude configuration profile by entering new credentials interactively.`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var profile string
		if len(args) > 0 {
			profile = args[0]
		}

		// If no profile provided, list available profiles
		if profile == "" {
			profiles, err := config.ListProfiles()
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			if len(profiles) == 0 {
				fmt.Println("No profiles found. Use 'claudes add' to create one.")
				os.Exit(1)
			}

			fmt.Println("Available profiles:")
			for i, p := range profiles {
				fmt.Printf("  %d. %s\n", i+1, p)
			}

			reader := bufio.NewReader(os.Stdin)
			for {
				fmt.Print("Select profile number or enter name: ")
				input, _ := reader.ReadString('\n')
				input = strings.TrimSpace(input)

				// Check if it's a number
				var selected int
				fmt.Sscanf(input, "%d", &selected)
				if selected > 0 && selected <= len(profiles) {
					profile = profiles[selected-1]
					break
				}

				// Check if it's a valid profile name
				found := false
				for _, p := range profiles {
					if p == input {
						profile = input
						found = true
						break
					}
				}

				if found {
					break
				}

				fmt.Println("Invalid selection. Please try again.")
			}
		}

		// Check if profile exists
		if !config.ProfileExists(profile) {
			fmt.Printf("Error: Profile '%s' not found. Use 'claudes add %s' to create it.\n", profile, profile)
			os.Exit(1)
		}

		// Load existing profile to get current values BEFORE asking input
		existingVars, err := config.LoadProfile(profile)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		// Parse existing values
		existingConfig := config.ProfileConfig{}
		for _, v := range existingVars {
			parts := strings.SplitN(v, "=", 2)
			if len(parts) == 2 {
				key := parts[0]
				value := parts[1]
				switch key {
				case "ANTHROPIC_AUTH_TOKEN":
					existingConfig.AuthToken = value
				case "ANTHROPIC_BASE_URL":
					existingConfig.BaseURL = value
				case "ANTHROPIC_MODEL":
					existingConfig.Model = value
				}
			}
		}

		// Mask auth token for display
		maskedToken := ""
		if len(existingConfig.AuthToken) > 8 {
			maskedToken = existingConfig.AuthToken[:4] + "********" + existingConfig.AuthToken[len(existingConfig.AuthToken)-4:]
		} else if existingConfig.AuthToken != "" {
			maskedToken = "********"
		}

		// Interactive input
		reader := bufio.NewReader(os.Stdin)

		fmt.Printf("\nUpdating profile '%s' (leave empty to keep current value)\n\n", profile)

		// Get auth token
		if maskedToken != "" {
			fmt.Printf("Enter ANTHROPIC_AUTH_TOKEN (current: %s): ", maskedToken)
		} else {
			fmt.Print("Enter ANTHROPIC_AUTH_TOKEN: ")
		}
		authTokenInput, _ := reader.ReadString('\n')
		authToken := strings.TrimSpace(authTokenInput)

		// Get base URL
		if existingConfig.BaseURL != "" {
			fmt.Printf("Enter ANTHROPIC_BASE_URL (current: %s): ", existingConfig.BaseURL)
		} else {
			fmt.Print("Enter ANTHROPIC_BASE_URL: ")
		}
		baseURLInput, _ := reader.ReadString('\n')
		baseURL := strings.TrimSpace(baseURLInput)

		// Get model
		if existingConfig.Model != "" {
			fmt.Printf("Enter ANTHROPIC_MODEL (current: %s): ", existingConfig.Model)
		} else {
			fmt.Print("Enter ANTHROPIC_MODEL: ")
		}
		modelInput, _ := reader.ReadString('\n')
		model := strings.TrimSpace(modelInput)

		// Use new values or keep existing
		updatedConfig := config.ProfileConfig{
			AuthToken: existingConfig.AuthToken,
			BaseURL:   existingConfig.BaseURL,
			Model:     existingConfig.Model,
		}

		if authToken != "" {
			updatedConfig.AuthToken = authToken
		}
		if baseURL != "" {
			updatedConfig.BaseURL = baseURL
		}
		if model != "" {
			updatedConfig.Model = model
		}

		// Save updated profile
		if err := config.SaveProfile(profile, updatedConfig); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		fmt.Printf("Profile '%s' updated successfully!\n", profile)
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
