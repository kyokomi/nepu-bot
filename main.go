package main

import (
	"os"

	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	docomo "github.com/kyokomi/go-docomo"
	"github.com/zenazn/goji"
	"math/rand"
)

var Kaomoji = []string{
	"(; ・∀・)",
	"(~_~;)",
	"(-_-;)",
	"?(°_°>)",
	"Σ(￣□￣;)",
	"( ｀・ω・´)",
	"m9( ﾟдﾟ)",
}

type Config struct {
	Name,
	ChatAdapter,
	StoreAdapter,
	HTTPAddr string
}

type Message struct {
	userID, userName, channelID, channelName, text string
}

type OutgoingMessage struct {
	Channel  string `json:"channel"`
	Username string `json:"username"`
	Text     string `json:"text"`
}

var logger = log.New(os.Stderr, "nepu-bot", log.Llongfile)

var sendURL = os.Getenv("SLACK_INCOMING_URL")

var random = rand.New(rand.NewSource(1))


func main() {

	bot := Config{
		Name:         "いーすん",
		ChatAdapter:  "slack",
		StoreAdapter: "memory",
		HTTPAddr:     os.Getenv("PORT"),
	}
	team := os.Getenv("VICTOR_SLACK_TEAM")
	token := os.Getenv("VICTOR_SLACK_TOKEN")
	d := docomo.New(os.Getenv("DOCOMO_APIKEY"))

	goji.Post("/hubot/slack-webhook", func(_ http.ResponseWriter, r *http.Request) {
		m := Message{
			userID:      r.PostFormValue("user_id"),
			userName:    r.PostFormValue("user_name"),
			channelID:   r.PostFormValue("channel_id"),
			channelName: r.PostFormValue("channel_name"),
			text:        r.PostFormValue("text"),
		}

		if !strings.Contains(m.text, bot.Name) {
			return
		}

		// 名前のみの場合は固定文言に置き換え
		t := strings.Replace(m.text, bot.Name, "", 1)
		if len(t) == 0 {
			t = "hello"
		}

		// 雑談API呼び出し
		res, err := d.SendZatsudan(m.userName, t)
		if err != nil {
			logger.Println(err)
			return
		}

		var resMap map[string]string
		if err := json.Unmarshal(res, &resMap); err != nil {
			logger.Println("Unmarshal ", err)
			return
		}

		// 顔文字をランダムで付与する
		idx := random.Int31n((int32)(len(Kaomoji) - 1))
		message := resMap["utt"] + Kaomoji[idx]
		// 結果を非同期でSlackへ
		go Send(bot.Name, team, token, m.channelID, message)
	})
	goji.Serve()
}

func Send(botName, team, token, channelID, msg string) {
	body, err := json.Marshal(&OutgoingMessage{
		Channel:  channelID,
		Username: botName,
		Text:     msg,
	})
	if err != nil {
		logger.Println("error sending to json marshal:", err)
	}

	//endpoint := fmt.Sprintf("https://%s.slack.com/services/hooks/hubot?token=%s", team, token)
	endpoint := fmt.Sprintf("%s", sendURL)

	if _, err := http.PostForm(endpoint, url.Values{"payload": {string(body)}}); err != nil {
		logger.Println("error sending to chat:", err)
	}
}
