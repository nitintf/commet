package ui

import (
	"fmt"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/bitcs/commet/internal/config"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type configState int

const (
	mainMenu configState = iota
	aiSettings
	gitSettings
	enteringText
	confirmingSave
	selectingModel
)

type configModel struct {
	state         configState
	config        *config.Config
	menuItems     []string
	cursor        int
	textInput     string
	currentField  string
	showingInput  bool
	message       string
	previousState configState
	hasChanges    bool
	modelOptions  []string
}

func NewConfigModel(cfg *config.Config) configModel {
	return configModel{
		state:     mainMenu,
		config:    cfg,
		menuItems: []string{"AI Settings", "Git Settings", "Save & Exit"},
		cursor:    0,
	}
}

func (m configModel) Init() tea.Cmd {
	return nil
}

func (m configModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch m.state {
		case mainMenu:
			return m.updateMainMenu(msg)
		case aiSettings:
			return m.updateAISettings(msg)
		case gitSettings:
			return m.updateGitSettings(msg)
		case enteringText:
			return m.updateTextInput(msg)
		case confirmingSave:
			return m.updateConfirmation(msg)
		case selectingModel:
			return m.updateModelSelection(msg)
		}
	}
	return m, nil
}

func (m configModel) updateMainMenu(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}
	case "down", "j":
		if m.cursor < len(m.menuItems)-1 {
			m.cursor++
		}
	case "enter":
		switch m.cursor {
		case 0:
			m.state = aiSettings
			m.cursor = 0
		case 1:
			m.state = gitSettings
			m.cursor = 0
		case 2:
			if m.hasChanges {
				m.state = confirmingSave
			} else {
				return m, tea.Quit
			}
		}
	}
	return m, nil
}

func (m configModel) updateAISettings(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	aiItems := []string{"Provider (press Enter to cycle)", "API Key (press Enter to edit)", "Model (press Enter to edit)", "← Back"}

	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit
	case "esc":
		m.state = mainMenu
		m.cursor = 0
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}
	case "down", "j":
		if m.cursor < len(aiItems)-1 {
			m.cursor++
		}
	case "enter":
		switch m.cursor {
		case 0:
			return m.selectProvider()
		case 1:
			m.currentField = "API Key"
			m.textInput = m.config.AI.APIKey
			m.previousState = aiSettings
			m.state = enteringText
			m.showingInput = true
			return m, nil
		case 2:
			m.modelOptions = m.config.AI.GetAvailableModels()
			m.previousState = aiSettings
			m.state = selectingModel
			m.cursor = 0
			// Find current model index
			currentModel := m.config.AI.Model
			if currentModel == "" {
				currentModel = m.config.GetDefaultModel()
			}
			for i, model := range m.modelOptions {
				if model == currentModel {
					m.cursor = i
					break
				}
			}
			return m, nil
		case 3:
			m.state = mainMenu
			m.cursor = 0
		}
	}
	return m, nil
}

func (m configModel) selectProvider() (tea.Model, tea.Cmd) {
	providers := []config.Provider{config.ProviderOpenAI, config.ProviderClaude, config.ProviderGoogle, config.ProviderGroq}
	currentIndex := 0
	for i, p := range providers {
		if p == m.config.AI.Provider {
			currentIndex = i
			break
		}
	}

	nextIndex := (currentIndex + 1) % len(providers)
	m.config.AI.Provider = providers[nextIndex]

	m.config.AI.Model = ""

	m.message = fmt.Sprintf("Provider changed to %s", m.config.AI.Provider)
	m.hasChanges = true

	return m, nil
}

