package main

import (
	"os"

	"github.com/kyokomi/nepu-bot/src/webapp"

	docomo "github.com/kyokomi/go-docomo"
	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"
	"net/http"
	"flag"
)

func main() {

	var apikey string
	flag.StringVar(&apikey, "d", os.Getenv("DOCOMO_APIKEY"), "ドコモのAPIKEY")
	var slackURL string
	flag.StringVar(&slackURL, "s", os.Getenv("SLACK_INCOMING_URL"), "SlackのIncomingのURL")
	flag.Parse()

	slackClient := &webapp.SlackClient{
		Name:         "いーすん",
		SlackIncomingURL: slackURL,
	}

	docomoClient := docomo.New(apikey)

	goji.Use(func(c *web.C, h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			c.Env["slack"] = slackClient
			c.Env["docomo"] = docomoClient
			h.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	})

	goji.Post("/hubot/slack-webhook", webapp.HubotSlackWebhook)

	goji.Serve()
}

