package main

import (
	"log"
	"net/http"
	"os"

	"github.com/thielepaul/go-app-chat/frontend"
	"github.com/thielepaul/go-app-chat/rpc"

	"github.com/NYTimes/gziphandler"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

func main() {
	app.Route("/", app.NewZeroComponentFactory(&frontend.Chat{}))
	app.RunWhenOnBrowser()

	if !app.IsServer {
		// early return helps to reduce code size of wasm binary due to dead code elimination
		return
	}

	http.Handle("/", gziphandler.GzipHandler(&app.Handler{
		Name:         "Chat",
		Description:  "A simple chat",
		Scripts:      []string{"https://cdn.tailwindcss.com"},
		Body:         func() app.HTMLBody { return app.Body().Class("dark:bg-gray-900 dark:text-white") },
		Resources:    embeddedResourceResolver{Handler: http.FileServer(http.FS(web))},
		LoadingLabel: "Loading Chat...",
	}))

	_ = rpc.NewBackend()

	port := ":8000"
	if os.Getenv("PORT") != "" {
		port = ":" + os.Getenv("PORT")
	}

	log.Printf("Listening on " + port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}
