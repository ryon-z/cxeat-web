package utils

import (
	"fmt"
	"strings"
)

// GetURL : GET Method 호출 시 사용할 URL을 조합
func GetURL(baseURL string, params []string) string {
	return fmt.Sprintf("%s?%s", baseURL, strings.Join(params, "&"))
}
