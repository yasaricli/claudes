package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"claudes/internal/config"

	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete [profile]",
	Short: "Delete an existing profile",
	Long:  `Delete an existing Claude configuration profile.`,
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
				fmt.Println("No profiles found.")
				os.Exit(1)
			}

			fmt.Println("Available profiles:")
			for i, p := range profiles {
				fmt.Printf("  %d. %s\n", i+1, p)
			}

			reader := bufio.NewReader(os.Stdin)
			for {
				fmt.Print("Select profile number or enter name to delete: ")
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
			fmt.Printf("Error: Profile '%s' not found.\n", profile)
			os.Exit(1)
		}

		// Confirm deletion
		reader := bufio.NewReader(os.Stdin)
		for {
			fmt.Printf("Are you sure you want to delete profile '%s'? (y/N): ", profile)
			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(strings.ToLower(input))

			if input == "y" || input == "yes" {
				break
			}

			if input == "" || input == "n" || input == "no" {
				fmt.Println("Deletion cancelled.")
				os.Exit(0)
			}

			fmt.Println("Invalid input. Please enter 'y' or 'n'.")
		}

		// Delete profile
		if err := config.DeleteProfile(profile); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		fmt.Printf("Profile '%s' deleted successfully!\n", profile)
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
