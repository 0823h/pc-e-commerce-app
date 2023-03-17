package common

import (
	"context"

	"github.com/zhenghaoz/gorse/client"
)

var Gorse *client.GorseClient

func GorseInit() {
	gorse := client.NewGorseClient("http://127.0.0.1:8087", "api_key")
	Gorse = gorse
	Gorse.InsertItem(context.Background(), client.Item{ItemId: "1", IsHidden: false,
		Categories: []string{"Shoes"}, Timestamp: "2023/03/18 12:22", Labels: []string{"Shoes labels"}, Comment: ""})
}

func GetGorse() *client.GorseClient {
	return Gorse
}
