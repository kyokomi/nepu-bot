package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/brettbuddin/victor"
)

func main() {
	bot := victor.New(victor.Config{
		Name:         "victor",
		ChatAdapter:  "shell",
		StoreAdapter: "memory",
		HTTPAddr:     ":8000",
	})

	bot.HandleCommandFunc("hello|hi|howdy", func(s victor.State) {
		s.Chat().Send(s.Message().ChannelID(), fmt.Sprintf("Hello, %s", s.Message().UserName()))
	})

	go bot.Run()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)
	<-sigs

	bot.Stop()
}
