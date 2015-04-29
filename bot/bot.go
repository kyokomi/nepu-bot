package bot

import (
	"github.com/kyokomi/go-docomo/docomo"
	"github.com/nlopes/slack"
	"golang.org/x/net/context"
)

type DocomoClient struct {
	*docomo.Client
}

// SlackClient is Slack IncomingURL Client.
type SlackClient struct {
	*slack.Slack
	Name  string
	Token string
}

type key string

const (
	slackClientKey key = "SlackClient"
	slackRTMKey    key = "SlackRTM"
)

func NewSlackClient(ctx context.Context, name string, token string) context.Context {
	c := SlackClient{}
	c.Slack = slack.New(token)
	c.Name = name
	c.Token = token
	return context.WithValue(ctx, slackClientKey, c)
}

func FromSlackClient(ctx context.Context) SlackClient {
	return ctx.Value(slackClientKey).(SlackClient)
}

func NewSlackRTM(ctx context.Context, protocol, origin string) context.Context {
	api := FromSlackClient(ctx)
	wsAPI, err := api.StartRTM(protocol, origin)
	if err != nil {
		return ctx
	}
	return context.WithValue(ctx, slackRTMKey, wsAPI)
}

func FromSlackRTM(ctx context.Context) *slack.SlackWS {
	return ctx.Value(slackRTMKey).(*slack.SlackWS)
}
