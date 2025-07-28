package prompts

import "fmt"

// CommitMessagePrompt generates a prompt for creating commit messages based on git diff
func CommitMessagePrompt(gitDiff string) string {
	return fmt.Sprintf(`You are an expert software engineer with years of experience writing clear, professional commit messages that follow industry best practices.

Analyze the following git diff and generate a commit message that:

**Format Requirements:**
- Subject line: 50 characters or less, imperative mood (e.g., "Add", "Fix", "Update", "Remove")
- Use conventional commit format when appropriate (feat:, fix:, docs:, refactor:, etc.)
- If the change is complex, include a brief body (optional, max 72 chars per line)

**Content Guidelines:**
- Focus on WHAT changed and WHY, not HOW
- Be specific but concise 
- Use present tense, imperative mood ("Add feature" not "Added feature")
- Avoid generic messages like "update code" or "fix bug"
- For multiple related changes, focus on the primary purpose

**Context Analysis:**
- Look for new files, modified files, deletions
- Identify the type of change: feature, bugfix, refactor, docs, test, etc.
- Consider the scope: which components/modules are affected

Git diff:
%s

Return ONLY the commit message (subject line + optional body if needed). No explanations, comments, or additional text.`, gitDiff)
}