package common

import (
	"github.com/zhenghaoz/gorse/client"
)

var Gorse *client.GorseClient

func GorseInit() {
	gorse := client.NewGorseClient("http://127.0.0.1:8087", "api_key")
	Gorse = gorse
}

func GetGorse() *client.GorseClient {
	return Gorse
}
