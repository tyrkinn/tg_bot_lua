package commands

import (
	"io"
	"net/http"

	lua "github.com/yuin/gopher-lua"
	tele "gopkg.in/telebot.v3"
)

type TgLuaContext struct {
	Funcs map[string]lua.LGFunction
	Tele  tele.Context
	L     *lua.LState
}

func (c *TgLuaContext) SetupGlobals() {
	c.Funcs = map[string]lua.LGFunction{
		"say": func(L *lua.LState) int {
			T := c.Tele
			message := L.CheckString(1)
			T.Send(message)
			return 0
		},
		"reply": func(L *lua.LState) int {
			T := c.Tele
			message := L.CheckString(1)
			T.Reply(message)
			return 0
		},
		"http": func(L *lua.LState) int {
			T := c.Tele
			uri := L.CheckString(1)
			resp, err := http.Get(uri)
			if err != nil {
				T.Send("Error occured while http call")
			}
			defer resp.Body.Close()
			body, err := io.ReadAll(resp.Body)
			L.Push(lua.LString(body))
			return 1
		},
	}
	for name, fn := range c.Funcs {
		c.L.SetGlobal(name, c.L.NewFunction(fn))
	}
}

func (LTCtx TgLuaContext) AttachToTGCtx(T tele.Context) {
	LTCtx.Tele = T
	LTCtx.SetupGlobals()
}

func (c *TgLuaContext) ExecLua(src string) error {
	err := c.L.DoString(src)
	return err
}
