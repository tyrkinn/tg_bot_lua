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

func DefineBotCommand(ctx *commands.TgLuaContext, bot *tele.Bot, commandName string, luaCode string) {
	L := ctx.L
	bot.Handle("/"+commandName, func(c tele.Context) error {
		ctx.AttachToTGCtx(c)
		src := luaCode
		if err := L.DoString(src); err != nil {
			c.Send("Error while trying to execute lua command: " + err.Error())
		}
		return nil
	})
}

func main() {

	pref := tele.Settings{
		Token:  os.Getenv("TG_BOT_TOKEN"),
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	LTCtx := commands.TgLuaContext{L: lua.NewState()}

	b, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}

	b.Handle("/hello", func(c tele.Context) error {
		LTCtx.AttachToTGCtx(c)
		if err := LTCtx.ExecLua("reply('Hi')"); err != nil {
			c.Send("Error while trying to execute lua command: " + err.Error())
		}
		return nil
	})

	b.Handle("/define", func(c tele.Context) error {
		LTCtx.AttachToTGCtx(c)
		args := c.Args()
		name := args[0]
		luaCommand := strings.Join(args[1:], " ")
		if len(args) < 2 {
			c.Send("You should provide 2 arguments: Command Name And Lua Code To It")
		}
		DefineBotCommand(&LTCtx, b, name, luaCommand)
		LTCtx.ExecLua("reply('Command `/" + name + "` created')")
		return nil
	})

	b.Start()

}
