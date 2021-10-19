package viewcontrollers

import (
	"fmt"
	"net/http"
	"strings"
	"yelloment-api/controllers"
	"yelloment-api/models"
	"yelloment-api/utils"

	"github.com/gin-gonic/gin"
)

// HandleMyPageUserInfo : my-page/user/info 뷰
func HandleMyPageUserInfo(c *gin.Context) {
	state := gin.H{}
	utils.CombineCustomStateGlobalState(c, &state)
	state["userInfo"] = ""

	isSuccess, result := controllers.GetOwnedUser(c)
	if !isSuccess {
		slackMsg := fmt.Sprintf("[front]HandleMyPageUserInfo::controllers.GetOwnedUser::%s", result.(string))
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

	// 생년월일 문자열로 변경.
	fmt.Println("userMst[0]", userMst[0])
	fmt.Println("userMst[0].BirthDay", userMst[0].BirthDay)
	userInfo := userMst[0]
	if userInfo.BirthDay != nil {
		splited := strings.Split(*userInfo.BirthDay, "T")
		birthDay := strings.ReplaceAll(splited[0], "-", "")
		userInfo.BirthDay = &birthDay
	}
	state["userInfo"] = userInfo
	state["activeNav"] = "myPage"
	state["navTitle"] = "정보수정"

	c.HTML(http.StatusOK, "my-page/user-info.html", state)
}
