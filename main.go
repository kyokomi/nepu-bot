package main

import (
	"os"

	"github.com/kyokomi/nepu-bot/src/webapp"

	"flag"
	"net/http"

	"fmt"
	"time"

	"strings"

	"github.com/k0kubun/pp"
	"github.com/kyokomi/go-docomo/docomo"
	"github.com/kyokomi/nepu-bot/bot"
	"github.com/kyokomi/slack"
	"github.com/zenazn/goji"
	"github.com/zenazn/goji/graceful"
	"github.com/zenazn/goji/web"
	"golang.org/x/net/context"
	"math/rand"
)

func main() {
	var apikey string
	flag.StringVar(&apikey, "d", os.Getenv("DOCOMO_APIKEY"), "ドコモのAPIKEY")
	var slackURL string
	flag.StringVar(&slackURL, "s", os.Getenv("SLACK_INCOMING_URL"), "SlackのIncomingのURL")
	var token string
	flag.StringVar(&token, "token", os.Getenv("SLACK_BOT_TOKEN"), "SlackのBotToken")
	flag.Parse()

	ctx := bot.BotContext{}
	ctx.Context = context.Background()

	slackClient := &bot.SlackClient{
		Name:             "いーすん",
		Token:            token,
		SlackIncomingURL: slackURL,
	}
	ctx.Slack = slackClient

	docomoClient := docomo.NewClient(apikey)
	ctx.Docomo = docomoClient

	if slackURL != "" {
		fmt.Println("start incoming bot ...")
		webIncoming(ctx)
	} else if token != "" {
		fmt.Println("start rtm bot ...")
		webSocket(ctx)
	} else {
		fmt.Println("not run...")
	}
}

func webIncoming(bot bot.BotContext) {
	goji.Use(func(c *web.C, h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			c.Env["context"] = bot
			h.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	})

	goji.Post("/hubot/slack-webhook", webapp.HubotSlackWebhook)
	graceful.PostHook(func() { fmt.Println("end") })
	goji.Serve()
}

func webSocket(ctx bot.BotContext) {
	chSender := make(chan slack.OutgoingMessage)
	chReceiver := make(chan slack.SlackEvent)

	api := slack.New(ctx.Slack.Token)
	api.SetDebug(true)
	wsAPI, err := api.StartRTM("", "http://example.com")
	if err != nil {
		fmt.Errorf("%s\n", err)
	}

	ctx.Context = context.WithValue(ctx.Context, "user", wsAPI.GetInfo().User)

	go wsAPI.HandleIncomingEvents(chReceiver)
	go wsAPI.Keepalive(20 * time.Second)
	go func(wsAPI *slack.SlackWS, chSender chan slack.OutgoingMessage) {
		for {
			select {
			case msg := <-chSender:
				wsAPI.SendMessage(&msg)
			}
		}
	}(wsAPI, chSender)

	for {
		select {
		case msg := <-chReceiver:
			fmt.Print("Event Received: ")
			switch msg.Data.(type) {
			case slack.HelloEvent:
			// TODO: デフォルトChannelに何か投げたい
			case *slack.MessageEvent:
				a := msg.Data.(*slack.MessageEvent)
				message := MessageResponse(ctx, a)
				if !a.IsStarred && message != "" {
					chSender <- *wsAPI.NewOutgoingMessage(message, a.ChannelId)
				}
			case *slack.PresenceChangeEvent:
				a := msg.Data.(*slack.PresenceChangeEvent)
				fmt.Printf("Presence Change: %v\n", a)
			case slack.LatencyReport:
				a := msg.Data.(slack.LatencyReport)
				fmt.Printf("Current latency: %v\n", a.Value)
			case *slack.SlackWSError:
				error := msg.Data.(*slack.SlackWSError)
				fmt.Printf("Error: %d - %s\n", error.Code, error.Msg)
			default:
				fmt.Printf("Unexpected: %v\n", msg.Data)
			}
		}
	}
}
var rd = rand.New(rand.NewSource(time.Now().UnixNano()))

func MessageResponse(ctx bot.BotContext, msEvent *slack.MessageEvent) string {
	user, _ := ctx.Value("user").(slack.UserDetails)
	if user.Id == msEvent.UserId {
		// 自分のやつはスルーする
		return ""
	}

	messageText := msEvent.Text

	// TODO: ここでキーワードでハンドリングとか
	if strings.Index(messageText, user.Id) != -1 {
		messageText = messageText[strings.Index(messageText, ":")+len(":"):]
		pp.Println("bot message ", messageText)
	} else if strings.Index(messageText, "いーすん") == -1 {
		a := int(rd.Int() % 5)
		fmt.Println("################## ", a)
		if a != 1 {
			return ""
		}
	} else {
		fmt.Println("################## ", "else")
	}

	// TODO: メッセージ生成（以前のやつ
	m := webapp.NewMessage(msEvent.UserId, msEvent.ChannelId, messageText)
	return webapp.CreateResMessage(ctx, m)
}
