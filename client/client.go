package main

import (
	"bufio"
	"fmt"
	"net"
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

func main() {
	res, _ := HandlePost("This is a post")
	fmt.Println(res)
	res, _ = HandleComment("0", "This is a comment")
	fmt.Println(res)
	res, _ = HandleLike("0-0")
	fmt.Println(res)
	res, _ = HandleFeed(0)
	fmt.Println(res)
	res, _ = HandleFetch("0")
	fmt.Println(res)
}
