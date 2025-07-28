package git

import (
	"fmt"
	"os/exec"
	"strings"
)

// GetUntrackedFiles returns list of untracked files
func GetUntrackedFiles() ([]string, error) {
	cmd := exec.Command("git", "ls-files", "--others", "--exclude-standard")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get untracked files: %w", err)
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