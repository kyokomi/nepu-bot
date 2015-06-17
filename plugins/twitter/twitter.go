package twitter

import (
	"strings"

	"github.com/go-xweb/log"
	"github.com/kyokomi/slackbot/plugins"
	"golang.org/x/net/context"
)

type pluginKey string

var accessToken string

func init() {
	plugins.AddPlugin(pluginKey("twitterImageSearchMessage"), TwitterImageSearchMessage{})
}

type TwitterImageSearchMessage struct {
}

func (r TwitterImageSearchMessage) CheckMessage(ctx context.Context, message string) (bool, string) {
	return strings.Contains(message, "いーすん画像"), message
}

func (r TwitterImageSearchMessage) DoAction(ctx context.Context, message string) bool {
	if accessToken == "" {
		token, err := newAccessToken("", "")
		if err != nil {
			log.Println(err)
			return true
		}
		accessToken = token
	}
	imageURLs, err := searchImages(accessToken, "イストワール", 1)
	if err != nil {
		log.Println(err)
		return true
	}

	for _, imageURL := range imageURLs {
		plugins.SendMessage(ctx, imageURL)
	}

	return false
}

var _ plugins.BotMessagePlugin = (*TwitterImageSearchMessage)(nil)
