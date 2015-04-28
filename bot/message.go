package bot

import (
	"strings"

	"log"
	"os"

	"github.com/k0kubun/pp"
	"github.com/kyokomi/go-docomo/docomo"
)

var logger = log.New(os.Stdout, "nepu-bot", log.Llongfile)

// Message is Slack Receive Message.
type Message struct {
	userID, userName, channelID, channelName, text string
}

func NewMessage(userID, channelID, text string) Message {
	var m Message
	m.userID = userID
	m.channelID = channelID
	m.text = text
	return m
}

func CreateResMessage(ctx BotContext, m Message) string {
	var resMessage string

	pp.Println(m)

	// 名前のみの場合は固定文言に置き換え
	text := strings.Replace(m.text, ctx.Slack.Name, "", 1)
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
		res, err := ctx.Docomo.Dialogue.Get(dq, true)
		if err != nil {
			logger.Println(err)
			return m.text
		}

		pp.Println(res)

		resMessage = res.Utt

	case containsArray(text, "おしえて", "教えて"):
		// 知識Q&A
		qa := docomo.KnowledgeQARequest{}
		qa.QAText = text
		for _, word := range []string{"おしえて", "教えて"} {
			qa.QAText = strings.Replace(qa.QAText, word, "", -1)
		}
		res, err := ctx.Docomo.KnowledgeQA.Get(qa)
		if err != nil {
			logger.Println(err)
			return m.text
		}

		if res.Success() {
			resMessage = res.Answers[0].AnswerText
		} else {
			resMessage = "はて?"
		}
	}

	pp.Println(resMessage + " " + GetKaomji())

	return resMessage + " " + GetKaomji()
}

func containsArray(s string, substrs ...string) bool {
	for i := 0; i < len(substrs); i++ {
		if strings.Index(s, substrs[i]) >= 0 {
			return true
		}
	}
	return false
}
