package plugins

import (
	"github.com/nlopes/slack"
	"golang.org/x/net/context"
)

var Plugins = map[interface{}]BotMessagePlugin{}

type BotMessagePlugin interface {
	CheckMessage(ctx context.Context, message string) (bool, string)
	DoAction(ctx context.Context, msEvent *slack.MessageEvent, message string, sendMessageFunc func(message string))
}
