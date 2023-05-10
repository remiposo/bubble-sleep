package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	fps      = 30
	padding  = 2
	maxWidth = 80
)

var helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).Render

type model struct {
	progress progress.Model
	sleep    time.Duration
	percent  float64
}

type tickMsg time.Time

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second/fps, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

type sleepMsg time.Time

func sleepCmd(d time.Duration) tea.Cmd {
	return tea.Tick(d, func(t time.Time) tea.Msg {
		return sleepMsg(t)
	})
}

func (m model) Init() tea.Cmd {
	return tea.Batch(sleepCmd(m.sleep), tickCmd())
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m, tea.Quit
	case tea.WindowSizeMsg:
		m.progress.Width = msg.Width - padding*2 - 4
		if m.progress.Width > maxWidth {
			m.progress.Width = maxWidth
		}
		return m, nil
	case tickMsg:
		var cmd tea.Cmd
		if m.percent <= 0.0 {
			cmd = nil
		} else {
			cmd = tickCmd()
			m.percent -= 1.0 / fps / m.sleep.Seconds()
		}
		return m, cmd
	case sleepMsg:
		// since tickCmd is delayed, percent will not be 0
		m.percent = 0.0
		return m, tea.Quit
	default:
		return m, nil
	}
}

func (m model) View() string {
	if m.percent <= 0.0 {
		return ""
	}
	pad := strings.Repeat(" ", padding)
	return "\n" +
		pad + m.progress.ViewAs(m.percent) + "\n\n" +
		pad + helpStyle("Press any key to exit") + "\n"
}

func main() {
	flag.Parse()
	if flag.NArg() != 1 {
		fmt.Println("Usage: sleep <seconds>")
		os.Exit(1)
	}
	sleep, err := strconv.Atoi(flag.Arg(0))
	if err != nil {
		fmt.Println("Usage: sleep <seconds>")
		os.Exit(1)
	}
	m := model{
		sleep: time.Second * time.Duration(sleep),
		progress: progress.NewModel(
			progress.WithDefaultGradient(),
			progress.WithoutPercentage(),
		),
		percent: 1.0,
	}
	if err := tea.NewProgram(m).Start(); err != nil {
		fmt.Println("Unexpected Error!!:", err)
		os.Exit(1)
	}
}
