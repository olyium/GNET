package client

import (
	"fmt"
	"net"
)

func ClientListener(IP string, PORT string) {

	Listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", IP, PORT))
	if err != nil {
		fmt.Println("failed to start client listener")
		return
	}

	fmt.Println("client listening on", Listener.Addr())

	for {
		Connection, err := Listener.Accept()
		if err != nil {
			continue
		}
		go ClientHandler(Connection)
	}
}
