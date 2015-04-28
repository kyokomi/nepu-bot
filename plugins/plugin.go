package plugins

import (
	"github.com/kyokomi/nepu-bot/bot"
	"github.com/kyokomi/slack"
)

var Plugins = map[interface{}]BotMessagePlugin{}

type BotMessagePlugin interface {
	CheckMessage(ctx bot.BotContext, message string) (bool, string)
	DoAction(ctx bot.BotContext, msEvent *slack.MessageEvent, message string, sendMessageFunc func(message string))
}