func (m configModel) updateGitSettings(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	gitItems := []string{"Auto Stage (press Enter to toggle)", "Show Diff (press Enter to toggle)", "Confirm Push (press Enter to toggle)", "Direct Commit (press Enter to toggle)", "Interactive Mode (press Enter to toggle)", "Use AI (press Enter to toggle)", "← Back"}

	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit
	case "esc":
		m.state = mainMenu
		m.cursor = 0
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}
	case "down", "j":
		if m.cursor < len(gitItems)-1 {
			m.cursor++
		}
	case "enter":
		switch m.cursor {
		case 0:
			m.config.Git.AutoStage = !m.config.Git.AutoStage
			m.message = fmt.Sprintf("Auto Stage %s", m.boolToString(m.config.Git.AutoStage))
			m.hasChanges = true
		case 1:
			m.config.Git.ShowDiff = !m.config.Git.ShowDiff
			m.message = fmt.Sprintf("Show Diff %s", m.boolToString(m.config.Git.ShowDiff))
			m.hasChanges = true
		case 2:
			m.config.Git.ConfirmPush = !m.config.Git.ConfirmPush
			m.message = fmt.Sprintf("Confirm Push %s", m.boolToString(m.config.Git.ConfirmPush))
			m.hasChanges = true
		case 3:
			m.config.Git.DirectCommit = !m.config.Git.DirectCommit
			m.message = fmt.Sprintf("Direct Commit %s", m.boolToString(m.config.Git.DirectCommit))
			m.hasChanges = true
		case 4:
			m.config.Git.Interactive = !m.config.Git.Interactive
			m.message = fmt.Sprintf("Interactive Mode %s", m.boolToString(m.config.Git.Interactive))
			m.hasChanges = true
		case 5:
			m.config.Git.UseAI = !m.config.Git.UseAI
			m.message = fmt.Sprintf("Use AI %s", m.boolToString(m.config.Git.UseAI))
			m.hasChanges = true
		case 6:
			m.state = mainMenu
			m.cursor = 0
		}
	}
	return m, nil
}

func (m configModel) updateModelSelection(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit
	case "esc":
		m.state = m.previousState
		m.cursor = 0
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}
	case "down", "j":
		if m.cursor < len(m.modelOptions)-1 {
			m.cursor++
		}
	case "enter":
		if m.cursor < len(m.modelOptions) {
			m.config.AI.Model = m.modelOptions[m.cursor]
			m.hasChanges = true
			m.state = m.previousState
			m.cursor = 0
		}
	}
	return m, nil
}

func (m configModel) updateTextInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit
	case "esc":
		m.state = m.previousState
		m.showingInput = false
		m.textInput = ""
	case "enter":
		m.applyTextInput()
		m.state = m.previousState
		m.showingInput = false
	case "backspace", "ctrl+h":
		if len(m.textInput) > 0 {
			m.textInput = m.textInput[:len(m.textInput)-1]
		}
	case "ctrl+u":
		m.textInput = ""
	case "ctrl+v", "cmd+v":
		// Paste from clipboard (works on both Windows/Linux and Mac)
		if clipboardText, err := clipboard.ReadAll(); err == nil {
			cleanText := cleanClipboardText(clipboardText)
			m.textInput += cleanText
		}
	case "ctrl+a", "cmd+a":
		// Select all (clear current input, will be replaced on next paste/type)
		m.textInput = ""
	default:
		// Accept all printable characters and common symbols
		key := msg.String()
		if len(key) == 1 && key >= " " && key <= "~" {
			m.textInput += key
		}
	}
	return m, nil
}

func (m configModel) applyTextInput() {
	switch m.currentField {
	case "API Key":
		m.config.AI.APIKey = m.textInput
		if m.textInput == "" {
			m.message = "API Key cleared"
		} else {
			m.message = "API Key updated"
		}
		m.hasChanges = true
	case "Model":
		m.config.AI.Model = m.textInput
		if m.textInput == "" {
			m.message = "Model reset to default"
		} else {
			m.message = fmt.Sprintf("Model set to %s", m.textInput)
		}
		m.hasChanges = true
	}
}

func (m configModel) updateConfirmation(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q", "n":
		return m, tea.Quit
	case "y", "enter":
		if err := m.config.Save(); err != nil {
			m.message = fmt.Sprintf("Error saving config: %v", err)
			m.state = mainMenu
			m.cursor = 0
		} else {
			m.message = "Configuration saved successfully!"
			m.hasChanges = false
			return m, tea.Quit
		}
	case "esc":
		m.state = mainMenu
		m.cursor = 0
	}
	return m, nil
}

