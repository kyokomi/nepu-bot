package main

import (
	"fmt"
	"os"
	"os/signal"

	"encoding/json"
	"log"

	"github.com/brettbuddin/victor"
	"github.com/kyokomi/nepu-bot/docomogo"
)

func main() {
	bot := victor.New(victor.Config{
		Name:         "いーすん",
		ChatAdapter:  "slack",
		StoreAdapter: "memory",
		HTTPAddr:     ":8000",
	})

	d := docomogo.NewDialogue()

	bot.HandleCommandFunc("hello|hi|howdy", func(s victor.State) {
		s.Chat().Send(s.Message().ChannelID(), fmt.Sprintf("Hello, %s", s.Message().UserName()))
	})
	bot.HandleCommandFunc(".*", func(s victor.State) {

		res := d.Send(s.Message().Text())

		var resMap map[string]string
		if err := json.Unmarshal(res, &resMap); err != nil {
			log.Fatalln("Unmarshal ", err)
		}

		s.Chat().Send(s.Message().ChannelID(), resMap["utt"])
	})

	go bot.Run()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)
	<-sigs

	bot.Stop()
}
