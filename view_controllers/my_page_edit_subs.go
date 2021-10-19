package viewcontrollers

import (
	"fmt"
	"net/http"
	"yelloment-api/controllers"
	"yelloment-api/models"
	"yelloment-api/utils"

	"github.com/gin-gonic/gin"
)

// HandleEditSubs : subs/edit 뷰
func HandleEditSubs(c *gin.Context) {
	state := gin.H{}
	utils.CombineCustomStateGlobalState(c, &state)
	state["deliveries"] = ""

	isSuccess, itemState := controllers.GetSubsItemsState(4, -12)
	if !isSuccess {
		slackMsg := fmt.Sprintf("[front]HandleEditSubs::controllers.GetSubsItemsState(4, -12)::%s", itemState["error"])
		utils.SendSlackMessage(utils.SlackChannel, slackMsg)
		state["errorMessage"] = "시스템 에러 발생"
		state["errorCode"] = "W0010"
		c.HTML(http.StatusBadRequest, "error.html", state)
		return
	}

	utils.CombineCustomStateGlobalState(c, &state)
	utils.CombineTwoGinH(&state, &itemState)

	// SubsNo 획득
	isSuccess, cookieResult := utils.GetIDFromCookie(c, "subscriptionNo")
	if !isSuccess {
		state["errorMessage"] = "선택한 구독이 없습니다. 구독관리로 이동하여 구독을 선택해주세요."
		c.HTML(http.StatusBadRequest, "error.html", state)
		return
	}
	subsNo := cookieResult.(int)

	// 구독 정보 획득
	isSucces, result := controllers.GetOwnedSubsMstUnit(c, subsNo)
	if !isSucces {
		slackMsg := fmt.Sprintf("[front]HandleEditSubs::controllers.GetOwnedSubsMstUnit::%s", result.(string))
		utils.SendSlackMessage(utils.SlackChannel, slackMsg)
		state["errorMessage"] = "시스템 에러 발생"
		state["errorCode"] = "W0011"
		c.HTML(http.StatusBadRequest, "error.html", state)
		return
	}

	subsInfos := *result.(*[]models.SubsMst)
	if len(subsInfos) < 1 {
		state["errorMessage"] = "잘못된 구독 정보 요청"
		c.HTML(http.StatusBadRequest, "error.html", state)
		return
	}
	subsInfo := subsInfos[0]
	state["subsInfo"] = subsInfo

	timeString := subsInfo.FirstDate
	isSuccess, korDow := utils.GetEngDow(timeString)
	if !isSuccess {
		slackMsg := fmt.Sprintf("[front]HandleEditSubs::utils.GetEngDow(timeString)::%s", korDow)
		utils.SendSlackMessage(utils.SlackChannel, slackMsg)
		state["errorMessage"] = "시스템 에러 발생"
		state["errorCode"] = "W0012"
		c.HTML(http.StatusBadRequest, "error.html", state)
		return
	}
	state["subsDayOfWeek"] = korDow
	state["activeNav"] = "myPage"
	state["navTitle"] = "상품변경"

	c.HTML(http.StatusOK, "my-page/edit-subs.html", state)
	return
}
