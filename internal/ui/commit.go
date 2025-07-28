package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	choices  []string
	cursor   int
	selected map[int]struct{}
}

func initialModel() model {
	return model{
		choices: []string{
			"feat: add new feature",
			"fix: resolve bug in authentication",
			"docs: update README with installation steps",
			"refactor: improve code structure",
		},
		selected: make(map[int]struct{}),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "enter", " ":
			fmt.Printf("Selected: %s\n", m.choices[m.cursor])
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() string {
	s := "Choose a commit message:\n\n"

	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
			choice = lipgloss.NewStyle().Foreground(lipgloss.Color("212")).Render(choice)
		}
		s += fmt.Sprintf("%s %s\n", cursor, choice)
	}

	s += "\nPress q to quit, ↑/↓ to navigate, enter to select.\n"
	return s
}

func RunCommitUI() error {
	p := tea.NewProgram(initialModel())
	_, err := p.Run()
	return err
}