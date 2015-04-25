package bot

import (
	"golang.org/x/net/context"
	"github.com/kyokomi/go-docomo/docomo"
)

type BotContext struct {
	context.Context
	Slack  *SlackClient
	Docomo *docomo.Client
}

// SlackClient is Slack IncomingURL Client.
type SlackClient struct {
	Name             string // TODO: 旧API
	SlackIncomingURL string // TODO: 旧API
	Token            string
}
