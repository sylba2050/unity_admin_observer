package unity_admin_observer

import (
	"encoding/json"
	"io/ioutil"
)

var Config Conf

type Conf struct {
	Mail     string `json:"mail"`
	Password string `json:"password"`
	LoginURL string `json:"login_url"`
	SalesURL string `json:"sales_url"`
}

func init() {
	bytes, err := ioutil.ReadFile("config.json")
	if err != nil {
		panic(err)
	}
	if err := json.Unmarshal(bytes, &Config); err != nil {
		panic(err)
	}
}
