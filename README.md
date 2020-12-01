# chat-server
Forked from github.com/plutov/packagemain/tree/master/20-tcp-chat

A (very) simple chat server written in Go.

# Changes introduced by this fork
## Commands
- /members
### WIP
- Get rid of the '/msg' command. Sending a message should be the norm/default, commands the exception
# User seperation by unique identifier rather than their IP-Address
# Commands
- /msg text - sends a message to all members of the room the issuing user is in
- /quit - ends the connection
- /join name - joins a room or creates it if it doesn't exist
- /rooms - lists all available rooms
- /nick name - changes the user's nick name ("anonymous" by default)
- /members - lists all users of the room the issuing user is in
