package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const MAIN_COLOR = "#FAFAFA"

func getPublicIP() (string, error) {
	resp, err := http.Get("https://api.ipify.org?format=text")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	ip, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(ip)), nil
}
func getPrayerInfo(ip string) (string, error) {
	resp, err := http.Get("https://islamicfinder.us/index.php/api/prayer_times?user_ip=" + ip)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	info, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", nil
	}
	return string(info), nil
}

type model struct {
	width  int
	height int
	ip     string
	info   string
}

func InitialModel() model {
	ip, err := getPublicIP()
	if err != nil {
		fmt.Printf("Could no get public ip: %v\n", err)
		os.Exit(1)
	}
	info, err := getPrayerInfo(ip)
	if err != nil {
		fmt.Printf("Could not get prayer info: %v\n", err)
		os.Exit(1)
	}
	return model{ip: ip, info: info}
}

func (model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m model) View() string {
	content := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(MAIN_COLOR)).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(MAIN_COLOR)).
		Padding(1).
		Width(50).
		Align(lipgloss.Center).
		Render(m.info)

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		content,
	)
}

func main() {
	p := tea.NewProgram(InitialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Some error happened %v\n", err)
		os.Exit(1)
	}
}
