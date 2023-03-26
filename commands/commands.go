package commands

import (
	"io"
	"net/http"

	lua "github.com/yuin/gopher-lua"
	tele "gopkg.in/telebot.v3"
)

func sendMessage(tgContext tele.Context) lua.LGFunction {
	return func(L *lua.LState) int {
		message := L.CheckString(1)
		tgContext.Send(message)
		return 0
	}
}

func replyWithText(tgContext tele.Context) lua.LGFunction {
	return func(L *lua.LState) int {
		message := L.CheckString(1)
		tgContext.Reply(message)
		return 0
	}
}

func httpget(tgContext tele.Context) lua.LGFunction {
	return func(L *lua.LState) int {
		uri := L.CheckString(1)
		resp, err := http.Get(uri)
		if err != nil {
			tgContext.Send("Error occured while http call")
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		L.Push(lua.LString(body))
		return 1
	}
}

type Command = func(tele.Context) lua.LGFunction

func attachCommand(L *lua.LState, commandName string, command Command, tContext tele.Context) {

	fn := L.NewFunction(command(tContext))
	L.SetGlobal(commandName, fn)
}

func AddCommandsToState(L *lua.LState, tgContext tele.Context) {

	attachCommand(
		L, "say", sendMessage, tgContext,
	)

	attachCommand(
		L, "reply", replyWithText, tgContext,
	)

	attachCommand(
		L, "http", httpget, tgContext,
	)

}
