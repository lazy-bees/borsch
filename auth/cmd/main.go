package main

import (
	"github.com/lazy-bees/borsch/auth/config"
	"github.com/lazy-bees/borsch/auth/server"
	"github.com/spf13/viper"
	"log"
)

func main() {
	if err := config.Init(); err != nil {
		log.Fatalf("%s", err.Error())
	}

	s := server.New()

	if err := s.Run(viper.GetString("port")); err != nil {
		log.Fatalf("%s", err.Error())
	}
}
