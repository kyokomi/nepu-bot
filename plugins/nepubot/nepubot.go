package nepubot

import (
	"math/rand"
	"strings"
	"time"

	"github.com/kyokomi/slackbot"
	"github.com/kyokomi/slackbot/plugins"
	"github.com/nlopes/slack"
	"golang.org/x/net/context"
)

type pluginKey string

func init() {
	plugins.AddPlugin(pluginKey("nepubot"), NepuMessage{})
}

var rd = rand.New(rand.NewSource(time.Now().UnixNano()))

type NepuMessage struct {
}

func (r NepuMessage) CheckMessage(ctx context.Context, message string) (bool, string) {
	api := slackbot.FromSlackClient(ctx)
	botUser := api.GetInfo().User
	botName := api.Name

	if strings.Index(message, botUser.Id) != -1 {
		// MessageにBotのIDが含まれる
		message = message[strings.Index(message, ":")+len(":"):]
	} else if strings.Index(message, botName) != -1 {
		// Messageに名前がふくまれている
		text := strings.Replace(message, botName, "", 1)
		if len(text) == 0 {
			text = "hello" // 名前だけの場合は固定で挨拶
		}
	} else {
		// その他は5回に1回だけ返信する
		a := int(rd.Int() % 5)
		if a != 1 {
			return false, ""
		}
	}

	return true, message
}

func (r NepuMessage) DoAction(ctx context.Context, msEvent *slack.MessageEvent, message string, sendMessageFunc func(message string)) bool {
	m := NewMessage(msEvent.UserId, msEvent.ChannelId, message)
	sendMessageFunc(DocomoAPIMessage(ctx, m))
	return false // stop not next
}

var _ plugins.BotMessagePlugin = (*NepuMessage)(nil)
