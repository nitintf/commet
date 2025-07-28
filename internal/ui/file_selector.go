package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type FileSelectorModel struct {
	files        []string
	selected     map[int]bool
	cursor       int
	finished     bool
	cancelled    bool
	selectedFiles []string
}

type keyMap struct {
	Up     key.Binding
	Down   key.Binding
	Select key.Binding
	Confirm key.Binding
	Cancel key.Binding
	Help   key.Binding
	Quit   key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Select},
		{k.Confirm, k.Cancel, k.Quit},
	}
}

var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	),
	Select: key.NewBinding(
		key.WithKeys(" ", "x"),
		key.WithHelp("space/x", "select"),
	),
	Confirm: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "confirm"),
	),
	Cancel: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "cancel"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c", "q"),
		key.WithHelp("q", "quit"),
	),
}

func NewFileSelectorModel(files []string) FileSelectorModel {
	return FileSelectorModel{
		files:    files,
		selected: make(map[int]bool),
		cursor:   0,
	}
}

func (m FileSelectorModel) Init() tea.Cmd {
	return nil
}

func (m FileSelectorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Quit):
			m.cancelled = true
			m.finished = true
			return m, tea.Quit

		case key.Matches(msg, keys.Cancel):
			m.cancelled = true
			m.finished = true
			return m, tea.Quit

		case key.Matches(msg, keys.Up):
			if m.cursor > 0 {
				m.cursor--
			}

		case key.Matches(msg, keys.Down):
			if m.cursor < len(m.files)-1 {
				m.cursor++
			}

		case key.Matches(msg, keys.Select):
			if m.selected[m.cursor] {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = true
			}

		case key.Matches(msg, keys.Confirm):
			// Collect selected files
			m.selectedFiles = []string{}
			for i, selected := range m.selected {
				if selected {
					m.selectedFiles = append(m.selectedFiles, m.files[i])
				}
			}
			m.finished = true
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m FileSelectorModel) View() string {
	if m.finished {
		return ""
	}

	var s strings.Builder

	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("212")).
		Bold(true)

	highlightStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("212"))

	selectedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("46"))

	s.WriteString(titleStyle.Render("📁 Select files to commit") + "\n\n")

	for i, file := range m.files {
		cursor := "  "
		checkbox := "☐"
		
		if m.cursor == i {
			cursor = "> "
		}
		
		if m.selected[i] {
			checkbox = "☑"
		}

		line := fmt.Sprintf("%s%s %s", cursor, checkbox, file)
		
		if m.cursor == i {
			if m.selected[i] {
				s.WriteString(selectedStyle.Render(line) + "\n")
			} else {
				s.WriteString(highlightStyle.Render(line) + "\n")
			}
		} else if m.selected[i] {
			s.WriteString(selectedStyle.Render(line) + "\n")
		} else {
			s.WriteString(line + "\n")
		}
	}

	s.WriteString("\n")
	s.WriteString("Navigate: ↑/↓  Select: space  Confirm: enter  Cancel: esc\n")
	
	selectedCount := len(m.selected)
	if selectedCount > 0 {
		s.WriteString(fmt.Sprintf("\n✨ %d file(s) selected", selectedCount))
	}

	return s.String()
}

func (m FileSelectorModel) GetSelectedFiles() []string {
	return m.selectedFiles
}

func (m FileSelectorModel) WasCancelled() bool {
	return m.cancelled
}

// RunFileSelector shows the interactive file selector and returns selected files
func RunFileSelector(files []string) ([]string, error) {
	if len(files) == 0 {
		return []string{}, fmt.Errorf("no files available for selection")
	}

	m := NewFileSelectorModel(files)
	p := tea.NewProgram(m, tea.WithAltScreen())
	
	finalModel, err := p.Run()
	if err != nil {
		return nil, fmt.Errorf("failed to run file selector: %w", err)
	}

	model := finalModel.(FileSelectorModel)
	
	if model.WasCancelled() {
		return nil, fmt.Errorf("file selection cancelled")
	}

	return model.GetSelectedFiles(), nil
}