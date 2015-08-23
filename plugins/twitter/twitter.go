package twitter

import (
	"strings"

	"github.com/go-xweb/log"
	"github.com/kyokomi/slackbot/plugins"
)

type Plugin struct {
	accessToken string
}

func (r Plugin) CheckMessage(_ plugins.BotEvent, message string) (bool, string) {
	return strings.Contains(message, "いーすん画像"), message
}

func (r Plugin) DoAction(event plugins.BotEvent, message string) bool {
	if r.accessToken == "" {
		token, err := newAccessToken("", "")
		if err != nil {
			log.Println(err)
			return true
		}
		r.accessToken = token
	}
	imageURLs, err := searchImages(r.accessToken, "イストワール", 1)
	if err != nil {
		log.Println(err)
		return true
	}

	for _, imageURL := range imageURLs {
		event.Reply(imageURL)
	}

	return false
}

var _ plugins.BotMessagePlugin = (*Plugin)(nil)
