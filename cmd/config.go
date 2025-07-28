package cmd

import (
	"fmt"

	"github.com/bitcs/commet/internal/config"
	"github.com/bitcs/commet/internal/ui"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration settings",
	Long:  `View and manage configuration settings for AI providers, API keys, and other options.`,
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration",
	Long:  `Display the current configuration settings with masked API keys for security.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load()
		if err != nil {
			fmt.Printf("Error loading config: %v\n", err)
			return
		}

		fmt.Println("Current Configuration:")
		fmt.Printf("  AI Provider: %s\n", cfg.AI.Provider)
		fmt.Printf("  API Key: %s\n", cfg.MaskAPIKey())

		model := cfg.AI.Model
		if model == "" {
			model = cfg.GetDefaultModel()
		}
		fmt.Printf("  Model: %s\n", model)
		
		fmt.Println("\nGit Settings:")
		fmt.Printf("  Auto Stage: %t\n", cfg.Git.AutoStage)
		fmt.Printf("  Show Diff: %t\n", cfg.Git.ShowDiff)
		fmt.Printf("  Confirm Push: %t\n", cfg.Git.ConfirmPush)
	},
}

var configSetCmd = &cobra.Command{
	Use:   "set",
	Short: "Set configuration values",
	Long:  `Set configuration values using command line flags or interactive mode.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load()
		if err != nil {
			fmt.Printf("Error loading config: %v\n", err)
			return
		}

		provider, _ := cmd.Flags().GetString("provider")
		apiKey, _ := cmd.Flags().GetString("api-key")
		model, _ := cmd.Flags().GetString("model")

		if provider == "" && apiKey == "" && model == "" {
			if err := ui.RunConfigUI(cfg); err != nil {
				fmt.Printf("Error: TUI not available (%v)\n", err)
				return
			}
			return
		}

		if provider != "" {
			p, err := config.ParseProvider(provider)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}
			cfg.AI.Provider = p
		}

		if apiKey != "" {
			cfg.AI.APIKey = apiKey
		}

		if model != "" {
			cfg.AI.Model = model
		}

		if err := cfg.Save(); err != nil {
			fmt.Printf("Error saving config: %v\n", err)
			return
		}

		fmt.Println("Configuration updated successfully!")
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configShowCmd)
	configCmd.AddCommand(configSetCmd)

	configSetCmd.Flags().StringP("provider", "p", "", "AI provider (openai, claude, google)")
	configSetCmd.Flags().StringP("api-key", "k", "", "API key for the AI provider")
	configSetCmd.Flags().StringP("model", "m", "", "AI model to use")
}
