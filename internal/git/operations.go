package git

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/bitcs/commet/internal/config"
)

func GetDiff(cfg *config.Config) (string, error) {
	// First try staged changes
	cmd := exec.Command("git", "diff", "--cached")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get staged changes: %w", err)
	}

	// If no staged changes and auto-stage is enabled, stage all changes
	if len(output) == 0 {
		if cfg.Git.AutoStage {
			if err := StageAllChanges(); err != nil {
				return "", fmt.Errorf("failed to auto-stage changes: %w", err)
			}
			// Get staged changes after auto-staging
			cmd = exec.Command("git", "diff", "--cached")
			output, err = cmd.Output()
			if err != nil {
				return "", fmt.Errorf("failed to get staged changes after auto-staging: %w", err)
			}
		} else {
			// Get working directory changes
			cmd = exec.Command("git", "diff")
			output, err = cmd.Output()
			if err != nil {
				return "", fmt.Errorf("failed to get working directory changes: %w", err)
			}
		}
	}

	return string(output), nil
}

func StageAllChanges() error {
	cmd := exec.Command("git", "add", ".")
	return cmd.Run()
}

func CreateCommit(message string) error {
	cmd := exec.Command("git", "commit", "-m", message)
	return cmd.Run()
}

func PushChanges() error {
	cmd := exec.Command("git", "push")
	return cmd.Run()
}

func GetStatus() (string, error) {
	cmd := exec.Command("git", "status", "--porcelain")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get git status: %w", err)
	}
	return string(output), nil
}

func IsRepository() bool {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	err := cmd.Run()
	return err == nil
}

func GetCurrentBranch() (string, error) {
	cmd := exec.Command("git", "branch", "--show-current")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get current branch: %w", err)
	}
	return string(output), nil
}

func StageFiles(files []string) error {
	if len(files) == 0 {
		return nil
	}

	args := append([]string{"add"}, files...)
	cmd := exec.Command("git", args...)
	return cmd.Run()
}

func UnstageFiles(files []string) error {
	if len(files) == 0 {
		return nil
	}

	args := append([]string{"reset", "HEAD", "--"}, files...)
	cmd := exec.Command("git", args...)
	return cmd.Run()
}

func GetDiffForFiles(files []string, staged bool) (string, error) {
	args := []string{"diff"}

	if staged {
		args = append(args, "--cached")
	}

	if len(files) > 0 {
		args = append(args, "--")
		args = append(args, files...)
	}

	cmd := exec.Command("git", args...)
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get diff for files: %w", err)
	}

	return string(output), nil
}

func GetUnstagedFiles() ([]string, error) {
	cmd := exec.Command("git", "diff", "--name-only")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get unstaged files: %w", err)
	}

	files := []string{}
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	for _, line := range lines {
		if line != "" {
			files = append(files, line)
		}
	}

	return files, nil
}

func GetStagedFiles() ([]string, error) {
	cmd := exec.Command("git", "diff", "--cached", "--name-only")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get staged files: %w", err)
	}

	files := []string{}
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	for _, line := range lines {
		if line != "" {
			files = append(files, line)
		}
	}

	return files, nil
}

// GetFileStatus returns the git status for a single file (M, A, D, etc.)
func GetFileStatus(filename string) (string, error) {
	cmd := exec.Command("git", "status", "--porcelain", filename)
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get file status: %w", err)
	}
	
	status := strings.TrimSpace(string(output))
	if len(status) >= 2 {
		return status[:2], nil
	}
	return "??", nil // untracked
}

// GetSingleFileDiff returns the diff for a single file
func GetSingleFileDiff(filename string) (string, error) {
	// Try staged diff first
	cmd := exec.Command("git", "diff", "--cached", filename)
	output, err := cmd.Output()
	if err == nil && len(output) > 0 {
		return string(output), nil
	}
	
	// Try unstaged diff
	cmd = exec.Command("git", "diff", filename)
	output, err = cmd.Output()
	if err == nil && len(output) > 0 {
		return string(output), nil
	}
	
	// For untracked files, show the entire file content as additions
	cmd = exec.Command("git", "ls-files", "--others", "--exclude-standard", filename)
	output, err = cmd.Output()
	if err == nil && len(output) > 0 {
		// This is an untracked file, show it as all additions
		cmd = exec.Command("cat", filename)
		content, err := cmd.Output()
		if err != nil {
			return "", fmt.Errorf("failed to read untracked file: %w", err)
		}
		
		// Format as diff-like output
		lines := strings.Split(string(content), "\n")
		var diff strings.Builder
		diff.WriteString(fmt.Sprintf("--- /dev/null\n+++ b/%s\n", filename))
		diff.WriteString("@@ -0,0 +1," + fmt.Sprintf("%d", len(lines)) + " @@\n")
		for _, line := range lines {
			diff.WriteString("+" + line + "\n")
		}
		return diff.String(), nil
	}
	
	return "", fmt.Errorf("no diff found for file: %s", filename)
}
