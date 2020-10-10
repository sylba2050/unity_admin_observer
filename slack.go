package unity_admin_observer

type Message struct {
	Channel string `json:"channel"`
	Text    string `json:"text"`
}
