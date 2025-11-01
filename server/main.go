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
	res, ok := parseCommand(msg)
	if !ok {
		conn.Write([]byte("Could not parse message...\n"))
		return
	}
	conn.Write([]byte(res))
}

func parsePostId(id_s string) *Post {
	ids := strings.Split(id_s, "-")
	id, err := strconv.Atoi(ids[0])
	if err != nil {
		fmt.Println("Could not convert to integer: ", ids[0])
		return nil
	}
	if id >= len(posts) {
		fmt.Printf("Post %v does not exist\n", id)
		return nil
	}
	cur_post := &posts[id]
	for id_i := 1; id_i < len(ids); id_i += 1 {
		id := ids[id_i]
		id_val, err := strconv.Atoi(id)
		if err != nil {
			fmt.Println("Could not convert to integer: ", id)
			return nil
		}
		if id_val >= len(cur_post.comments) {
			fmt.Printf("Post %v does not exist\n", id_val)
			return nil
		}
		cur_post = &cur_post.comments[id_val]
	}
	return cur_post
}

func parseCommand(command string) (string, bool) {
	words := strings.Fields(command)
	var res string
	if len(words) < 2 {
		fmt.Println("Not a valid command")
		return "", false
	}
	switch words[0] {
	// example: POST <message>
	case "POST":
		msg := strings.Join(words[1:], " ")
		post := Post{id: uint64(len(posts)), content: msg}
		posts = append(posts, post)
		res = "Post created successfully..."
		break
	// example: FEED <page count>
	case "FEED":
		res = fmt.Sprintf(res, "%v", posts)
		break
	// example: FETCH <post_id-comment_id...>
	case "FETCH":
		if len(words) < 2 {
			fmt.Println("Not enough arguments to command ", words[0])
			return "", false
		}
		post := parsePostId(words[1])
		if post == nil {
			fmt.Println("Could not parse ", words[1])
			return "", false
		}
		res = fmt.Sprintf(res, "%v", post)
		break
	// example: COMMENT <id-comment_id...> <message>
	case "COMMENT":
		if len(words) < 3 {
			fmt.Println("Not enough arguments to command ", words[0])
			return "", false
		}
		msg := strings.Join(words[2:], " ")
		post := parsePostId(words[1])
		if post == nil {
			fmt.Println("Could not parse ", words[1])
			return "", false
		}
		fmt.Println(post)
		comment := Post{id: uint64(len(post.comments)), content: msg}
		post.comments = append(post.comments, comment)
		res = "Comment posted successfully..."
		break
	// example: LIKE <id-comment_id..>
	case "LIKE":
		post := parsePostId(words[1])
		if post == nil {
			fmt.Println("Could not parse ", words[1])
			return "", false
		}
		post.likes += 1
		res = "Post liked successfully..."
		break
	default:
		log.Fatalf("Unknown command: %s\n", words[0])
	}
	return res, true
}

func main() {
	list, err := net.Listen("tcp", PORT)
	if err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
	defer list.Close()
	fmt.Printf("Now accepting connections on port %v...\n", PORT)
	for {
		conn, err := list.Accept()
		if err != nil {
			fmt.Println("Could not accept incoming request: ", err)
			continue
		}
		go handleConnection(conn)
	}
}
