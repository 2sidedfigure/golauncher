package main

import (
	"log"
	"net/http"
	"strings"

	"try/thunder"

	"code.google.com/p/go.net/websocket"
)

func Listen(address string, launcher thunder.Launcher) error {
	mux := http.NewServeMux()

	mux.Handle("/control", websocket.Handler(func(ws *websocket.Conn) {
		launcher.LedOn()
		launcher.Stop()

		for {
			var (
				input string
				quit  bool
			)

			if err := websocket.Message.Receive(ws, &input); err != nil {
				log.Printf("ERROR: %s", err)
				break
			}

			switch strings.ToUpper(input) {
			case "DOWN":
				launcher.Down()
			case "UP":
				launcher.Up()
			case "LEFT":
				launcher.Left()
			case "RIGHT":
				launcher.Right()
			case "FIRE":
				launcher.Fire()
			case "QUIT":
				quit = true
			default:
				launcher.Stop()
			}

			if quit {
				break
			}
		}

		launcher.Stop()
		launcher.LedOff()
	}))
	mux.Handle("/", http.FileServer(FS(false)))

	return http.ListenAndServe(address, mux)
}
