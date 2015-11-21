package nepubot

import (
	"math/rand"
	"strings"
	"time"

	"github.com/kyokomi/slackbot/plugins"

	godocomo "github.com/kyokomi/go-docomo/docomo"
	"github.com/kyokomi/slackbot"
	"github.com/kyokomi/slackbot/plugins/docomo"
)

const (
	randomMessageCount = 5 // その他は5回に1回だけ返信する
)

var rd = rand.New(rand.NewSource(time.Now().UnixNano()))

type plugin struct {
	docomo  plugins.BotMessagePlugin
	Plugins plugins.PluginManager
}

func NewPlugin(pm plugins.PluginManager, docomoClient *godocomo.Client, repository slackbot.Repository) plugins.BotMessagePlugin {
	return &plugin{
		docomo:  docomo.NewPlugin(docomoClient, repository),
		Plugins: pm,
	}
}

func (r *plugin) CheckMessage(event plugins.BotEvent, message string) (bool, string) {
	if ok, message := r.docomo.CheckMessage(event, message); ok {
		return true, message
	}

	a := int(rd.Int() % randomMessageCount)
	if a == 1 {
		return true, message
	}

	return false, ""
}

func (r *plugin) DoAction(event plugins.BotEvent, message string) bool {
	if strings.Index(message, "静かに") != -1 {
		r.Plugins.StopReply()
		go func() {
			// 5分黙ってもらう
			time.Sleep(5 * time.Minute)
			r.Plugins.StartReply()
		}()
	}
	return r.docomo.DoAction(plugins.NewBotEvent(NewNepubotEvent(rd, event.GetMessageSender()),
		event.BotID(),
		event.BotName(),
		event.SenderID(),
		event.SenderName(),
		event.BaseText(),
		event.Channel(),
	), message)
}

func (r *plugin) Help() string {
	return `nepu-bot: ネプテューヌbot
	超次元ゲイムネプテューヌのいーすんことイストワールのbot

	パッシブスキル:

		5回に1回くらいの割合で会話に混ざる

	アクティブスキル:

		静かに: 5分間Slackへの投稿を控える
	`
}

var _ plugins.BotMessagePlugin = (*plugin)(nil)
