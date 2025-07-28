package ui

import (
	"fmt"

	"github.com/atotto/clipboard"
)

func DisplayCommitMessage(commitMsg string, isAIGenerated bool) {
	if isAIGenerated {
		fmt.Printf("\n%s\n\n", commitMsg)
	}

	clipboard.WriteAll(commitMsg)
}
