package schat

import (
	"log"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type errMsg error

type model struct {
	textarea      textarea.Model
	spinner       spinner.Model
	communication []string
	terminalWidth int
	sessionId     string
	processing    bool
	err           error
}

func Run() {
	m := initialModel()
	p := tea.NewProgram(m)

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

func initialModel() model {
	sp := spinner.New()
	sp.Spinner = spinner.Dot
	sp.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#838ba7"))

	ti := textarea.New()
	ti.ShowLineNumbers = false
	ti.Prompt = ""
	ti.SetHeight(3)
	ti.Focus()

	return model{
		textarea: ti,
		spinner:  sp,
		err:      nil,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(textarea.Blink, spinner.Tick)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// Store terminal width when window resizes
		m.terminalWidth = msg.Width - 4
		m.textarea.SetWidth(m.terminalWidth)

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlS:
			bodyStyle := lipgloss.NewStyle().
				Border(lipgloss.BlockBorder(), false, false, false, true).
				BorderForeground(lipgloss.Color("#838ba7")).
				Width(m.terminalWidth + 2).
				Padding(1)

			value := m.textarea.Value()
			if len(value) > 0 && value[len(value)-1] == '\n' {
				value = value[:len(value)-1]
			}
			m.communication = append(
				m.communication,
				bodyStyle.Render(value),
				"",
			)
			m.textarea.Reset()
			m.processing = true
			m.textarea.Blur()

			// TODO: Call remote API here to get data based on the user message
			// This is the best place because:
			// 1. The message has just been captured and stored in m.communication
			// 2. We have the raw user input from m.textarea.Value()
			// 3. The Bubbletea pattern supports async operations by returning commands
			// 4. We can return a tea.Cmd that performs the HTTP request and sends back a message with the result
			// 5. This keeps the UI responsive while the API call is in progress

		case tea.KeyCtrlQ:
			return m, tea.Quit
		case tea.KeyCtrlL:
			m.processing = false
			m.communication = []string{}
			m.textarea.Reset()
			return m, tea.Batch(append(cmds, m.textarea.Focus())...)
		case tea.KeyEsc:
			m.processing = false
			m.textarea.Reset()
			return m, tea.Batch(append(cmds, m.textarea.Focus())...)
		default:
			if m.textarea.Focused() {
				cmd = m.textarea.Focus()
				cmds = append(cmds, cmd)
			}
		}

	case errMsg:
		m.err = msg
		return m, nil
	}

	m.textarea, cmd = m.textarea.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if m.terminalWidth == 0 {
		return "Loading ..."
	}

	borderStyle := lipgloss.NewStyle().
		Padding(0, 1).
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("#838ba7"))

	processingStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#838ba7"))

	parts := m.communication[:]
	if m.processing {
		parts = append(
			parts,
			lipgloss.JoinHorizontal(
				lipgloss.Right, m.spinner.View(),
				processingStyle.Render("Processing..."),
			),
			"",
		)
	}
	parts = append(
		parts,
		"Your prompt:",
		borderStyle.Render(m.textarea.View()),
		"(ctrl+s to send, ctrl+q to quit)",
	)

	return lipgloss.JoinVertical(lipgloss.Top, parts...)
}
