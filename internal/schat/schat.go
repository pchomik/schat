package schat

import (
	"log"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type errMsg error

type model struct {
	textarea textarea.Model
	err      error
}

func Run() {
	p := tea.NewProgram(initialModel())

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}

}

func initialModel() model {

	style := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#000")).
		Background(lipgloss.Color("243"))

	ti := textarea.New()
	// ti.Placeholder = "Once upon a time..."
	ti.ShowLineNumbers = false
	ti.Prompt = ""

	ti.FocusedStyle.Base = style
	ti.FocusedStyle.Text = style
	ti.FocusedStyle.CursorLine = style
	ti.FocusedStyle.Prompt = style
	ti.FocusedStyle.EndOfBuffer = style
	ti.FocusedStyle.Placeholder = style

	ti.BlurredStyle = ti.FocusedStyle

	ti.SetWidth(98)
	ti.SetHeight(3)
	ti.Focus()

	return model{
		textarea: ti,
		err:      nil,
	}
}

func (m model) Init() tea.Cmd {
	return textarea.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			if m.textarea.Focused() {
				m.textarea.Blur()
			}
		case tea.KeyCtrlJ:
			return m, tea.Quit
		case tea.KeyCtrlC:
			return m, tea.Quit
		default:
			if !m.textarea.Focused() {
				cmd = m.textarea.Focus()
				cmds = append(cmds, cmd)
			}
		}

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	m.textarea, cmd = m.textarea.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	style := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#000")).
		Background(lipgloss.Color("243")).
		Width(100).
		Padding(1)

	centeredStyle := lipgloss.NewStyle().
		Width(200).
		Align(lipgloss.Center)

	ti := style.Render(m.textarea.View())
	msg := lipgloss.NewStyle().
		Width(100).
		MarginTop(1).
		Render("(ctrl+enter to send, ctrl+c to quit)")
	return centeredStyle.Render(ti + "\n" + msg + "\n\n")
}
