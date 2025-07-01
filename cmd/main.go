package main

import (
	"app/cmd/server"
	"fmt"
)

func main() {
	// env
	// ...

	// app
	// - config
	cfg := &server.ConfigServerChi{
		ServerAddress:  ":8080",
		LoaderFilePath: "docs/db/sections.json",
	}
	app := server.NewServerChi(cfg)
	// - run
	if err := app.Run(); err != nil {
		fmt.Println(err)
		return
	}
}
