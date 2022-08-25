package main

type CommandID int

const (
	CMD_NICK CommandID = iota
	CMD_JOIN
	CMD_ROOMS
	CMD_MSG
	CMD_QUIT
	CMD_HELP
)

type Command struct {
	id     CommandID
	client *Client
	args   []string
}

func NewCommand(id CommandID, client *Client, args []string) *Command {
	return &Command{
		id:     id,
		client: client,
		args:   args,
	}
}
