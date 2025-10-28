package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
)

const (
	PORT = ":8080"
)

type Post struct {
	id       uint64
	content  string
	author   string
	likes    uint64
	comments []Post
}

var posts []Post

func handleConnection(conn net.Conn) {
	defer conn.Close()
	msg, _ := bufio.NewReader(conn).ReadString('\n')
	parseCommand(msg)
	conn.Write([]byte("Message recieved...\n"))
}

func parseCommand(command string) {
	words := strings.Fields(command)
	switch words[0] {
	// example: POST <message>
	case "POST":
		msg := strings.Join(words[1:], " ")
		post := Post{id: uint64(len(posts)), content: msg}
		posts = append(posts, post)
		fmt.Println(msg)
		break
	// example: FEED <page count>
	case "FEED":
		fmt.Println(posts)
		break
	// example: COMMENT <id> <message>
	case "COMMENT":
		msg := strings.Join(words[2:], " ")
		id, _ := strconv.Atoi(words[1])
		comment := Post{id: uint64(len(posts[id].comments)), content: msg}
		posts[id].comments = append(posts[id].comments, comment)
		break
	// example: LIKE <id>
	case "LIKE":
		id, _ := strconv.Atoi(words[1])
		posts[id].likes += 1
		break
	default:
		log.Fatalf("Unknown command: %s\n", words[0])
	}
}

func main() {
	list, _ := net.Listen("tcp", PORT)
	defer list.Close()
	fmt.Println("Now accepting connections...")
	for {
		conn, _ := list.Accept()
		go handleConnection(conn)
	}
}
