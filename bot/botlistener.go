package bot

import (
	"fmt"
	"net"
)

func BotListener(IP string, PORT string) {

	BotListener, ERR := net.Listen("tcp", fmt.Sprintf("%s:%s", IP, PORT))
	if ERR != nil {
		fmt.Print("failed to start bot listener\n")
	}

	fmt.Print("bot listening on ", BotListener.Addr(), "\n")

	for {
		Connection, ERR := BotListener.Accept()
		if ERR != nil {
			fmt.Print("failed to accept connection from bot\n")
		}
		fmt.Print("new bot connection\n")
		go BotHandler(Connection)
	}

}
