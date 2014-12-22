package webapp

import (
	"net/http"
	"strings"
	"encoding/json"

	"github.com/zenazn/goji/web"
	"log"
	"os"
	"math/rand"
	"net/url"
	"github.com/kyokomi/nepu-bot/src/config"

	docomo "github.com/kyokomi/go-docomo"
)

var sendURL = os.Getenv("SLACK_INCOMING_URL")

var logger = log.New(os.Stdout, "nepu-bot", log.Llongfile)
var random = rand.New(rand.NewSource(1))

var Kaomoji = []string{
	"(; ・∀・)",
	"(~_~;)",
	"(-_-;)",
	"?(°_°>)",
	"Σ(￣□￣;)",
	"( ｀・ω・´)",
	"m9( ﾟдﾟ)",
}

type Message struct {
	userID, userName, channelID, channelName, text string
}

type OutgoingMessage struct {
	Channel  string `json:"channel"`
	Username string `json:"username"`
	Text     string `json:"text"`
}

func HubotSlackWebhook(c web.C, _ http.ResponseWriter, r *http.Request) {

	bot := c.Env["bot"].(*config.BotConfig)
	d := c.Env["docomo"].(*docomo.DocomoClient)

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

	var resMessage string
	if containsArray(m.text, []string{"おしえて", "教えて"}) {
		// 知識Q&A
		qa := docomo.QARequest{}
		qa.QAText = m.text
		for _, word := range []string{"おしえて", "教えて"} {
			qa.QAText = strings.Replace(qa.QAText, word, "", -1)
		}
		res, err := d.SendQA(&qa)
		if err != nil {
			logger.Println(err)
			return
		}
		resMessage = res.Answers[0].AnswerText

	} else {
		// その他は全部雑談

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
		resMessage = resMap["utt"]
	}

	// 顔文字をランダムで付与する
	idx := random.Int31n((int32)(len(Kaomoji) - 1))
	message := resMessage + " " + Kaomoji[idx]
	// 結果を非同期でSlackへ
	go Send(bot.Name, m.channelID, message)
}

func Send(botName, channelID, msg string) {
	body, err := json.Marshal(&OutgoingMessage{
		Channel:  channelID,
		Username: botName,
		Text:     msg,
	})
	if err != nil {
		logger.Println("error sending to json marshal:", err)
	}

	if _, err := http.PostForm(sendURL, url.Values{"payload": {string(body)}}); err != nil {
		logger.Println("error sending to chat:", err)
	}
}

//

func containsArray(s string, substrs []string) bool {
	for _, substr := range substrs {
		if strings.Index(s, substr) >= 0 {
			return true
		}
	}
	return false
}
