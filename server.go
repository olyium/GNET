package main

import (
	"gnet/bot"
	"gnet/client"
)

var (
	IP          string = "localhost"
	CLIENT_PORT string = "141"
	BOT_PORT    string = "142"
	BINS_PORT   string = "143"
)

func main() {
	go client.ClientListener(IP, CLIENT_PORT)
	go bot.BotListener(IP, BOT_PORT)
	select {}
}
