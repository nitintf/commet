package cmd

import (
	"fmt"

	"github.com/bitcs/commet/internal/ui"
	"github.com/spf13/cobra"
)

var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Generate AI-powered commit messages",
	Long: `Generate commit messages using AI based on your staged changes.
This command will analyze your git diff and suggest appropriate commit messages.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Starting commit message generator...")
		if err := ui.RunCommitUI(); err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(commitCmd)
}