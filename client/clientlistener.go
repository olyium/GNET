package client

import (
	"fmt"
	"net"
)

func ClientListener(IP string, PORT string) {

	ClientListener, ERR := net.Listen("tcp", fmt.Sprintf("%s:%s", IP, PORT))
	if ERR != nil {
		fmt.Print("failed to start client listener\n")
	}

	fmt.Print("client listening on ", ClientListener.Addr(), "\n")

	for {
		Connection, ERR := ClientListener.Accept()
		if ERR != nil {
			fmt.Print("failed to accept connection from client\n")
		}
		fmt.Print("new client connection\n")
		go ClientLoginHandler(Connection)
	}

}
