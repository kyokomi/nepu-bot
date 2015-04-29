package nepubot

import (
	"strings"

	"github.com/kyokomi/go-docomo/docomo"
	"golang.org/x/net/context"
)

func NewContext(ctx context.Context, apiKey string) context.Context {
	return docomo.NewContext(ctx, apiKey)
}

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

func DocomoAPIMessage(ctx context.Context, m Message) string {
	d := docomo.FromContext(ctx)

	text := m.text
	var resMessage string
	switch {
	default:
		// その他は全部雑談

		// 雑談API呼び出し
		dq := docomo.DialogueRequest{}
		dq.Nickname = &m.userName
		dq.Utt = &text
		res, err := d.Dialogue.Get(dq, true)
		if err != nil {
			return m.text
		}

		resMessage = res.Utt

	case containsArray(text, "おしえて", "教えて"):
		// 知識Q&A
		qa := docomo.KnowledgeQARequest{}
		qa.QAText = text
		for _, word := range []string{"おしえて", "教えて"} {
			qa.QAText = strings.Replace(qa.QAText, word, "", -1)
		}
		res, err := d.KnowledgeQA.Get(qa)
		if err != nil {
			return m.text
		}

		if res.Success() {
			resMessage = res.Answers[0].AnswerText
		} else {
			resMessage = "はて?"
		}
	}

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
