package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type choice struct {
	title  string
	action func()
}

type model struct {
	choices []choice
	cursor  int
}

// Styling
var (
	cursorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("205")).SetString(">")
	normalStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	selectedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("229")).Bold(true)
	headerStyle   = lipgloss.NewStyle().Bold(true).Underline(true).Foreground(lipgloss.Color("33"))
	footerStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Italic(true)
	outputStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Background(lipgloss.Color("235")).Padding(1, 2)
)

func InitialModel() model {
	return model{
		choices: []choice{
			{"Create Room", func() { fmt.Println(outputStyle.Render("You chose to create a room!")) }},
			{"Join Room", func() { fmt.Println(outputStyle.Render("You chose to join a room!")) }},
			{"Direct Message", func() { fmt.Println(outputStyle.Render("You chose to send a direct message!")) }},
			{"Chats", func() { fmt.Println(outputStyle.Render("You chose to view your chats!")) }},
			{"Friend Lists", func() { fmt.Println(outputStyle.Render("You chose to view your friend lists!")) }},
			{"Friend Requests", func() { fmt.Println(outputStyle.Render("You chose to view your friend requests!")) }},
			{"Settings", func() { fmt.Println(outputStyle.Render("You chose settings!")) }},
			{"Logout", func() { fmt.Println(outputStyle.Render("You chose to logout!")) }},
		},
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
		case "enter":
			m.choices[m.cursor].action()
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() string {
	s := headerStyle.Render("What should we buy at the market?") + "\n\n"

	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = cursorStyle.String()
		}

		title := normalStyle.Render(choice.title)
		if m.cursor == i {
			title = selectedStyle.Render(choice.title)
		}

		s += fmt.Sprintf("%s %s\n", cursor, title)
	}

	s += "\n" + footerStyle.Render("↑/↓ to navigate • Enter to select • q to quit") + "\n"
	return s
}
