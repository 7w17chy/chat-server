package main

type Room struct {
	name    string
	members map[string]*Client
}

// Returns a string slice with the nick names of all current members of the room.
func (room *Room) Members() []string {
	clientNicks := make([]string, 1)
	for _, client := range room.members {
		clientNicks = append(clientNicks, client.nick)
	}

	return clientNicks
}

// Send message to all members of the room.
func (r *Room) Broadcast(sender *Client, msg string) {
	for id, m := range r.members {
		// don't send message to sender
		if sender.id != id {
			m.Msg(msg)
		}
	}
}
