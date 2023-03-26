package main

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/tyrkinn/lua_tg_bot/commands"
	lua "github.com/yuin/gopher-lua"
	tele "gopkg.in/telebot.v3"
)

func DefineBotCommand(bot *tele.Bot, L *lua.LState, commandName string, luaCode string) {
	bot.Handle("/"+commandName, func(c tele.Context) error {
		commands.AddCommandsToState(L, c)
		src := luaCode
		if err := L.DoString(src); err != nil {
			c.Send("Error while trying to execute lua command: " + err.Error())
		}
		return nil
	})
}

func main() {

	L := lua.NewState()
	defer L.Close()

	pref := tele.Settings{
		Token:  os.Getenv("TG_BOT_TOKEN"),
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}
	b, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}

	b.Handle("/hello", func(c tele.Context) error {
		commands.AddCommandsToState(L, c)
		src := "reply('Привет, залупа')"
		if err := L.DoString(src); err != nil {
			c.Send("Error while trying to execute lua command: " + err.Error())
		}
		return nil
	})

	b.Handle("/define", func(c tele.Context) error {
		commands.AddCommandsToState(L, c)
		args := c.Args()
		name := args[0]
		luaCommand := strings.Join(args[1:], " ")
		if len(args) < 2 {
			c.Send("You should provide 2 arguments: Command Name And Lua Code To It")
		}
		DefineBotCommand(b, L, name, luaCommand)
		return nil
	})

	b.Start()

}
