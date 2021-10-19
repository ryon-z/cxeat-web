package viewcontrollers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
	"yelloment-api/config"
	envutil "yelloment-api/env_util"
	"yelloment-api/utils"

	"github.com/gin-gonic/gin"
)

// AccessTokenResp : accessToken response
type AccessTokenResp struct {
	Code    int       `json:"code"`
	Expire  time.Time `json:"expire"`
	Token   string    `json:"token"`
	Message string    `json:"message"`
}

// HandleLogin : handel login
func HandleLogin(c *gin.Context, baseURL string) {
	code := c.Query("code")
	var token string = ""
	var expire string = ""

	var state = gin.H{}
	utils.CombineCustomStateGlobalState(c, &state)

	if code != "" {
		origin := envutil.GetGoDotEnvVariable("PROTOCOL") + "://" + c.Request.Host
		baseURL = fmt.Sprintf("%s%s?code=%s", origin, baseURL, code)
		isSuccess, response := utils.CustomRequest(baseURL, "GET", [][]string{}, [][]string{}, []*http.Cookie{})
		if !isSuccess {
			slackMsg := fmt.Sprintf("[front]HandleLogin::CustomRequest(%s)::%s", baseURL, response)
			utils.SendSlackMessage(utils.SlackChannel, slackMsg)
			state["errorCode"] = 1022
			state["errorMessage"] = config.ExposedErrorMessageMap[1]
			c.HTML(http.StatusBadRequest, "error.html", state)
			return
		}

		fmt.Println("login response", response)

		var accessTokenResp AccessTokenResp
		if err := json.Unmarshal([]byte(response), &accessTokenResp); err != nil {
			fmt.Println("unmarshal 실패", err.Error())
			slackMsg := fmt.Sprintf("[front]HandleLogin::json.Unmarshal([]byte(%s), &accessTokenResp)::%s", response, err.Error())
			utils.SendSlackMessage(utils.SlackChannel, slackMsg)
			state["errorCode"] = 9999
			state["errorMessage"] = "시스템 오류 발생"
			c.HTML(http.StatusBadRequest, "error.html", state)
			return
		}
		if accessTokenResp.Code != 200 {
			fmt.Println(accessTokenResp.Message)
			_, errProto := config.UnstringifyError(accessTokenResp.Message)
			var index int
			if errProto.YelloCode == 9004 || errProto.YelloCode == 9014 || errProto.YelloCode == 9002 ||
				errProto.YelloCode == 5001 {
				index = errProto.YelloCode
			} else {
				index = int(errProto.YelloCode / 1000)
			}
			fmt.Println("errProto.Message", errProto.Message)
			fmt.Println("errProto.YelloCode", errProto.YelloCode)
			fmt.Println("index", index)
			state["errorCode"] = errProto.YelloCode
			state["errorMessage"] = config.ExposedErrorMessageMap[index]
			state["activeRedirectLogout"] = "yes"
			if errProto.YelloCode == 5001 {
				newSignupDateString := errProto.Message
				state["errorMessage"] = strings.ReplaceAll(config.ExposedErrorMessageMap[index], "newSignupDate", newSignupDateString)
			}
			c.HTML(http.StatusBadRequest, "error.html", state)
			return
		}

		token = accessTokenResp.Token
		expire = accessTokenResp.Expire.String()

		// Temporary, 개발 중일 경우에만 httpOnly False
		expireSecString := envutil.GetGoDotEnvVariable("AUTH_TOKEN_EXPIRE_SEC")
		expireSec, err := strconv.Atoi(expireSecString)
		if err != nil {
			state["errorCode"] = 9998
			state["errorMessage"] = "시스템 에러 발생"
			slackMsg := fmt.Sprintf("[front]HandleLogin::strconv.Atoi(expireSecString)::%s", err.Error())
			utils.SendSlackMessage(utils.SlackChannel, slackMsg)
			c.HTML(http.StatusBadRequest, "error.html", state)
			return
		}

		c.SetCookie("jwt", token, expireSec, "/", "", false, true)
		c.SetCookie("expire", expire, expireSec, "/", "", false, false)

		c.Redirect(http.StatusTemporaryRedirect, "/web/signup")
	}

	state["navTitle"] = "로그인"
	c.HTML(http.StatusOK, "login.html", state)

	return
}

// HandleKakaoLogin : Login redirect Controller
func HandleKakaoLogin(c *gin.Context) {
	HandleLogin(c, "/user/kakao")
}

// HandleNaverLogin : Login redirect Controller
func HandleNaverLogin(c *gin.Context) {
	HandleLogin(c, "/user/naver")
}
