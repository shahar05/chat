package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"time"
)

func main() {

	conn, err := net.Dial("tcp", "localhost:8888")

	for err != nil {
		conn, err = net.Dial("tcp", "localhost:8888")
		time.Sleep(5 * time.Second)
	}
	defer conn.Close()

	go listenToResponse(conn)

	stdinReader := bufio.NewReader(os.Stdin)
	fmt.Println(">> type /help")

	for {
		fmt.Println(">> waiting for command ")
		command, _ := stdinReader.ReadString('\n')
		fmt.Fprintf(conn, command+"\n")
	}

}

func listenToResponse(conn net.Conn) {
	connReader := bufio.NewReader(conn)
	for {
		response, err := connReader.ReadString('\n')
		if err != nil {
			errAndDie(err)
		}
		fmt.Println(response)
	}
}

func errAndDie(err error) {
	fmt.Println(err.Error())
	os.Exit(1)
}
