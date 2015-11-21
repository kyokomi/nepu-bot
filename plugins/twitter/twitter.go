package twitter

import (
	"log"
	"strings"

	"github.com/kyokomi/slackbot/plugins"
)

type plugin struct {
	accessToken string
}

func (r plugin) CheckMessage(_ plugins.BotEvent, message string) (bool, string) {
	return strings.Contains(message, "いーすん画像"), message
}

func (r plugin) DoAction(event plugins.BotEvent, message string) bool {
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

func (p *plugin) Help() string {
	return `twitter: いーすん画像表示
	いーすん画像:

		ネプテューヌシリーズの画像をTwitterから検索する。
	`
}

var _ plugins.BotMessagePlugin = (*plugin)(nil)
