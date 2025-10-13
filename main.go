package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

const MAIN_COLOR = "#FAFAFA"


type PrayerInfo struct {
	Code   int    `json:"code"`
	Status string `json:"status"`
	Data   struct {
		Timings struct {
			Fajr       string `json:"Fajr"`
			Sunrise    string `json:"Sunrise"`
			Dhuhr      string `json:"Dhuhr"`
			Asr        string `json:"Asr"`
			Sunset     string `json:"Sunset"`
			Maghrib    string `json:"Maghrib"`
			Isha       string `json:"Isha"`
			Imsak      string `json:"Imsak"`
			Midnight   string `json:"Midnight"`
			Firstthird string `json:"Firstthird"`
			Lastthird  string `json:"Lastthird"`
		} `json:"timings"`
	} `json:"data"`
}

func getPrayerInfo() (PrayerInfo, error) {
	resp, err := http.Get("https://api.aladhan.com/v1/timings/13-10-2025?latitude=31.9461222&longitude=35.923844&method=23&&timezonestring=Asia%2FAmman")
	if err != nil {
		return PrayerInfo{}, err
	}
	defer resp.Body.Close()

	jsonString, err := io.ReadAll(resp.Body)
	if err != nil {
		return PrayerInfo{}, nil
	}
	jsonString = []byte(strings.ReplaceAll(string(jsonString), "%", ""))
	var result PrayerInfo
	err = json.Unmarshal(jsonString, &result)
	if err != nil {
		fmt.Printf("Could not unmarshal the response: %v\n", err)
		fmt.Printf("jsonString: %v\n", jsonString)
		os.Exit(1)
	}
	if result.Code != 200 {
		fmt.Printf("Error integrating got none 200 status %v\n", result)
		os.Exit(1)
	}
	return result, nil
}

type model struct {
	width   int
	height  int
	info    PrayerInfo
	spin    spinner.Model
	loading bool
}

func fetchData() tea.Msg {
	info, err := getPrayerInfo()
	if err != nil {
		fmt.Printf("Could not get prayer info: %v\n", err)
		os.Exit(1)
	}
	return info
}

func InitialModel() model {
	spin := spinner.New(spinner.WithSpinner(spinner.Globe))
	spin.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("69"))
	spin.Spinner.FPS = time.Second / 10
	return model{spin: spin, loading: true}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		m.spin.Tick,
		fetchData,
	)
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
		default:
			return m, nil
		}
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spin, cmd = m.spin.Update(msg)
		return m, cmd
	case PrayerInfo:
		m.info = msg
		m.loading = false
		return m, nil
	}
	return m, nil
}

var helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render

func (m model) View() string {
	if m.loading {
		return fmt.Sprintf("%s fetching prayers times based on your location\n %s", m.spin.View(), helpStyle("• q/ctrl+c: exit\n"))
	}
	t := table.New()
	t.Headers("Prayer", "time")
	t.Row("Fajr", m.info.Data.Timings.Fajr)
	t.Row("Dhuhr", m.info.Data.Timings.Dhuhr)
	t.Row("Asr", m.info.Data.Timings.Asr)
	t.Row("Maghrib", m.info.Data.Timings.Maghrib)
	t.Row("Isha", m.info.Data.Timings.Isha)
	return t.Render() + "\n" + helpStyle("• q/ctrl+c: exit\n")

	// return lipgloss.Place(
	// 	m.width,
	// 	m.height,
	// 	lipgloss.Center,
	// 	lipgloss.Center,
	// 	content,
	// )
}

func main() {
	m := InitialModel()
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Some error happened %v\n", err)
		os.Exit(1)
	}
}