func (m configModel) View() string {
	var s strings.Builder

	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("212")).
		Bold(true)

	highlightStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("212"))

	s.WriteString(titleStyle.Render("Commet Configuration") + "\n\n")

	switch m.state {
	case mainMenu:
		s.WriteString("Select configuration category:\n\n")
		for i, item := range m.menuItems {
			cursor := " "
			if m.cursor == i {
				cursor = ">"
				s.WriteString(fmt.Sprintf("%s %s\n",
					highlightStyle.Render(cursor),
					highlightStyle.Render(item)))
			} else {
				s.WriteString(fmt.Sprintf("%s %s\n", cursor, item))
			}
		}
		s.WriteString("\nUse ↑/↓ to navigate, Enter to select, Q to quit")
		if m.hasChanges {
			s.WriteString("\n⚠️ You have unsaved changes! Use 'Save & Exit' to persist them.")
		} else {
			s.WriteString("\n💡 Make changes and use 'Save & Exit' to persist them.")
		}

	case aiSettings:
		s.WriteString("AI Settings:\n\n")
		aiItems := []string{
			fmt.Sprintf("Provider: %s", m.config.AI.Provider),
			fmt.Sprintf("API Key: %s", m.config.MaskAPIKey()),
			fmt.Sprintf("Model: %s", m.getDisplayModel()),
			"← Back",
		}

		for i, item := range aiItems {
			cursor := " "
			if m.cursor == i {
				cursor = ">"
				s.WriteString(fmt.Sprintf("%s %s\n",
					highlightStyle.Render(cursor),
					highlightStyle.Render(item)))
			} else {
				s.WriteString(fmt.Sprintf("%s %s\n", cursor, item))
			}
		}
		s.WriteString("\nUse ↑/↓ to navigate, Enter to select/edit, Esc to go back")

	case gitSettings:
		s.WriteString("Git Settings:\n\n")
		gitItems := []string{
			fmt.Sprintf("Auto Stage: %s", m.boolToString(m.config.Git.AutoStage)),
			fmt.Sprintf("Show Diff: %s", m.boolToString(m.config.Git.ShowDiff)),
			fmt.Sprintf("Confirm Push: %s", m.boolToString(m.config.Git.ConfirmPush)),
			fmt.Sprintf("Direct Commit: %s", m.boolToString(m.config.Git.DirectCommit)),
			fmt.Sprintf("Interactive Mode: %s", m.boolToString(m.config.Git.Interactive)),
			fmt.Sprintf("Use AI: %s", m.boolToString(m.config.Git.UseAI)),
			"← Back",
		}

		for i, item := range gitItems {
			cursor := " "
			if m.cursor == i {
				cursor = ">"
				s.WriteString(fmt.Sprintf("%s %s\n",
					highlightStyle.Render(cursor),
					highlightStyle.Render(item)))
			} else {
				s.WriteString(fmt.Sprintf("%s %s\n", cursor, item))
			}
		}
		s.WriteString("\nUse ↑/↓ to navigate, Enter to toggle, Esc to go back")

	case enteringText:
		s.WriteString(fmt.Sprintf("Enter %s:\n\n", m.currentField))

		displayInput := m.textInput
		if m.currentField == "API Key" && len(m.textInput) > 0 {
			displayInput = strings.Repeat("*", len(m.textInput))
		}

		// Add cursor indicator
		s.WriteString(fmt.Sprintf("> %s_\n", displayInput))
		s.WriteString("\nType or paste (Ctrl+V/Cmd+V), Ctrl+A to clear, Enter to save, Esc to cancel")

	case selectingModel:
		s.WriteString("Select Model:\n\n")

		for i, model := range m.modelOptions {
			cursor := " "
			if m.cursor == i {
				cursor = ">"
				s.WriteString(fmt.Sprintf("%s %s\n",
					highlightStyle.Render(cursor),
					highlightStyle.Render(model)))
			} else {
				s.WriteString(fmt.Sprintf("%s %s\n", cursor, model))
			}
		}
		s.WriteString("\nUse ↑/↓ to navigate, Enter to select, Esc to go back")

	case confirmingSave:
		s.WriteString("Save Commet Configuration?\n\n")
		s.WriteString("Your changes will be saved to ~/.commet.yaml\n\n")
		s.WriteString("Press 'y' or Enter to save, 'n' or Esc to cancel")
	}

	return s.String()
}

func (m configModel) getDisplayModel() string {
	model := m.config.AI.Model
	if model == "" {
		return m.config.GetDefaultModel() + " (default)"
	}
	return model
}

func (m configModel) boolToString(b bool) string {
	if b {
		return "✓ enabled"
	}
	return "✗ disabled"
}

// cleanClipboardText removes unwanted characters from clipboard content
func cleanClipboardText(text string) string {
	// Remove newlines, carriage returns, and tabs
	cleanText := strings.ReplaceAll(text, "\n", "")
	cleanText = strings.ReplaceAll(cleanText, "\r", "")
	cleanText = strings.ReplaceAll(cleanText, "\t", "")

	// Only keep printable ASCII characters
	var filteredText string
	for _, char := range cleanText {
		if char >= ' ' && char <= '~' {
			filteredText += string(char)
		}
	}

	return filteredText
}

func RunConfigUI(cfg *config.Config) error {
	m := NewConfigModel(cfg)
	p := tea.NewProgram(m, tea.WithAltScreen())
	_, err := p.Run()
	return err
}
