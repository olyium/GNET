package main

import (
	"gnet/bins"
	"gnet/bot"
	"gnet/client"
)

var (
	IP          = "localhost"
	CLIENT_PORT = "5141"
	BOT_PORT    = "5142"
	BINS_PORT   = "5143"
)

func main() {
	go client.ClientListener(IP, CLIENT_PORT)
	go bot.BotListener(IP, BOT_PORT)
	go bins.BinListener(IP, BINS_PORT)
	select {}
}
