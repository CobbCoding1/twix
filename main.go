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
	msg, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Println("Could not read message from connection: ", err)
		return
	}
	ok := parseCommand(msg)
	if !ok {
		conn.Write([]byte("Could not parse message...\n"))
		return
	}
	conn.Write([]byte("Message recieved...\n"))
}

func parseCommand(command string) bool {
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
		id, err := strconv.Atoi(words[1])
		if err != nil {
			fmt.Println(words[0], " Could not convert to integer: ", words[1])
			return false
		}
		comment := Post{id: uint64(len(posts[id].comments)), content: msg}
		posts[id].comments = append(posts[id].comments, comment)
		break
	// example: LIKE <id>
	case "LIKE":
		id, err := strconv.Atoi(words[1])
		if err != nil {
			fmt.Println(words[0], " Could not convert to integer: ", words[1])
			return false
		}
		posts[id].likes += 1
		break
	default:
		log.Fatalf("Unknown command: %s\n", words[0])
	}
	return true
}

func main() {
	list, err := net.Listen("tcp", PORT)
	if err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
	defer list.Close()
	fmt.Println("Now accepting connections...")
	for {
		conn, err := list.Accept()
		if err != nil {
			fmt.Println("Could not accept incoming request: ", err)
			continue
		}
		go handleConnection(conn)
	}
}
