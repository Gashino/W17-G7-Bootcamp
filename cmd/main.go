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
	}
	app := server.NewServerChi(cfg)
	// - run
	if err := app.Run(); err != nil {
		fmt.Println(err)
		return
	}
}
