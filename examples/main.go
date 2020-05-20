// This is an example program showing the usage of hellabot
package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"autronix.com/hellabot"
	log "gopkg.in/inconshreveable/log15.v2"
)

var serv = flag.String("server", os.Getenv("SERVER_HOST"), "hostname and port for irc server to connect to")
var nick = flag.String("nick", "hellabot", "nickname for the bot")

func main() {
	flag.Parse()

	botOptions := func(bot *hbot.Bot) {
		bot.SSL = true
		//bot.SASL = true
		bot.Password = os.Getenv("SERVER_PASS")
		bot.TLSConfig.InsecureSkipVerify = true
		bot.V3Cap = true
		bot.CapReqs = append([]string{":server-time", "message-tags"})
	}

	/*hijackSession := func(bot *hbot.Bot) {
		bot.HijackSession = true
	}*/
	channels := func(bot *hbot.Bot) {
		bot.Channels = strings.Split(os.Getenv("CHANNELS"), ",")
	}
	irc, err := hbot.NewBot(*serv, *nick, channels, botOptions)
	if err != nil {
		panic(err)
	}

	irc.AddTrigger(SayInfoMessage)
	irc.AddTrigger(LongTrigger)
	irc.Logger.SetHandler(log.StdoutHandler)
	// logHandler := log.LvlFilterHandler(log.LvlInfo, log.StdoutHandler)
	// or
	// irc.Logger.SetHandler(logHandler)
	// or
	// irc.Logger.SetHandler(log.StreamHandler(os.Stdout, log.JsonFormat()))

	// Start up bot (this blocks until we disconnect)
	irc.Run()
	fmt.Println("Bot shutting down.")
}

// This trigger replies Hello when you say hello
var SayInfoMessage = hbot.Trigger{
	func(bot *hbot.Bot, m *hbot.Message) bool {
		return m.Command == "PRIVMSG" && m.Content == "-info"
	},
	func(irc *hbot.Bot, m *hbot.Message) bool {
		irc.Reply(m, "Hello")
		return false
	},
}

// This trigger replies Hello when you say hello
var LongTrigger = hbot.Trigger{
	func(bot *hbot.Bot, m *hbot.Message) bool {
		return m.Command == "PRIVMSG" && m.Content == "-long"
	},
	func(irc *hbot.Bot, m *hbot.Message) bool {
		irc.Reply(m, "This is the first message")
		time.Sleep(5 * time.Second)
		irc.Reply(m, "This is the second message")

		return false
	},
}
