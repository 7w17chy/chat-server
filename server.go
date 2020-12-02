package main

import (
	"fmt"
	"log"
	"time"
	"net"
	"strings"
	"github.com/segmentio/ksuid"
)

type Server struct {
	rooms    map[string]*Room
	commands chan Command
	running  bool
}

func NewServer() *Server {
	return &Server{
		rooms:    make(map[string]*Room),
		commands: make(chan Command),
		running:  true,
	}
}

func (s *Server) Init() {
	s.rooms["#welcome"] = &Room{
		name: "#welcome",
		members: make(map[string]*Client),
	}
}

func (s *Server) Run() {
	for cmd := range s.commands {
		switch cmd.id {
		case CMD_NICK:
			s.Nick(cmd.client, cmd.args[1])
		case CMD_JOIN:
			s.Join(cmd.client, cmd.args[1])
		case CMD_ROOMS:
			s.ListRooms(cmd.client)
		case CMD_MSG:
			s.Msg(cmd.client, cmd.args)
		case CMD_QUIT:
			s.Quit(cmd.client)
		case CMD_LISTMSG:
			s.ListMembers(cmd.client)
		case CMD_SHTDWN:
			s.Shutdown()
		}
	}
}

// Shut down server gracefully.
func (server *Server) Shutdown() {
	// Close all client's connections
	for _, room := range server.rooms {
		room.GeneralMessage("Server is shutting down.")
		for _, client := range room.members {
			server.Kick(client)
		}
	}

	// Exit superloop
	server.running = false
}

func (s *Server) NewClient(conn net.Conn) {
	log.Printf("new client has joined: %s", conn.RemoteAddr().String())

	// generate new client id
	uid, err := ksuid.NewRandomWithTime(time.Now())
	if err != nil {
		log.Printf("Failed to create new client id!")
		// and then it will crash and burn... 'cause of nil
	}

	c := &Client{
		conn:     conn,
		nick:     "anonymous",
		id:       uid.String(),
		commands: s.commands,
	}

	// put them in a default dedicated 'welcome'-room
	s.Join(c, "#welcome")

	c.ReadInput()
}

func (s *Server) Nick(c *Client, nick string) {
	c.nick = nick
	c.Msg(fmt.Sprintf("all right, I will call you %s", nick))
}

func (s *Server) Join(c *Client, roomName string) {
	r, ok := s.rooms[roomName]
	if !ok {
		r = &Room{
			name:    roomName,
			members: make(map[string]*Client),
		}
		s.rooms[roomName] = r
	}
	r.members[c.id] = c

	s.removeClientFromRoom(c)
	c.room = r

	r.Broadcast(c, fmt.Sprintf("%s joined the room", c.nick))

	c.Msg(fmt.Sprintf("welcome to %s", roomName))
}

// List all currently active rooms.
func (s *Server) ListRooms(c *Client) {
	var rooms []string
	for name := range s.rooms {
		rooms = append(rooms, name)
	}

	c.Msg(fmt.Sprintf("available rooms: %s", strings.Join(rooms, ", ")))
}

// List all current members of the room.
func (s *Server) ListMembers(c *Client) {
	var room = s.rooms[c.room.name]
	var members []string
	for _, nick := range room.Members() {
		members = append(members, nick)
	}

	c.Msg(fmt.Sprintf("Members of the current room: %s", strings.Join(members[1:], ", ")))
}

// Send a message to all members of the room the sending user is currently in.
func (s *Server) Msg(c *Client, args []string) {
	//msg := strings.Join(args[1:len(args)], " ")
	msg := strings.Join(args[:len(args)], " ")
	c.room.Broadcast(c, c.nick+": "+msg)
}

// Leave the current room. All chat data for the current user will be lost.
func (s *Server) Quit(c *Client) {
	log.Printf("client has left the chat: %s", c.id)

	s.removeClientFromRoom(c)

	c.Msg("sad to see you go =(")
	c.conn.Close()
}

// Remove client from their current room and kick them from the server by closing their connection.
func (server *Server) Kick(client *Client) {
	server.removeClientFromRoom(client)
	client.conn.Close()
}

// Remove specified user from their current room.
func (s *Server) removeClientFromRoom(c *Client) {
	if c.room != nil {
		oldRoom := s.rooms[c.room.name]
		delete(s.rooms[c.room.name].members, c.id)
		oldRoom.Broadcast(c, fmt.Sprintf("%s has left the room", c.nick))
	}
}
