package unity_admin_observer

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type Message struct {
	Channel string `json:"channel"`
	Text    string `json:"text"`
}

func SendSlackMessage(packages []string, nowSales []int) {
	text := buildSlackNotificationText(packages, nowSales)
	message := Message{Channel: Config.SlackChannelName, Text: text}
	rawSendJSON, err := json.Marshal(message)
	if err != nil {
		panic(err)
	}
	req, err := http.NewRequest("POST", Config.SlackURL, bytes.NewBuffer(rawSendJSON))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", Config.SlackToken))
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
}

func buildSlackNotificationText(packages []string, nowSales []int) string {
	if len(packages) != len(nowSales) {
		panic(errors.New("len(packages) != len(nowSales)"))
	}
	var text string
	for _, user := range Config.SlackUserIDs {
		text += fmt.Sprintf("<@%s> ", user)
	}
	for i := 0; i < len(packages); i++ {
		text += "\n"
		text += fmt.Sprintf("%sが購入されました。現在の累計販売数は%d個です。", packages[i], nowSales[i])
	}

	return text
}
