package webapp

import (
	"encoding/json"
	"net/http"
	"strings"

	"log"
	"math/rand"
	"net/url"
	"os"

	"github.com/kyokomi/go-docomo/docomo"
	"github.com/zenazn/goji/web"
)

var logger = log.New(os.Stdout, "nepu-bot", log.Llongfile)
var random = rand.New(rand.NewSource(1))

// Kaomoji 顔文字
var Kaomoji = []string{
	"(; ・∀・)",
	"(~_~;)",
	"(-_-;)",
	"?(°_°>)",
	"Σ(￣□￣;)",
	"( ｀・ω・´)",
	"m9( ﾟдﾟ)",
}

// SlackClient is Slack IncomingURL Client.
type SlackClient struct {
	Name             string // TODO: 旧API
	SlackIncomingURL string // TODO: 旧API
	Token            string
}

// Message is Slack Receive Message.
type Message struct {
	userID, userName, channelID, channelName, text string
}

// OutgoingMessage is Slack PostRequest Message.
type OutgoingMessage struct {
	Channel  string `json:"channel"`
	Username string `json:"username"`
	Text     string `json:"text"`
}

// HubotSlackWebhook hubotとしてSlackのWebhockを待ち受けるHandleFunc
func HubotSlackWebhook(c web.C, _ http.ResponseWriter, r *http.Request) {

	slackClient := c.Env["slack"].(*SlackClient)
	docomoClient := c.Env["docomo"].(*docomo.Client)

	m := Message{
		userID:      r.PostFormValue("user_id"),
		userName:    r.PostFormValue("user_name"),
		channelID:   r.PostFormValue("channel_id"),
		channelName: r.PostFormValue("channel_name"),
		text:        r.PostFormValue("text"),
	}
	var resMessage string

	// 名前のみの場合は固定文言に置き換え
	text := strings.Replace(m.text, slackClient.Name, "", 1)
	if len(text) == 0 {
		text = "hello"
	}

	switch {
	default:
		// その他は全部雑談

		// 雑談API呼び出し
		dq := docomo.DialogueRequest{}
		dq.Nickname = &m.userName
		dq.Utt = &text
		res, err := docomoClient.Dialogue.Get(dq, true)
		if err != nil {
			logger.Println(err)
			return
		}

		resMessage = res.Utt

	case containsArray(text, "おしえて", "教えて"):
		// 知識Q&A
		qa := docomo.KnowledgeQARequest{}
		qa.QAText = text
		for _, word := range []string{"おしえて", "教えて"} {
			qa.QAText = strings.Replace(qa.QAText, word, "", -1)
		}
		res, err := docomoClient.KnowledgeQA.Get(qa)
		if err != nil {
			logger.Println(err)
			return
		}

		if res.Success() {
			resMessage = res.Answers[0].AnswerText
		} else {
			resMessage = "はて?"
		}
	}

	// 顔文字をランダムで付与する
	idx := random.Int31n((int32)(len(Kaomoji) - 1))
	message := resMessage + " " + Kaomoji[idx]
	// 結果を非同期でSlackへ
	go slackClient.send(m.channelID, message)
}

func (s SlackClient) send(channelID, msg string) {
	ms := &OutgoingMessage{
		Channel:  channelID,
		Username: s.Name,
		Text:     msg,
	}
	body, err := json.Marshal(ms)
	if err != nil {
		logger.Println("error sending to json marshal:", err)
	}

	if _, err := http.PostForm(s.SlackIncomingURL, url.Values{"payload": {string(body)}}); err != nil {
		logger.Println("error sending to chat:", err)
	}
}

func containsArray(s string, substrs ...string) bool {
	for i := 0; i < len(substrs); i++ {
		if strings.Index(s, substrs[i]) >= 0 {
			return true
		}
	}
	return false
}
