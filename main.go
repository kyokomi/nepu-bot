package main

import (
	"fmt"
	"os"
	"os/signal"

	"encoding/json"
	"log"

	"github.com/brettbuddin/victor"
	docomo "github.com/kyokomi/go-docomo"
	"github.com/k0kubun/pp"
	"strings"
)

var logger = log.New(os.Stderr, "nepu-bot", log.Llongfile)

func main() {
	bot := victor.New(victor.Config{
		Name:         "いーすん",
		ChatAdapter:  "slack",
		StoreAdapter: "memory",
		HTTPAddr:     ":8000",
	})

	d := docomo.New()

	bot.HandleCommandFunc("hello|hi|howdy", func(s victor.State) {
		s.Chat().Send(s.Message().ChannelID(), fmt.Sprintf("Hello, %s", s.Message().UserName()))
	})
	bot.HandleCommandFunc("image .*", func(s victor.State) {
		pp.Println(s.Message())


		// image 以降を取得する
		d.SendImage(strings.TrimRight(s.Message().Text(), "image "))

	})
	bot.HandleCommandFunc(".*", func(s victor.State) {
		pp.Println(s.Message())

		res, err := d.SendZatsudan(s.Message().UserName(), s.Message().Text())
		if err != nil {
			logger.Println(err)
			return
		}

		var resMap map[string]string
		if err := json.Unmarshal(res, &resMap); err != nil {
			logger.Println("Unmarshal ", err)
			return
		}

		// Send Slack
		s.Chat().Send(s.Message().ChannelID(), resMap["utt"])
	})

	go bot.Run()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)
	<-sigs

	bot.Stop()
}
