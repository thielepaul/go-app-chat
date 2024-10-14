package frontend

import (
	"log"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"github.com/thielepaul/go-app-chat/rpc"
)

type Chat struct {
	app.Compo
	rpc      *rpc.Frontend
	inputVal string
}

func (c *Chat) Render() app.UI {
	messages := []app.UI{}
	if c.rpc != nil {
		for i := len(c.rpc.Messages) - 1; i >= 0; i-- {
			messages = append(messages, app.Li().Class("p-2 border-b border-gray-200 dark:border-gray-700").Style("scroll-snap-align", "start").Text(c.rpc.Messages[i]))
		}
	}
	return app.Div().Class("container mx-auto p-4 dark:bg-gray-900 dark:text-white h-screen flex flex-col").Body(
		app.Div().Class("flex-shrink-0").Body(
			app.H1().Class("text-3xl font-bold mb-4").Text("Simple Chat!"),
			app.P().Class("mb-4").Text("This is an anonymous chat:"),
		),
		app.Ul().ID("messages-container").Class("flex-grow flex flex-col-reverse overflow-auto mb-4").Body(messages...).Style("scroll-snap-type", "y mandatory"),
		app.Div().Class("flex-shrink-0 w-full bg-white dark:bg-gray-900 p-4 flex").Body(
			app.Input().Class("flex-grow p-2 border border-gray-300 rounded mr-2 dark:border-gray-700 dark:bg-gray-800 dark:text-white").
				Type("text").Value(c.inputVal).OnInput(c.ValueTo(&c.inputVal)).
				OnKeyPress(func(ctx app.Context, e app.Event) {
					if e.Get("key").String() == "Enter" {
						c.OnClick(ctx, e)
					}
				}),
			app.Button().Class("p-2 bg-blue-500 text-white rounded dark:bg-blue-700").Text("Send").OnClick(c.OnClick),
		),
	)
}

func (c *Chat) OnClick(ctx app.Context, e app.Event) {
	if c.rpc.Client == nil {
		log.Println("rpc client is nil")
		return
	}
	if _, err := rpc.Call(c.rpc.Client, (&rpc.Backend{}).AddMessage, c.inputVal); err != nil {
		log.Println("addmessage error: ", err)
	}
	c.inputVal = ""
}

func (c *Chat) OnMount(ctx app.Context) {
	var err error
	c.rpc, err = rpc.NewFrontend(ctx, app.Window().URL())
	if err != nil {
		log.Panicf("error creating rpc client: %v", err)
	}
}
