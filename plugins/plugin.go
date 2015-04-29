package plugins

import (
	"github.com/nlopes/slack"
	"golang.org/x/net/context"
)

var (
	ctx     context.Context
	Plugins = map[interface{}]BotMessagePlugin{}
)

func init() {
	ctx = context.Background()
}

type BotMessagePlugin interface {
	CheckMessage(ctx context.Context, message string) (bool, string)
	DoAction(ctx context.Context, msEvent *slack.MessageEvent, message string, sendMessageFunc func(message string))
}

func Context() context.Context {
	return ctx
}
