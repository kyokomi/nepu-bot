package nepubot

import (
	"math/rand"
	"strings"
	"time"

	"github.com/kyokomi/slackbot/plugins"
	"github.com/kyokomi/slackbot/slackctx"
	"golang.org/x/net/context"
)

const (
	randomMessageCount = 5 // その他は5回に1回だけ返信する
)

type pluginKey string

func init() {
	plugins.AddPlugin(pluginKey("nepubot"), NepuMessage{})
}

var rd = rand.New(rand.NewSource(time.Now().UnixNano()))

type NepuMessage struct {
}

func (r NepuMessage) CheckMessage(ctx context.Context, message string) (bool, string) {
	api := slackctx.FromSlackClient(ctx)
	botName := api.Name

	if strings.Index(message, botName) != -1 {
		// MessageにBotのIDが含まれる
		message = message[strings.Index(message, ":")+len(":"):]
	} else if strings.Index(message, botName) != -1 {
		// Messageに名前がふくまれている
		text := strings.Replace(message, botName, "", 1)
		if len(text) == 0 {
			text = "hello" // 名前だけの場合は固定で挨拶
		}
		message = text
	} else {
		a := int(rd.Int() % randomMessageCount)
		if a != 1 {
			return false, ""
		}
	}

	return true, message
}

func (r NepuMessage) DoAction(ctx context.Context, message string) bool {
	msEvent := slackctx.FromMessageEvent(ctx)

	if strings.Index(message, "静かに") != -1 {
		plugins.Stop()
		go func() {
			// 5分黙ってもらう
			time.Sleep(5 * time.Minute)
			plugins.Start()
		}()
	}

	m := NewMessage(msEvent.BotID, msEvent.Channel, message)
	plugins.SendMessage(ctx, DocomoAPIMessage(ctx, m))

	return false // stop not next
}

var _ plugins.BotMessagePlugin = (*NepuMessage)(nil)
