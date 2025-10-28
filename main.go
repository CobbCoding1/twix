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
	if len(words) < 2 {
		fmt.Println("Not a valid command")
		return false
	}
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
	// example: FETCH <post_id-comment_id...>
	case "FETCH":
		fmt.Println(posts)
		break
	// example: COMMENT <id-comment_id...> <message>
	case "COMMENT":
		if len(words) < 3 {
			fmt.Println("Not enough arguments to command ", words[0])
			return false
		}
		msg := strings.Join(words[2:], " ")
		id, err := strconv.Atoi(words[1])
		if err != nil {
			fmt.Println(words[0], " Could not convert to integer: ", words[1])
			return false
		}
		if id >= len(posts) {
			fmt.Printf("Post %v does not exist\n", id)
			return false
		}
		comment := Post{id: uint64(len(posts[id].comments)), content: msg}
		posts[id].comments = append(posts[id].comments, comment)
		break
	// example: LIKE <id-comment_id..>
	case "LIKE":
		id, err := strconv.Atoi(words[1])
		if err != nil {
			fmt.Println(words[0], " Could not convert to integer: ", words[1])
			return false
		}
		if id >= len(posts) {
			fmt.Printf("Post %v does not exist\n", id)
			return false
		}
		if len(words) == 3 {
			comment_id, err := strconv.Atoi(words[2])
			if err != nil {
				fmt.Println(words[0], " Could not convert to integer: ", words[2])
				return false
			}
			if comment_id >= len(posts[id].comments) {
				fmt.Printf("Comment %v on post %v does not exist\n", comment_id, id)
				return false
			}
			posts[id].comments[comment_id].likes += 1
			break
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
