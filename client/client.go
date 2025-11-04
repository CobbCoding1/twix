package main

import (
	"bufio"
	"fmt"
	"log"
	"net"

	tea "github.com/charmbracelet/bubbletea"
)

// "FEED"
func HandleFeed(pageNum int) (string, error) {
	conn, _ := net.Dial("tcp", "localhost:8080")
	defer conn.Close()
	fmt.Fprintf(conn, "FEED %v\n", pageNum)
	return bufio.NewReader(conn).ReadString('\n')
}

// "POST"
func HandlePost(content string) (string, error) {
	conn, _ := net.Dial("tcp", "localhost:8080")
	defer conn.Close()
	fmt.Fprintf(conn, "POST %v\n", content)
	return bufio.NewReader(conn).ReadString('\n')
}

// "FETCH"
func HandleFetch(post string) (string, error) {
	conn, _ := net.Dial("tcp", "localhost:8080")
	defer conn.Close()
	fmt.Fprintf(conn, "FETCH %v\n", post)
	return bufio.NewReader(conn).ReadString('\n')
}

// "COMMENT"
func HandleComment(post, content string) (string, error) {
	conn, _ := net.Dial("tcp", "localhost:8080")
	defer conn.Close()
	fmt.Fprintf(conn, "COMMENT %v %v\n", post, content)
	return bufio.NewReader(conn).ReadString('\n')
}

// "LIKE"
func HandleLike(post string) (string, error) {
	conn, _ := net.Dial("tcp", "localhost:8080")
	defer conn.Close()
	fmt.Fprintf(conn, "LIKE %v\n", post)
	return bufio.NewReader(conn).ReadString('\n')
}

type model struct {
	cursor int
	width  int
	height int
	text   string
}

func (m model) Init() tea.Cmd {
	return nil
}

func initialModel() model {
	return model{
		text: "",
	}
}

func (m model) View() string {
	s := fmt.Sprintf("cursor %v\n", m.cursor)
	s += fmt.Sprintf("1. cobbcoding\nThis is a post.\n\n")
	s += fmt.Sprintf("2. stam\nThis is a post 2.\n\n")
	for range m.height - 8 {
		s += "\n"
	}
	s += fmt.Sprintf(":%v|", m.text)
	return s
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case tea.KeyMsg:
		switch msg.(tea.KeyMsg).String() {
		case "ctrl+c":
			return m, tea.Quit
		case "backspace":
			if m.cursor > 0 {
				m.cursor--
				m.text = m.text[:m.cursor]
			}
			break
		case "enter":
			m.text = ""
			m.cursor = 0
			break
		default:
			m.text += msg.(tea.KeyMsg).String()
			m.cursor += 1
		}
	case tea.WindowSizeMsg:
		m.width = msg.(tea.WindowSizeMsg).Width
		m.height = msg.(tea.WindowSizeMsg).Height
		break
	}
	return m, nil
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	_, err := p.Run()
	if err != nil {
		log.Fatal("Could not start TUI")
	}
}
