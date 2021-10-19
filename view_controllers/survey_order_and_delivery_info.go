package viewcontrollers

import (
	"fmt"
	"net/http"
	"yelloment-api/controllers"
	"yelloment-api/models"
	"yelloment-api/utils"

	"github.com/gin-gonic/gin"
)

// ResultParams : 설문조사 결과에서 넘어온 인자
type ResultParams struct {
	Params string `form:"result-param"`
}

// HandleOrderAndDeliveryInfo : order-and-delivery-info 뷰
func HandleOrderAndDeliveryInfo(c *gin.Context) {
	state := gin.H{}
	utils.CombineCustomStateGlobalState(c, &state)

	state["codeLabels"] = ""
	state["basicDelivery"] = ""
	state["basicCard"] = ""
	state["userAddresses"] = ""
	state["userCards"] = ""
	state["isOrderAndDelivery"] = "yes"

	// 설문조사 결과 인자 획득
	var resultParams ResultParams
	if err := c.ShouldBind(&resultParams); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	state["resultParams"] = resultParams.Params

	// 뷰에서 항목 표시를 위한 코드 라벨
	isSuccess, labelResult := controllers.GetCodesFromCache([]string{"ITEM_CATEGORY", "DELIVERY_PERIOD", "DELIVERY_DOW", "BOX_TYPE", "ORDER_TYPE"})
	if !isSuccess {
		slackMsg := fmt.Sprintf("[front]HandleOrderAndDeliveryInfo::controllers.GetCodesFromCache([]string{'ITEM_CATEGORY', 'DELIVERY_PERIOD', 'DELIVERY_DOW', 'BOX_TYPE', 'ORDER_TYPE'})::%s",
			labelResult.(string))
		utils.SendSlackMessage(utils.SlackChannel, slackMsg)
		state["errorMessage"] = "시스템 에러 발생"
		state["errorCode"] = "W0018"
		c.HTML(http.StatusBadRequest, "error.html", state)
		return
	}
	var codeLabels = *labelResult.(*[]models.CodeMst)
	state["codeLabels"] = codeLabels

	// 기본 주소지
	isSuccess, addressResult := controllers.GetOwnedBasicUserAddress(c)
	if !isSuccess {
		slackMsg := fmt.Sprintf("[front]HandleOrderAndDeliveryInfo::controllers.GetOwnedBasicUserAddress::%s", addressResult.(string))
		utils.SendSlackMessage(utils.SlackChannel, slackMsg)
		state["errorMessage"] = "시스템 에러 발생"
		state["errorCode"] = "W0004"
		c.HTML(http.StatusBadRequest, "error.html", state)
		return
	}
	var basicAddress = *addressResult.(*[]models.UserAddress)
	if len(basicAddress) >= 1 {
		state["basicDelivery"] = basicAddress[0]
	}

	// 기본 카드
	isSuccess, cardResult := controllers.GetOwnedBasicUserCard(c)
	if !isSuccess {
		slackMsg := fmt.Sprintf("[front]HandleOrderAndDeliveryInfo::controllers.GetOwnedBasicUserCard::%s", cardResult.(string))
		utils.SendSlackMessage(utils.SlackChannel, slackMsg)
		state["errorMessage"] = "시스템 에러 발생"
		state["errorCode"] = "W0003"
		c.HTML(http.StatusBadRequest, "error.html", state)
		return
	}
	var userCard = *cardResult.(*[]models.UserCard)
	if len(userCard) >= 1 {
		state["basicCard"] = userCard[0]
	}

	// 유저 소유 주소지 목록
	isSuccess, result := controllers.GetOwnedUserAddress(c)
	if !isSuccess {
		slackMsg := fmt.Sprintf("[front]HandleOrderAndDeliveryInfo::controllers.GetOwnedUserAddress::%s", result.(string))
		utils.SendSlackMessage(utils.SlackChannel, slackMsg)
		state["errorMessage"] = "시스템 에러 발생"
		state["errorCode"] = "W0002"
		c.HTML(http.StatusBadRequest, "error.html", state)
		return
	}
	var userAddresses = *result.(*[]models.UserAddress)
	state["userAddresses"] = userAddresses

	// 유저 소유 카드 목록
	isSuccess, result = controllers.GetOwnedUserCard(c)
	if !isSuccess {
		slackMsg := fmt.Sprintf("[front]HandleOrderAndDeliveryInfo::controllers.GetOwnedUserCard::%s", result.(string))
		utils.SendSlackMessage(utils.SlackChannel, slackMsg)
		state["errorMessage"] = "시스템 에러 발생"
		state["errorCode"] = "W0001"
		c.HTML(http.StatusBadRequest, "error.html", state)
		return
	}

	var userCards = *result.(*[]models.UserCard)
	state["userCards"] = userCards
	state["navTitle"] = "설문조사 결과"

	// 브라우저 캐시 금지
	c.Header("Cache-Control", "no-cache, no-store, must-revalidate")

	c.HTML(http.StatusOK, "survey/order-and-delivery-info.html", state)
}
