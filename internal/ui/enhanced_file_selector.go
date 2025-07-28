package ui

import (
	"fmt"
	"strings"

	"github.com/bitcs/commet/internal/git"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type enhancedFileModel struct {
	files         []string
	cursor        int
	selected      map[int]bool
	currentDiff   string
	loadingDiff   bool
	windowWidth   int
	windowHeight  int
	diffScroll    int
	diffLines     []string
	fileStatus    map[string]string // Maps filename to status (staged, unstaged, untracked)
	initialStaged map[string]bool   // Track which files were initially staged
	quitted       bool              // Track if user quit the selection
}

var (
	listStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("240")).
		Padding(0, 1).
		Width(40)

	diffStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("240")).
		Padding(0, 1)

	selectedStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("15")).
		Background(lipgloss.Color("237")).
		Bold(true)

	cursorStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("11"))

	headerStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("14")).
		Bold(true).
		Align(lipgloss.Center)

	// Diff syntax colors - more muted and terminal-friendly
	addedStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("10"))

	removedStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("9"))

	hunkStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("6"))

	contextStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("7"))

	// File status indicators
	stagedIndicatorStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("10")).
		Bold(true)

	unstagedIndicatorStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("11")).
		Bold(true)

	untrackedIndicatorStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("12")).
		Bold(true)
)

func NewEnhancedFileSelector(files []string, stagedFiles, unstagedFiles, untrackedFiles []string) enhancedFileModel {
	fileStatus := make(map[string]string)
	initialStaged := make(map[string]bool)
	selected := make(map[int]bool)
	
	// Create sets for quick lookup
	stagedSet := make(map[string]bool)
	unstagedSet := make(map[string]bool)
	untrackedSet := make(map[string]bool)
	
	for _, file := range stagedFiles {
		stagedSet[file] = true
	}
	for _, file := range unstagedFiles {
		unstagedSet[file] = true
	}
	for _, file := range untrackedFiles {
		untrackedSet[file] = true
	}
	
	// Map file status and pre-select staged files
	// Handle each file in the combined list
	for i, file := range files {
		// Determine primary status with clear priority
		isStaged := stagedSet[file]
		isUnstaged := unstagedSet[file]
		isUntracked := untrackedSet[file]
		
		if isStaged {
			// File has staged changes - show as staged and pre-select
			fileStatus[file] = "staged"
			initialStaged[file] = true
			selected[i] = true
		} else if isUnstaged {
			// File has only unstaged changes - show as modified, don't pre-select
			fileStatus[file] = "unstaged"
			// Don't pre-select unstaged files
		} else if isUntracked {
			// File is untracked - show as new, don't pre-select  
			fileStatus[file] = "untracked"
			// Don't pre-select untracked files
		} else {
			// This shouldn't happen if our file lists are correct
			fileStatus[file] = "unknown"
		}
	}
	
	return enhancedFileModel{
		files:         files,
		selected:      selected,
		fileStatus:    fileStatus,
		initialStaged: initialStaged,
	}
}

func (m enhancedFileModel) Init() tea.Cmd {
	return m.loadDiffForCurrentFile()
}

func (m enhancedFileModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.windowWidth = msg.Width
		m.windowHeight = msg.Height
		// Update diff style width based on window size
		diffStyle = diffStyle.Width(msg.Width - 50) // Leave space for file list + margins

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitted = true
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
				return m, m.loadDiffForCurrentFile()
			}

		case "down", "j":
			if m.cursor < len(m.files)-1 {
				m.cursor++
				return m, m.loadDiffForCurrentFile()
			}

		case " ":
			// Toggle selection
			m.selected[m.cursor] = !m.selected[m.cursor]

		case "enter":
			// Confirm selection
			if len(m.selected) == 0 {
				// If nothing selected, select current file
				m.selected[m.cursor] = true
			}
			return m, tea.Quit

		case "left", "h":
			// Scroll diff up
			if m.diffScroll > 0 {
				m.diffScroll--
			}

		case "right", "l":
			// Scroll diff down
			maxScroll := len(m.diffLines) - (m.windowHeight - 8) // Account for borders and headers
			if maxScroll < 0 {
				maxScroll = 0
			}
			if m.diffScroll < maxScroll {
				m.diffScroll++
			}
		}

	case diffLoadedMsg:
		m.currentDiff = string(msg)
		m.diffLines = strings.Split(string(msg), "\n")
		m.loadingDiff = false
		m.diffScroll = 0 // Reset scroll when loading new diff
	}

	return m, nil
}

