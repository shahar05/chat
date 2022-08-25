package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"reflect"
	"strconv"
	"strings"
)

type Server struct {
	rooms    map[string]*Room
	commands chan Command
}

func newServer() *Server {
	return &Server{
		rooms:    map[string]*Room{},
		commands: make(chan Command),
	}
}

func (s *Server) run() {
	for cmd := range s.commands {
		switch cmd.id {
		case CMD_NICK:
			s.nick(cmd.client, cmd.args)
		case CMD_MSG:
			s.msg(cmd.client, cmd.args)
		case CMD_ROOMS:
			s.listRooms(cmd.client, cmd.args)
		case CMD_QUIT:
			s.quit(cmd.client, cmd.args)
		case CMD_JOIN:
			s.join(cmd.client, cmd.args)
		case CMD_HELP:
			s.help(cmd.client, cmd.args)

		}
	}
}
func (s *Server) handleConnection(conn net.Conn) {
	log.Printf("new Client connection: %s", conn.RemoteAddr().String())
	c := NewClient(conn, s.commands)
	c.readInput()
}

func (s *Server) help(c *Client, args []string) {
	c.msg(createList(COMMAND_MAP, "List of commands:"))
}

func (s *Server) nick(c *Client, args []string) {
	c.nickname = args[1]
	c.msg(fmt.Sprintf("all right, I will call you %s", c.nickname))
}

func (s *Server) join(c *Client, args []string) {
	roomName := args[1]
	r, exists := s.rooms[roomName]
	if !exists {
		r = newRoom(roomName)
		s.rooms[roomName] = r
	}
	//Add member
	r.members[c.conn.RemoteAddr()] = c
	c.joinRoom(r)
	r.broadcast(c, fmt.Sprintf("%s has joined the room", c.nickname))
	c.msg(fmt.Sprintf("Welcome to %s", r.name))

}
func createList(genericMap interface{}, title string) string {
	var sb strings.Builder
	sb.WriteString(title)
	sb.WriteString("\n")

	mapKeys := reflect.ValueOf(genericMap).MapKeys()
	i := 0

	for _, k := range mapKeys {
		i++
		sb.WriteString(strconv.Itoa(i) + ". ")
		sb.WriteString(k.String())
		sb.WriteString("\n")

	}

	return sb.String()
}

// func createList(strList []string, title string) string {
// 	var sb strings.Builder
// 	sb.WriteString(title)
// 	sb.WriteString("\n")
// 	for i, str := range strList {
// 		sb.WriteString(strconv.Itoa(i) + ". ")
// 		sb.WriteString(str)
// 		sb.WriteString("\n")
// 	}

// 	return sb.String()
// }

func (s *Server) listRooms(c *Client, args []string) {
	c.msg(createList(s.rooms, "List of rooms:"))
}

func (s *Server) msg(c *Client, args []string) {
	if c.room == nil {
		c.err(errors.New("you must join the room first"))
		return
	}
	c.room.broadcast(c, c.nickname+": "+strings.Join(args[1:], " "))
}
func (s *Server) quit(c *Client, args []string) {
	log.Printf("client has disconnected %s", c.conn.RemoteAddr().String())
	c.quitCurrentRoom()
	c.conn.Close()
}

func (c *Client) quitCurrentRoom() {
	if c.room != nil {
		delete(c.room.members, c.conn.RemoteAddr())
		c.room.broadcast(c, fmt.Sprintf("%s has left the room", c.nickname))
	}
}
