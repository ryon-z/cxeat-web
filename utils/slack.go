package utils

import (
	"fmt"
	"net/http"
	"strings"
)

var SlackChannel string

func SendSlackMessage(channel string, msg string) (isSuccess bool) {
	url := "https://hooks.slack.com/services/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"
	data := [][]string{}

	switch channel {
	case "system":
		url = "https://hooks.slack.com/services/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"
	case "operation":
		url = "https://hooks.slack.com/services/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"
	case "debug":
		url = "https://hooks.slack.com/services/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"
	case "internal":
		url = "https://hooks.slack.com/services/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"
		data = append(data, []string{"type", "mrkdwn"})
	default:
		url = "https://hooks.slack.com/services/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"
	}
	msg = strings.ReplaceAll(msg, "'", "")
	msg = strings.ReplaceAll(msg, "\"", "")
	data = append(data, []string{"text", msg})

	isSuccess, err := CustomRequest(
		url,
		"POST",
		data,
		[][]string{{"Content-Type", "application/json"}},
		[]*http.Cookie{},
	)

	fmt.Println("isSuccess", isSuccess)
	fmt.Println("err", err)

	return isSuccess
}
