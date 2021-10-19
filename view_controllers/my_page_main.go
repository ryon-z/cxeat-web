package viewcontrollers

import (
	"fmt"
	"net/http"
	"yelloment-api/controllers"
	"yelloment-api/models"
	"yelloment-api/utils"

	"github.com/gin-gonic/gin"
)

// HandleMyPageMain : my-page/main 뷰
func HandleMyPageMain(c *gin.Context) {
	state := gin.H{}
	utils.CombineCustomStateGlobalState(c, &state)
	state["userInfo"] = ""

	isSuccess, result := controllers.GetOwnedUser(c)
	if !isSuccess {
		slackMsg := fmt.Sprintf("[front]HandleMyPageMain::controllers.GetOwnedUser::%s", result.(string))
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
	state["userInfo"] = userMst[0]
	state["activeNav"] = "myPage"
	state["navTitle"] = "마이페이지"

	c.HTML(http.StatusOK, "my-page/main.html", state)
}
