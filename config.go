package unity_admin_observer

import (
	"encoding/json"
	"io/ioutil"
)

var Config Conf

type Conf struct {
	Mail             string   `json:"mail"`
	Password         string   `json:"password"`
	LoginURL         string   `json:"login_url"`
	SalesURL         string   `json:"sales_url"`
	SlackURL         string   `json:"slack_url"`
	SlackToken       string   `json:"slack_token"`
	SlackChannelName string   `json:"slack_channel_name"`
	SlackUserIDs     []string `json:"slack_user_ids"`
}

func init() {
	bytes, err := ioutil.ReadFile("/home/siruba_2050/unity_admin_observer/config.json")
	if err != nil {
		panic(err)
	}
	if err := json.Unmarshal(bytes, &Config); err != nil {
		panic(err)
	}
}
