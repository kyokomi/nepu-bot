package main

import (
	"os"

	"github.com/kyokomi/nepu-bot/src/webapp"

	docomo "github.com/kyokomi/go-docomo"
	"github.com/zenazn/goji"
	"github.com/kyokomi/nepu-bot/src/config"
	"github.com/zenazn/goji/web"
	"net/http"
)

func main() {

	botConfig := &config.BotConfig{
		Name:         "いーすん",
		ChatAdapter:  "slack",
		StoreAdapter: "memory",
		HTTPAddr:     os.Getenv("PORT"),
	}
	docomoClient := docomo.New(os.Getenv("DOCOMO_APIKEY"))

	goji.Use(func(c *web.C, h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			c.Env["bot"] = botConfig
			c.Env["docomo"] = docomoClient
			h.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	})

	goji.Post("/hubot/slack-webhook", webapp.HubotSlackWebhook)
	goji.Serve()
}

