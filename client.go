package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

type Client struct {
	conn     net.Conn
	nick     string
	id       string
	room     *Room
	commands chan<- Command
}

func (c *Client) ReadInput() {
	for {
		msg, err := bufio.NewReader(c.conn).ReadString('\n')
		if err != nil {
			return
		}

		msg = strings.Trim(msg, "\r\n")

		args := strings.Split(msg, " ")
		cmd := strings.TrimSpace(args[0])

		switch cmd {
		case "/nick":
			c.commands <- Command{
				id:     CMD_NICK,
				client: c,
				args:   args,
			}
		case "/join":
			c.commands <- Command{
				id:     CMD_JOIN,
				client: c,
				args:   args,
			}
		case "/rooms":
			c.commands <- Command{
				id:     CMD_ROOMS,
				client: c,
			}
		case "/quit":
			c.commands <- Command{
				id:     CMD_QUIT,
				client: c,
			}
		case "/members":
			c.commands <- Command{
				id:     CMD_LISTMSG,
				client: c,
			}
		default:
			if strings.HasPrefix(cmd, "/") {
				c.err(fmt.Errorf("Unknown command: %s", cmd))
			}
			// FIXME message can be endlessly large
			c.commands <- Command{
				id:     CMD_MSG,
				client: c,
				args:   args,
			}
		}
	}
}

func (c *Client) err(err error) {
	c.conn.Write([]byte("err: " + err.Error() + "\n"))
}

func (c *Client) Msg(msg string) {
	c.conn.Write([]byte("> " + msg + "\n"))
}
