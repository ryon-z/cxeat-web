package utils

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// CustomRequest : 커스텀 요청
func CustomRequest(baseURL string, method string, data [][]string, headers [][]string, cookies []*http.Cookie) (isSuccess bool, result string) {
	// Check method is valid
	switch strings.ToLower(method) {
	case "get":
	case "post":
	case "patch":
	case "delete":
	default:
		return false, fmt.Sprintf("Method is not allowed, method: %s\n", method)
	}

	// Get Content Type
	var contentType string
	for _, value := range headers {
		if value[0] == "Content-Type" && value[1] == "application/json" {
			contentType = "application/json"
		}
	}

	// Set Data Values
	var dataPhrase string
	var dataContainer []string
	for _, value := range data {
		if len(value) != 2 {
			return false, "length of data length is not 2"
		}
		dataKey := value[0]
		dataValue := value[1]
		var row string
		switch contentType {
		case "application/json":
			row = fmt.Sprintf("'%s':'%s'", dataKey, dataValue)
		default:
			row = fmt.Sprintf("%s=%s", dataKey, dataValue)
		}
		dataContainer = append(dataContainer, row)
	}
	switch contentType {
	case "application/json":
		dataPhrase = "{" + strings.Join(dataContainer, ",") + "}"
	default:
		dataPhrase = strings.Join(dataContainer, "&")
	}

	// Create request
	client := &http.Client{}
	req, err := http.NewRequest(strings.ToUpper(method), baseURL, strings.NewReader(dataPhrase))
	if err != nil {
		return false, err.Error()
	}
	if len(cookies) > 0 {
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}
	}

	// Set Headers
	for _, value := range headers {
		if len(value) != 2 {
			return false, "length of header length is not 2"
		}
		headerKey := value[0]
		headerValue := value[1]
		req.Header.Add(headerKey, headerValue)
	}

	// Get response
	resp, err := client.Do(req)
	if err != nil {
		return false, err.Error()
	}

	// Stringify response.Body
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err.Error()
	}

	return true, string(body)
}
