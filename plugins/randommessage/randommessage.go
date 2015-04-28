package randommessage

import (
	"math/rand"
	"strings"
	"time"

	"github.com/k0kubun/pp"
	"github.com/kyokomi/nepu-bot/bot"
	"github.com/kyokomi/nepu-bot/plugins"
	"github.com/kyokomi/slack"
)

func init() {
	plugins.Plugins["randomMessage"] = RandomMessage{}
}

var rd = rand.New(rand.NewSource(time.Now().UnixNano()))

type RandomMessage struct {
}

func (r RandomMessage) CheckMessage(ctx bot.BotContext, message string) (bool, string) {
	botUser, _ := ctx.Value("user").(slack.UserDetails)

	if strings.Index(message, botUser.Id) != -1 {
		message = message[strings.Index(message, ":")+len(":"):]
		pp.Println("bot message ", message)
	} else if strings.Index(message, "いーすん") == -1 {
		a := int(rd.Int() % 5)
		if a != 1 {
			return false, ""
		}
	}

	pp.Println(message)
	return true, message
}

func (r RandomMessage) DoAction(ctx bot.BotContext, msEvent *slack.MessageEvent, message string, sendMessageFunc func(message string)) {
	m := bot.NewMessage(msEvent.UserId, msEvent.ChannelId, message)
	sendMessageFunc(bot.CreateResMessage(ctx, m))
}

var _ plugins.BotMessagePlugin = (*RandomMessage)(nil)
