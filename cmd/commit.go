package cmd

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/atotto/clipboard"
	"github.com/bitcs/commet/internal/config"
	"github.com/bitcs/commet/internal/git"
	"github.com/bitcs/commet/internal/llm"
	"github.com/bitcs/commet/internal/ui"
	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	interactiveMode bool
	directCommit    bool
	useAI           bool
)

var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Generate AI-powered commit messages",
	Long: `Generate commit messages using AI based on your staged changes.
This command will analyze your git diff and suggest appropriate commit messages.

Examples:
  commet commit                    # Commit all staged/unstaged changes
  commet commit -i                 # Interactive file selection`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load()
		if err != nil {
			fmt.Printf("Error loading config: %v\n", err)
			return
		}

		if cfg.AI.APIKey == "" {
			fmt.Println("Error: No API key configured. Please run 'commet config set' first.")
			return
		}

		var gitDiff string

		// Check if interactive mode should be used (flag or config)
		shouldUseInteractive := interactiveMode || cfg.Git.Interactive

		if shouldUseInteractive {
			gitDiff, err = handleInteractiveFileSelection(cfg)
		} else {
			s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
			s.Suffix = " Analyzing git changes..."
			s.Start()
			gitDiff, err = git.GetDiff(cfg)
			s.Stop()
		}

		if err != nil {
			if err.Error() == "cancelled" {
				// User cancelled file selection - exit silently
				return
			}
			color.Red("Error getting git diff: %v\n", err)
			return
		}

		if gitDiff == "" {
			color.Yellow("No changes detected. Make sure you have staged changes or unstaged changes to commit.")
			return
		}

		color.Green("Found changes to commit")

		var commitMsg string

		// Check if AI should be used
		shouldUseAI := useAI || cfg.Git.UseAI

		if shouldUseAI {
			service, err := llm.NewService(cfg)
			if err != nil {
				fmt.Printf("Error creating LLM service: %v\n", err)
				return
			}

			s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
			s.Suffix = fmt.Sprintf(" Generating commit message using %s...", cfg.AI.Provider)
			s.Start()

			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			commitMsg, err = service.GenerateCommitMessage(ctx, gitDiff)
			s.Stop()

			if err != nil {
				color.Red("Error generating commit message: %v\n", err)
				return
			}
		} else {
			commitMsg, err = getManualCommitMessage()
			if err != nil {
				color.Red("Error getting commit message: %v\n", err)
				return
			}
		}

		if shouldUseAI {
			fmt.Printf("\n%s\n\n", commitMsg)
		}

		clipboard.WriteAll(commitMsg)

		shouldCommit := directCommit || cfg.Git.DirectCommit || !shouldUseAI

		if shouldCommit {
			if err := git.CreateCommit(commitMsg); err != nil {
				color.Red("Error creating commit: %v\n", err)
				return
			}
			color.Green("Commit created successfully!")
		} else {
			return
		}

		if cfg.Git.ConfirmPush {
			if askForConfirmation("\nDo you want to push the changes?") {
				if err := git.PushChanges(); err != nil {
					color.Red("Error pushing changes: %v\n", err)
					return
				}
				color.Green("Changes pushed successfully!")
			}
		}

	},
}

func askForConfirmation(question string) bool {
	color.Cyan("%s (y/N): ", question)
	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		return false
	}

	response = strings.ToLower(strings.TrimSpace(response))
	return response == "y" || response == "yes"
}

func handleInteractiveFileSelection(cfg *config.Config) (string, error) {
	unstagedFiles, err := git.GetUnstagedFiles()
	if err != nil {
		return "", fmt.Errorf("failed to get unstaged files: %w", err)
	}

	stagedFiles, err := git.GetStagedFiles()
	if err != nil {
		return "", fmt.Errorf("failed to get staged files: %w", err)
	}

	untrackedFiles, err := git.GetUntrackedFiles()
	if err != nil {
		return "", fmt.Errorf("failed to get untracked files: %w", err)
	}

	allFiles := make(map[string]bool)
	for _, file := range unstagedFiles {
		allFiles[file] = true
	}
	for _, file := range stagedFiles {
		allFiles[file] = true
	}
	for _, file := range untrackedFiles {
		allFiles[file] = true
	}

	availableFiles := make([]string, 0, len(allFiles))
	for file := range allFiles {
		availableFiles = append(availableFiles, file)
	}

	if len(availableFiles) == 0 {
		return "", fmt.Errorf("no files with changes found")
	}

	selectedFiles, filesToUnstage, err := ui.RunEnhancedFileSelector(availableFiles, stagedFiles, unstagedFiles, untrackedFiles)
	if err != nil {
		if err.Error() == "user cancelled file selection" {
			// User pressed 'q' to quit - exit gracefully without showing error
			return "", fmt.Errorf("cancelled")
		}
		return "", fmt.Errorf("file selection failed: %w", err)
	}

	if len(selectedFiles) == 0 {
		return "", fmt.Errorf("no files selected")
	}

	// Unstage files that were deselected
	if len(filesToUnstage) > 0 {
		color.Cyan("Unstaging deselected files: %v", filesToUnstage)
		if err := git.UnstageFiles(filesToUnstage); err != nil {
			return "", fmt.Errorf("failed to unstage files: %w", err)
		}
	}

	var filesToStage []string
	for _, file := range selectedFiles {
		needsStaging := false
		for _, unstaged := range unstagedFiles {
			if file == unstaged {
				needsStaging = true
				break
			}
		}
		if !needsStaging {
			for _, untracked := range untrackedFiles {
				if file == untracked {
					needsStaging = true
					break
				}
			}
		}
		if needsStaging {
			filesToStage = append(filesToStage, file)
		}
	}

	if len(filesToStage) > 0 {
		color.Cyan("Staging selected files: %v", filesToStage)
		if err := git.StageFiles(filesToStage); err != nil {
			return "", fmt.Errorf("failed to stage files: %w", err)
		}
	}

	return git.GetDiffForFiles(selectedFiles, true)
}

func getManualCommitMessage() (string, error) {
	color.Cyan("Enter commit message: ")

	reader := bufio.NewReader(os.Stdin)
	message, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("failed to read commit message: %w", err)
	}

	message = strings.TrimSpace(message)
	if message == "" {
		color.Red("Error: Commit message cannot be empty")
		return "", fmt.Errorf("commit message cannot be empty")
	}

	return message, nil
}

func init() {
	commitCmd.Flags().BoolVarP(&interactiveMode, "interactive", "i", false, "Interactive file selection mode")
	commitCmd.Flags().BoolVarP(&directCommit, "yes", "y", false, "Commit directly without confirmation")
	commitCmd.Flags().BoolVarP(&useAI, "ai", "a", false, "Use AI to generate commit message")
	rootCmd.AddCommand(commitCmd)
}
