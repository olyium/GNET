package bot

import (
	"fmt"
	"net"
)

type Bot struct {
	Addr string
	Conn net.Conn
}

var Bots []Bot

func BotHandler(Connection net.Conn) {

	for _, Bot := range Bots {
		if Bot.Addr == Connection.RemoteAddr().String() {
			fmt.Print("bot attempted to duplicate connection so we close now\n")
			Connection.Close()
		}
	}

}
