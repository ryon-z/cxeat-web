package viewcontrollers

import (
	"fmt"
	"net/http"
	"yelloment-api/controllers"
	"yelloment-api/middleware/mixin"
	"yelloment-api/models"
	"yelloment-api/utils"

	"github.com/gin-gonic/gin"
)

// HandleInsertDelivery : insert-delivery 뷰
func HandleInsertDelivery(c *gin.Context) {
	state := gin.H{}
	utils.CombineCustomStateGlobalState(c, &state)
	state["basicDelivery"] = ""

	isSuccess, result := controllers.GetOwnedBasicUserAddress(c)
	if !isSuccess {
		slackMsg := fmt.Sprintf("[front]HandleInsertDelivery::controllers.GetOwnedBasicUserAddress::%s", result.(string))
		utils.SendSlackMessage(utils.SlackChannel, slackMsg)
		state["errorMessage"] = "시스템 에러 발생"
		state["errorCode"] = "W0004"
		c.HTML(http.StatusBadRequest, "error.html", state)
		return
	}

	var basicAddress = *result.(*[]models.UserAddress)
	if len(basicAddress) >= 1 {
		state["basicDelivery"] = basicAddress[0]
	}
	state["activeNav"] = "myPage"
	state["navTitle"] = "배송지 추가"

	c.HTML(http.StatusOK, "insert-delivery.html", state)
}

// GetSubsMstUsingAddress : 입력 받은 주소가 사용되고 있는 구독 리스트 리턴
func GetSubsMstUsingAddress(userNo int, addressNo int) (isSuccess bool, result interface{}) {
	wherePhrase := fmt.Sprintf(`UserNo = %d AND AddressNo = %d AND (StatusCode = "normal" OR StatusCode = "pause")`, userNo, addressNo)
	return mixin.List("subsMst", wherePhrase, "")
}

// HandleModifyDelivery : modify-delivery 뷰
func HandleModifyDelivery(c *gin.Context) {
	state := gin.H{}
	utils.CombineCustomStateGlobalState(c, &state)
	state["addressInfo"] = ""
	state["usedAddressOtherSubs"] = false

	// addressNo 획득
	isSuccess, cookieResult := utils.GetIDFromCookie(c, "addressNo")
	if !isSuccess {
		state["errorMessage"] = "선택된 배송지가 없습니다. 배송지관리로 이동하여 배송지를 선택해주세요."
		c.HTML(http.StatusBadRequest, "error.html", state)
		return
	}
	addressNo := cookieResult.(int)

	isSuccess, result := controllers.GetOwnedUserAddressUnit(c, addressNo)
	if !isSuccess {
		slackMsg := fmt.Sprintf("[front]HandleInsertDelivery::controllers.GetOwnedUserAddressUnit::%s", result.(string))
		utils.SendSlackMessage(utils.SlackChannel, slackMsg)
		state["errorMessage"] = "시스템 에러 발생"
		state["errorCode"] = "W0005"
		c.HTML(http.StatusBadRequest, "error.html", state)
		return
	}

	var addressInfo = *result.(*[]models.UserAddress)
	if len(addressInfo) != 1 {
		state["errorMessage"] = "시스템 에러 발생"
		state["errorCode"] = "W0006"
		c.HTML(http.StatusBadRequest, "error.html", state)
		return
	}
	state["addressInfo"] = addressInfo[0]

	isSuccess, result = GetSubsMstUsingAddress(addressInfo[0].UserNo, addressInfo[0].AddressNo)
	if !isSuccess {
		slackMsg := fmt.Sprintf("[front]HandleInsertDelivery::GetSubsMstUsingAddress(%d, %d)::%s",
			addressInfo[0].UserNo, addressInfo[0].AddressNo, result.(string))
		utils.SendSlackMessage(utils.SlackChannel, slackMsg)
		state["errorMessage"] = "시스템 에러 발생"
		state["errorCode"] = "W0007"
		c.HTML(http.StatusBadRequest, "error.html", state)
		return
	}
	var subsMstUsingAddresses = *result.(*[]models.SubsMst)
	if len(subsMstUsingAddresses) >= 1 {
		state["usedAddressOtherSubs"] = true
	}

	state["activeNav"] = "myPage"
	state["navTitle"] = "배송정보"

	c.HTML(http.StatusOK, "modify-delivery.html", state)
}
