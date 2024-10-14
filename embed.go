package main

import (
	"embed"
	"net/http"
)

//go:generate env GOARCH=wasm GOOS=js go build -trimpath -o web/app.wasm

//go:embed web
var web embed.FS

type embeddedResourceResolver struct {
	http.Handler
}

func (r embeddedResourceResolver) Resolve(location string) string {
	if location == "" {
		return "/"
	}
	return location
}
