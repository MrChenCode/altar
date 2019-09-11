package test

import "testing"

func TestTopicInfo(t *testing.T) {
	RunTestApi(t, &Request{
		Path:   "/topicinfo",
		Method: Get,
		RespResult: []string{
			"add_time", "desc", "image", "platform_type", "wxqq_title",
			"similarity", "title", "type", "weibo_desc", "wxqq_desc",
		},
	})
}
