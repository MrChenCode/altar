package test

import (
	"testing"
)

func TestBookInfo(t *testing.T) {
	RunTestApi(t, &Request{
		Path:   "/gameinfo",
		Method: Get,
		RespResult: map[string]interface{}{
			"gameid":   float64(100),
			"gamename": "三国战纪",
		},
	})
}
