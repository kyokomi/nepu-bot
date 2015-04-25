package main

import (
	"os"

	"github.com/kyokomi/nepu-bot/src/webapp"

	"flag"
	"net/http"

	"github.com/kyokomi/go-docomo/docomo"
	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"
	"golang.org/x/net/context"
	"fmt"
	"github.com/zenazn/goji/graceful"
)

type BotContext struct {
	context.Context
	Slack  *webapp.SlackClient
	Docomo *docomo.Client
}

func main() {
	var apikey string
	flag.StringVar(&apikey, "d", os.Getenv("DOCOMO_APIKEY"), "ドコモのAPIKEY")
	var slackURL string
	flag.StringVar(&slackURL, "s", os.Getenv("SLACK_INCOMING_URL"), "SlackのIncomingのURL")
	var token string
	flag.StringVar(&token, "token", os.Getenv("SLACK_BOT_TOKEN"), "SlackのBotToken")
	flag.Parse()

	ctx := BotContext{}
	ctx.Context = context.Background()

	slackClient := &webapp.SlackClient{
		Name:             "いーすん",
		Token:            token,
		SlackIncomingURL: slackURL,
	}
	ctx.Slack = slackClient

	docomoClient := docomo.NewClient(apikey)
	ctx.Docomo = docomoClient

	if slackURL != "" {
		fmt.Println("start incoming bot ...")
		webIncoming(slackClient, docomoClient)
	} else if token != "" {
		fmt.Println("start rtm bot ...")
		webSocket(ctx)
	} else {
		fmt.Println("not run...")
	}
}

func webIncoming(slackClient *webapp.SlackClient, docomoClient *docomo.Client) {
	goji.Use(func(c *web.C, h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			c.Env["slack"] = slackClient
			c.Env["docomo"] = docomoClient
			h.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	})

	goji.Post("/hubot/slack-webhook", webapp.HubotSlackWebhook)
	graceful.PostHook(func() { fmt.Println("end") }
	goji.Serve()
}

func webSocket(ctx BotContext) {
	// TODO:
}
