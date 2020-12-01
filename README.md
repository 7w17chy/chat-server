# chat-server
Forked from github.com/plutov/packagemain/tree/master/20-tcp-chat

A (very) simple chat server written in Go.

# Changes introduced by this fork
## For upcoming changes/features, please take a look at the TODO file!
## Commands
- /members
- /msg text - command is no longer necessary. If not prefixed by a '/', text will be handled as message not as a command.
# User seperation by unique identifier rather than their IP-Address
# Commands
- /quit - quits the current room, and if in no room, ends the connection
- /join name - joins a room or creates it if it doesn't exist
- /rooms - lists all available rooms
- /nick name - changes the user's nick name ("anonymous" by default)
- /members - lists all users of the room the issuing user is in