func (m enhancedFileModel) View() string {
	if m.windowWidth == 0 {
		return "Loading..."
	}

	// File list section
	fileListContent := headerStyle.Render("Files to commit") + "\n\n"
	
	for i, file := range m.files {
		cursor := "  "
		if m.cursor == i {
			cursor = "❯ "
		}

		checkbox := "☐ "
		if m.selected[i] {
			checkbox = "☑ "
		}

		// Add status indicator
		var statusIndicator string
		status := m.fileStatus[file]
		switch status {
		case "staged":
			statusIndicator = stagedIndicatorStyle.Render("●")
		case "unstaged":
			statusIndicator = unstagedIndicatorStyle.Render("◯")
		case "untracked":
			statusIndicator = untrackedIndicatorStyle.Render("✦")
		default:
			statusIndicator = " " // Padding for alignment
		}

		line := cursor + checkbox + statusIndicator + " " + file
		
		if m.cursor == i {
			line = cursorStyle.Render(line)
		}
		if m.selected[i] {
			line = selectedStyle.Render(line)
		}

		fileListContent += line + "\n"
	}

	fileListContent += "\n" + lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Render("Space: select, Enter: confirm\nh/l: scroll diff, q: quit\n● staged ◯ modified ✦ added")

	fileList := listStyle.Render(fileListContent)

	// Diff section
	diffContent := headerStyle.Render("File Changes")
	
	if m.loadingDiff {
		diffContent += "\n\nLoading diff..."
	} else if len(m.diffLines) > 0 {
		// Calculate available height for diff content
		maxDiffHeight := m.windowHeight - 8 // Account for borders, headers, and help text
		if maxDiffHeight < 5 {
			maxDiffHeight = 5
		}
		
		// Get visible lines based on scroll position
		startLine := m.diffScroll
		endLine := startLine + maxDiffHeight
		if endLine > len(m.diffLines) {
			endLine = len(m.diffLines)
		}
		
		visibleLines := m.diffLines[startLine:endLine]
		diffContent += "\n\n" + m.formatDiffLines(visibleLines)
		
		// Add scroll indicator if needed
		if len(m.diffLines) > maxDiffHeight {
			scrollInfo := fmt.Sprintf("\n[%d-%d of %d lines] h/l: scroll", 
				startLine+1, endLine, len(m.diffLines))
			diffContent += lipgloss.NewStyle().
				Foreground(lipgloss.Color("240")).
				Render(scrollInfo)
		}
	} else {
		diffContent += "\n\nNo changes to display"
	}

	diffPanel := diffStyle.Height(m.windowHeight - 6).Render(diffContent)

	// Combine both panels side by side
	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		fileList,
		"  ", // spacing
		diffPanel,
	)
}

func (m enhancedFileModel) formatDiffLines(lines []string) string {
	var formatted strings.Builder

	for i, line := range lines {
		if strings.HasPrefix(line, "+") && !strings.HasPrefix(line, "+++") {
			// Added lines in muted green
			formatted.WriteString(addedStyle.Render(line))
		} else if strings.HasPrefix(line, "-") && !strings.HasPrefix(line, "---") {
			// Removed lines in muted red
			formatted.WriteString(removedStyle.Render(line))
		} else if strings.HasPrefix(line, "@@") {
			// Hunk headers in cyan
			formatted.WriteString(hunkStyle.Render(line))
		} else {
			// Context lines in default color
			formatted.WriteString(contextStyle.Render(line))
		}
		
		// Add newline except for the last line to avoid extra spacing
		if i < len(lines)-1 {
			formatted.WriteString("\n")
		}
	}

	return formatted.String()
}

func (m enhancedFileModel) loadDiffForCurrentFile() tea.Cmd {
	if len(m.files) == 0 || m.cursor >= len(m.files) {
		return nil
	}

	filename := m.files[m.cursor]
	m.loadingDiff = true

	return tea.Cmd(func() tea.Msg {
		diff, err := git.GetSingleFileDiff(filename)
		if err != nil {
			return diffLoadedMsg("Error loading diff: " + err.Error())
		}
		return diffLoadedMsg(diff)
	})
}

type diffLoadedMsg string

func (m enhancedFileModel) GetSelectedFiles() []string {
	var selected []string
	for i, isSelected := range m.selected {
		if isSelected && i < len(m.files) {
			selected = append(selected, m.files[i])
		}
	}
	return selected
}

func (m enhancedFileModel) GetFilesToUnstage() []string {
	var toUnstage []string
	for i, file := range m.files {
		// Only unstage files that were:
		// 1. Initially staged (in initialStaged map)
		// 2. Currently unselected by user
		// 3. Actually marked as "staged" status
		if m.initialStaged[file] && !m.selected[i] && m.fileStatus[file] == "staged" {
			toUnstage = append(toUnstage, file)
		}
	}
	return toUnstage
}

func RunEnhancedFileSelector(files, stagedFiles, unstagedFiles, untrackedFiles []string) ([]string, []string, error) {
	model := NewEnhancedFileSelector(files, stagedFiles, unstagedFiles, untrackedFiles)
	p := tea.NewProgram(model, tea.WithAltScreen())
	
	finalModel, err := p.Run()
	if err != nil {
		return nil, nil, err
	}

	if m, ok := finalModel.(enhancedFileModel); ok {
		if m.quitted {
			return nil, nil, fmt.Errorf("user cancelled file selection")
		}
		return m.GetSelectedFiles(), m.GetFilesToUnstage(), nil
	}

	return nil, nil, fmt.Errorf("unexpected model type")
}