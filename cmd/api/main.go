package main

import "github.com/Kaa-dan/webrtc-websocket-server.git/internal/config"

func main() {

	cfg := config.MustLoad()

	//logging config (excluding sensitive values)
	if cfg.LogLevel == "debug" || cfg.Environment == "development" {
		cfg.PrintConfig()
	}

}
