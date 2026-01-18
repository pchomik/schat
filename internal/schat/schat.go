package schat

import (
	"fmt"
	"log"
	"schat/internal/providers"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
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
	provider      string
	err           error
	providerChan  chan providerResponse
}

type providerResponse struct {
	response string
	err      error
}

func waitForApiResponse(ch <-chan providerResponse) tea.Cmd {
	return func() tea.Msg {
		select {
		case res := <-ch:
			return res
		default:
			return nil
		}
	}
}

func Run(provider string) {
	m := initialModel(provider)
	p := tea.NewProgram(m)

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

func initialModel(provider string) model {
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
		provider: provider,
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

	case providerResponse:
		m.processing = false
		m.providerChan = nil // Clean up channel

		if msg.err != nil {
			m.err = msg.err
			m.communication = append(m.communication, "Error: "+msg.err.Error())
		} else {
			msg, err := glamour.Render(msg.response, "dark")
			if err != nil {
				m.err = err
				m.communication = append(m.communication, "Error: "+err.Error(), "")
			} else {
				m.communication = append(m.communication, strings.TrimSpace(msg), "")
			}
		}
		cmd = m.textarea.Focus()
		cmds = append(cmds, cmd)

	case spinner.TickMsg:
		var spinnerCmd tea.Cmd
		m.spinner, spinnerCmd = m.spinner.Update(msg)

		// If we are processing, keep checking the API channel
		if m.processing && m.providerChan != nil {
			cmds = append(cmds, waitForApiResponse(m.providerChan))
		}

		return m, tea.Batch(spinnerCmd, tea.Batch(cmds...))

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

			m.providerChan = make(chan providerResponse)

			go func() {
				switch m.provider {
				case "opencode-cli":
					response := providers.NewOpenCodeCli().Run(value)
					m.providerChan <- providerResponse{response: response}
				case "cursor":
					response := providers.NewOpenCodeCli().Run(value)
					m.providerChan <- providerResponse{response: response}
				default:
					m.providerChan <- providerResponse{err: fmt.Errorf("unknown provider: %s", m.provider)}
				}

			}()

			return m, tea.Batch(waitForApiResponse(m.providerChan), m.spinner.Tick)
		case tea.KeyCtrlQ:
			return m, tea.Quit
		case tea.KeyCtrlL:
			m.processing = false
			m.communication = []string{}
			m.textarea.Reset()
			return m, tea.Batch(append(cmds, m.textarea.Focus(), m.spinner.Tick)...)
		case tea.KeyEsc:
			m.processing = false
			m.textarea.Reset()
			return m, tea.Batch(append(cmds, m.textarea.Focus(), m.spinner.Tick)...)
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

	textStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#838ba7"))

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
		// textStyle.Render("Your prompt:"),
		borderStyle.Render(m.textarea.View()),
		textStyle.Render("| ctrl+s - send | ctrl+q - quit | ctrl+l - new | esc - clear |"),
	)

	return lipgloss.JoinVertical(lipgloss.Top, parts...)
}
