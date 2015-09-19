package nepubot

import (
	"math/rand"

	"github.com/kyokomi/slackbot/plugins"
)

// Kaomoji 顔文字
var kaomojiMap = []string{
	"(; ・∀・)",
	"(~_~;)",
	"(-_-;)",
	"?(°_°>)",
	"Σ(￣□￣;)",
	"( ｀・ω・´)",
	"m9( ﾟдﾟ)",
}

type NepubotEvent struct {
	random        *rand.Rand
	messageSender plugins.MessageSender
}

func NewNepubotEvent(rd *rand.Rand, ms plugins.MessageSender) plugins.MessageSender {
	return &NepubotEvent{
		random:        rd,
		messageSender: ms,
	}
}

func (n NepubotEvent) kaomoji() string {
	idx := n.random.Int31n((int32)(len(kaomojiMap) - 1))
	return kaomojiMap[idx]
}

func (n *NepubotEvent) SendMessage(message string, channel string) {
	n.messageSender.SendMessage(message+n.kaomoji(), channel)
}

var _ plugins.MessageSender = (*NepubotEvent)(nil)
