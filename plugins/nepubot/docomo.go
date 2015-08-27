package nepubot

import (
	"strings"

	"github.com/kyokomi/go-docomo/docomo"
)

func (n *Plugin) DocomoAPIMessage(userName, text string) string {
	var resMessage string
	switch {
	default:
		// その他は全部雑談

		// 雑談API呼び出し
		dq := docomo.DialogueRequest{}
		dq.Nickname = &userName
		dq.Utt = &text
		res, err := n.Docomo.Dialogue.Get(dq, true)
		if err != nil {
			return text
		}

		resMessage = res.Utt

	case containsArray(text, "おしえて", "教えて"):
		// 知識Q&A
		qa := docomo.KnowledgeQARequest{}
		qa.QAText = text
		for _, word := range []string{"おしえて", "教えて"} {
			qa.QAText = strings.Replace(qa.QAText, word, "", -1)
		}
		res, err := n.Docomo.KnowledgeQA.Get(qa)
		if err != nil {
			return text
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
