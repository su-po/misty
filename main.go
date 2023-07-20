package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/joho/godotenv"
)

type Raindrop struct {
	ID     int32  `json:"_id"`
	Format string `json:"type"`
	Title  string `json:"title"`
	URL    string `json:"link"`
}

type RaindropCollection struct {
	Result bool       `json:"result"`
	Items  []Raindrop `json:"items"`
}

type Raindrops []Raindrop

type model struct {
	raindrops Raindrops
	cursor    int // current feed position
	target    map[int]struct{}
}

func open(url string) {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("xdg-open", url)

	case "darwin":
		cmd = exec.Command("open", url)

	case "windows":
		cmd = exec.Command("cmd", "/c", "start")
	default:
		fmt.Println("I may not handle your case yet. Please let me know!")
		return

	}

	//cmd := exec.Command("open", m.raindrops[m.cursor].URL)

	err := cmd.Start()

	cmd.Wait()
	if err != nil {
		fmt.Println("Issue occurred during attempt to open url\n", err)
		//				}

	}

	// :FIXME: We should probably return the error here.
}

func buildModel() model {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	token := string(os.Getenv("APITOKEN"))
	raindrops := getRaindrops(token)
	return model{
		raindrops: raindrops,
		cursor:    0,
		target:    make(map[int]struct{}), // it should capture what you've selected and open it up in your web browser?
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
				m.cursor -= 1
			}

		case "down", "j":
			if m.cursor < len(m.raindrops)-1 {
				m.cursor += 1
			}

		case " ":
			_, ok := m.target[m.cursor]
			if ok {
				delete(m.target, m.cursor)
			} else {
				m.target[m.cursor] = struct{}{}
			}
		case "enter":
			_, ok := m.target[m.cursor]
			if ok {
				cmd := exec.Command("open", m.raindrops[m.cursor].URL)

				err := cmd.Start()

				if err != nil {
					fmt.Println("Issue occurred during attempt to open url\n", err)
				}
			} else {
				m.target[m.cursor] = struct{}{}
			}
		}
	}
	return m, nil
}

func (m model) View() string {
	// HEADER
	header := "Articles\n"
	header += "Quit -> ctrl + c & q | Select -> space | Open Link -> enter\n"

	for i, article := range m.raindrops {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		// If selected
		checked := " "
		if _, ok := m.target[i]; ok {

			checked = "ø"
		}
		// render
		header += fmt.Sprintf("%s [%s] %s\n", cursor, checked, article.Title)

		if i > 10 {
			break
		}
	}
	return header
}
func getRaindrops(token string) Raindrops {

	var items []Raindrop
	cursor := 0
	for {
		url := "https://api.raindrop.io/rest/v1/raindrops/0?search=type%3Aarticle&perpage=50&page=" + fmt.Sprintf("%d", cursor)
		req, _ := http.NewRequest("GET", url, nil)
		headerValue := "Bearer " + token
		req.Header.Add("Authorization", headerValue)

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Fatal("Something went wrong\n." + err.Error())

		}

		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Fatal("Failed to read the response body")
		}

		var raindropCollection RaindropCollection
		err = json.Unmarshal(body, &raindropCollection)
		if err != nil {
			log.Fatal("Program failed.", err)
		}

		if len(raindropCollection.Items) < 1 {
			break
		}
		items = append(items, raindropCollection.Items...)
		cursor += 1
	}
	return items
	// j, _ := json.MarshalIndent(items, "", " ")
}

func main() {
	program := tea.NewProgram(buildModel())

	if err := program.Start(); err != nil {
		fmt.Println("An error occurred.", err)

	}
	// fmt.Println(string(j))
}
