package nepubot

import (
	"math/rand"
	"strings"
	"time"

	"github.com/kyokomi/go-docomo/docomo"
	"github.com/kyokomi/slackbot/plugins"
)

const (
	randomMessageCount = 5 // その他は5回に1回だけ返信する
)

var rd = rand.New(rand.NewSource(time.Now().UnixNano()))

type Plugin struct {
	Docomo  *docomo.Client
	Plugins plugins.PluginManager // TODO: ちょと無理矢理...
}

func (r *Plugin) CheckMessage(event plugins.BotEvent, message string) (bool, string) {
	botID := event.BotID()

	if strings.Index(message, botID) != -1 {
		// MessageにBotのIDが含まれる
		message = message[strings.Index(message, ":")+len(":"):]
	} else if strings.Index(message, botID) != -1 {
		// Messageに名前がふくまれている
		text := strings.Replace(message, botID, "", 1)
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

func (r *Plugin) DoAction(event plugins.BotEvent, message string) bool {
	if strings.Index(message, "静かに") != -1 {
		r.Plugins.StopReply()
		go func() {
			// 5分黙ってもらう
			time.Sleep(5 * time.Minute)
			r.Plugins.StartReply()
		}()
	}

	event.Reply(r.DocomoAPIMessage(event.SenderName(), message))

	return false // stop not next
}

var _ plugins.BotMessagePlugin = (*Plugin)(nil)
