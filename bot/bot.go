package bot

import (
	"github.com/kyokomi/go-docomo/docomo"
	"golang.org/x/net/context"
)

type BotContext struct {
	context.Context
	Slack  *SlackClient
	Docomo *docomo.Client
}

// SlackClient is Slack IncomingURL Client.
type SlackClient struct {
	Name  string
	Token string
}
