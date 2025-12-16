package bot

import (
	"fmt"
	"net"
)

func BotListener(IP string, PORT string) {

	Listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", IP, PORT))
	if err != nil {
		fmt.Println("failed to start bot listener")
		return
	}

	fmt.Println("bot listening on", Listener.Addr())

	for {
		Connection, err := Listener.Accept()
		if err != nil {
			continue
		}
		go BotHandler(Connection)
	}
}
