package global

import (
	envutil "yelloment-api/env_util"

	"github.com/gin-gonic/gin"
)

// CssRandomVersion : 버전 업 시 css 캐싱 방지를 위한 랜덤 버전
var CssRandomVersion string

// globalState : 전역 변수
var globalState = gin.H{
	"loggedIn":             "no",
	"activeNav":            "no",
	"navTitle":             "큐잇",
	"titleSuffix":          " - 큐잇",
	"needNavBack":          "yes",
	"newLineSymbol":        envutil.GetGoDotEnvVariable("NEW_LINE_SYMBOL"),
	"cssRandomVersion":     "0",
	"activeRedirectLogout": "no",
	"kakaoClientID":        envutil.GetGoDotEnvVariable("KAKAO_CLIENT_ID"),
	"naverClientID":        envutil.GetGoDotEnvVariable("NAVER_CLIENT_ID"),
}

// GlobalStateKeys : 전역 변수 키 값
var GlobalStateKeys = []string{}

// InitGlobalState : Web view에서 사용하는 전역 state 세팅
func InitGlobalState(c *gin.Context) {
	for key, value := range globalState {
		if key == "cssRandomVersion" {
			value = CssRandomVersion
		}
		c.Set(key, value)
		GlobalStateKeys = append(GlobalStateKeys, key)
	}
}

// UpdateGlobalState : 전역 WebState 업데이트
func UpdateGlobalState(c *gin.Context) {
	// 최초 실행 시 초기화
	_, exists := c.Get("loggedIn")
	if !exists {
		InitGlobalState(c)
	}

	// Login 전역변수 갱신
	c.Set("loggedIn", "no")

	_, jwtErr := c.Cookie("jwt")
	_, expireErr := c.Cookie("expire")
	if jwtErr == nil && expireErr == nil {
		c.Set("loggedIn", "yes")
	}
}
