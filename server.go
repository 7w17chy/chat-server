package main

import (
	"fmt"
	"log"
	"time"
	"net"
	"strings"
	"github.com/segmentio/ksuid"
)

type server struct {
	rooms    map[string]*room
	commands chan command
	running  bool
}

func newServer() *server {
	return &server{
		rooms:    make(map[string]*room),
		commands: make(chan command),
		running:  true,
	}
}

func (s *server) init() {
	s.rooms["#welcome"] = &room{
		name: "#welcome",
		members: make(map[string]*client),
	}
}

func (s *server) run() {
	for cmd := range s.commands {
		switch cmd.id {
		case CMD_NICK:
			s.nick(cmd.client, cmd.args[1])
		case CMD_JOIN:
			s.join(cmd.client, cmd.args[1])
		case CMD_ROOMS:
			s.listRooms(cmd.client)
		case CMD_MSG:
			s.msg(cmd.client, cmd.args)
		case CMD_QUIT:
			s.quit(cmd.client)
		case CMD_LISTMSG:
			s.ListMembers(cmd.client)
		}
	}
}

func (s *server) shutdown() {
	// TODO
	return
}

func (s *server) newClient(conn net.Conn) {
	log.Printf("new client has joined: %s", conn.RemoteAddr().String())

	// generate new client id
	uid, err := ksuid.NewRandomWithTime(time.Now())
	if err != nil {
		log.Printf("Failed to create new client id!")
		// and then it will crash and burn... 'cause of nil
	}

	c := &client{
		conn:     conn,
		nick:     "anonymous",
		id:       uid.String(),
		commands: s.commands,
	}

	// put them in a default dedicated 'welcome'-room
	s.join(c, "#welcome")

	c.readInput()
}

func (s *server) nick(c *client, nick string) {
	c.nick = nick
	c.msg(fmt.Sprintf("all right, I will call you %s", nick))
}

func (s *server) join(c *client, roomName string) {
	r, ok := s.rooms[roomName]
	if !ok {
		r = &room{
			name:    roomName,
			members: make(map[string]*client),
		}
		s.rooms[roomName] = r
	}
	r.members[c.id] = c

	s.removeClientFromRoom(c)
	c.room = r

	r.broadcast(c, fmt.Sprintf("%s joined the room", c.nick))

	c.msg(fmt.Sprintf("welcome to %s", roomName))
}

// List all currently active rooms.
func (s *server) listRooms(c *client) {
	var rooms []string
	for name := range s.rooms {
		rooms = append(rooms, name)
	}

	c.msg(fmt.Sprintf("available rooms: %s", strings.Join(rooms, ", ")))
}

// List all current members of the room.
func (s *server) ListMembers(c *client) {
	var room = s.rooms[c.room.name]
	var members []string
	for _, nick := range room.Members() {
		members = append(members, nick)
	}

	c.msg(fmt.Sprintf("Members of the current room: %s", strings.Join(members[1:], ", ")))
}

// Send a message to all members of the room the sending user is currently in.
func (s *server) msg(c *client, args []string) {
	//msg := strings.Join(args[1:len(args)], " ")
	msg := strings.Join(args[:len(args)], " ")
	c.room.broadcast(c, c.nick+": "+msg)
}

// Leave the current room. All chat data for the current user will be lost.
func (s *server) quit(c *client) {
	log.Printf("client has left the chat: %s", c.id)

	s.removeClientFromRoom(c)

	c.msg("sad to see you go =(")
	c.conn.Close()
}

// Remove specified user from their current room.
func (s *server) removeClientFromRoom(c *client) {
	if c.room != nil {
		oldRoom := s.rooms[c.room.name]
		delete(s.rooms[c.room.name].members, c.id)
		oldRoom.broadcast(c, fmt.Sprintf("%s has left the room", c.nick))
	}
}
