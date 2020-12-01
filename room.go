package main

type room struct {
	name    string
	members map[string]*client
}

// Returns a string slice with the nick names of all current members of the room.
func (room *room) Members() []string {
	clientNicks := make([]string, 1)
	for _, client := range room.members {
		clientNicks = append(clientNicks, client.nick)
	}

	return clientNicks
}

// Send message to all members of the room.
func (r *room) broadcast(sender *client, msg string) {
	for id, m := range r.members {
		// don't send message to sender
		if sender.id != id {
			m.msg(msg)
		}
	}
}
