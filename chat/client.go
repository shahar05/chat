package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

var COMMAND_MAP map[string]CommandID = map[string]CommandID{
	"/join":  CMD_JOIN,
	"/rooms": CMD_ROOMS,
	"/nick":  CMD_NICK,
	"/quit":  CMD_QUIT,
	"/msg":   CMD_MSG,
	"/help":  CMD_HELP,
}

type Client struct {
	conn     net.Conn
	nickname string
	room     *Room
	commands chan<- Command
}

func NewClient(conn net.Conn, commands chan<- Command) *Client {
	return &Client{
		nickname: "anonymous",
		conn:     conn,
		commands: commands,
		room:     nil,
	}
}

func (c *Client) readInput() {
	for {
		msg, err := bufio.NewReader(c.conn).ReadString('\n')
		if err != nil {
			log.Printf("unable to read client message due to: %s", err.Error())
			return
		}

		log.Println(msg)

		msg = strings.Trim(msg, "\r\n")
		args := strings.Split(msg, " ")
		cmd := strings.TrimSpace(args[0])

		command, exists := COMMAND_MAP[cmd]
		if exists {
			c.commands <- *NewCommand(command, c, args)
		} else {
			c.err(fmt.Errorf("unknown command: %s ", cmd))
		}
		c.msg("Waiting for more commands")
	}
}

func (c *Client) err(err error) {
	c.msg("ERR: " + err.Error())
}

func (c *Client) msg(msg string) {
	c.conn.Write([]byte("> " + msg + "\n"))
}

func (c *Client) joinRoom(r *Room) {
	c.quitCurrentRoom() // leave current room before join
	c.room = r
}
