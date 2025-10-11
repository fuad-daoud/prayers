package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss/table"
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

type PrayerInfo struct {
	Results struct {
		Fajr    string `json:"Fajr"`
		Duha    string `json:"Duha"`
		Dhuhr   string `json:"Dhuhr"`
		Asr     string `json:"Asr"`
		Maghrib string `json:"Maghrib"`
		Isha    string `json:"Isha"`
	} `json:"results"`
	Settings struct {
		Name     string `json:"name"`
		Location struct {
			City    string `json:"city"`
			State   string `json:"state"`
			Country string `json:"country"`
		} `json:"location"`
		Latitude     string `json:"latitude"`
		Longitude    string `json:"longitude"`
		Timezone     string `json:"timezone"`
		Method       int    `json:"method"`
		Juristic     int    `json:"juristic"`
		HighLatitude int    `json:"high_latitude"`
		FajirRule    struct {
			Type  int `json:"type"`
			Value int `json:"value"`
		} `json:"fajir_rule"`
		MaghribRule struct {
			Type  int `json:"type"`
			Value int `json:"value"`
		} `json:"maghrib_rule"`
		IshaRule struct {
			Type  int `json:"type"`
			Value int `json:"value"`
		} `json:"isha_rule"`
		TimeFormat int `json:"time_format"`
	} `json:"settings"`
	Success bool `json:"success"`
}

func getPrayerInfo(ip string) (PrayerInfo, error) {
	resp, err := http.Get("https://islamicfinder.us/index.php/api/prayer_times?user_ip=" + ip)
	if err != nil {
		return PrayerInfo{}, err
	}
	defer resp.Body.Close()

	jsonString, err := io.ReadAll(resp.Body)
	if err != nil {
		return PrayerInfo{}, nil
	}
	var result PrayerInfo
	err = json.Unmarshal(jsonString, &result)
	if err != nil {
		fmt.Printf("Could not unmarshal the response: %v\n", err)
		fmt.Printf("jsonString: %v\n", jsonString)
		os.Exit(1)
	}
	return result, nil
}

type model struct {
	width  int
	height int
	ip     string
	info   PrayerInfo
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
	t := table.New()
	t.Headers("Prayer", "time")
	t.Row("Fajr", m.info.Results.Fajr)
	t.Row("Dhuhr", m.info.Results.Dhuhr)
	t.Row("Asr", m.info.Results.Asr)
	t.Row("Maghrib", m.info.Results.Maghrib)
	t.Row("Isha", m.info.Results.Isha)
	return t.Render() + "\n"

	// return lipgloss.Place(
	// 	m.width,
	// 	m.height,
	// 	lipgloss.Center,
	// 	lipgloss.Center,
	// 	content,
	// )
}

func main() {
	p := tea.NewProgram(InitialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Some error happened %v\n", err)
		os.Exit(1)
	}
}
