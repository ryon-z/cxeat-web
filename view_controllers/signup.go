package viewcontrollers

import (
	"fmt"
	"net/http"
	"strings"
	"yelloment-api/controllers"
	"yelloment-api/models"
	"yelloment-api/utils"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// HandleSignup : 회원가입 처리
func HandleSignup(c *gin.Context) {
	state := gin.H{}
	utils.CombineCustomStateGlobalState(c, &state)

	claim := jwt.ExtractClaims(c)
	// 최초 로그인이 아닌 경우
	if claim["isNew"].(bool) == false {
		pagePath, err := c.Cookie("pageAfterLogin")
		if pagePath == "/web/login" || err != nil {
			pagePath = "/web/main"
		}
		c.Redirect(http.StatusTemporaryRedirect, pagePath)
		return
	}

	// 회원가입
	isSuccess, result := controllers.GetOwnedUser(c)
	if !isSuccess {
		slackMsg := fmt.Sprintf("[front]HandleSignup::controllers.GetOwnedUser::%s", result.(string))
		utils.SendSlackMessage(utils.SlackChannel, slackMsg)
		state["errorMessage"] = "시스템 에러 발생"
		state["errorCode"] = "W0013"
		c.HTML(http.StatusBadRequest, "error.html", state)
		return
	}

	var userMst = *result.(*[]models.UserMst)
	if len(userMst) < 1 {
		state["errorMessage"] = "잘못된 유저 정보 요청"
		c.HTML(http.StatusBadRequest, "error.html", state)
		return
	}

	userMst[0].UserName = strings.ReplaceAll(userMst[0].UserName, "(emoji)", "")

	state["userInfo"] = userMst[0]
	state["navTitle"] = "추가 회원 정보"
	state["isSignup"] = "yes"

	c.HTML(http.StatusOK, "signup.html", state)
	return
}
