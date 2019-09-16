package test

import "testing"

func TestTopicInfo(t *testing.T) {
	RunTestApi(t, &Request{
		Path:   "/topicinfo",
		Method: Get,
		RespResult: []string{
			"tid", "title", "desc", "type", "image", "icon", "platform_type",
			"wxqq_title", "wxqq_desc", "weibo_desc", "add_time", "list",
		},
	})
}
