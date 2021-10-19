package viewcontrollers

import (
	"fmt"
	"net/http"
	"yelloment-api/controllers"
	"yelloment-api/models"
	"yelloment-api/utils"

	"github.com/gin-gonic/gin"
)

// HandleManageDelivery : manage-delivery 뷰
func HandleManageDelivery(c *gin.Context) {
	state := gin.H{}
	utils.CombineCustomStateGlobalState(c, &state)
	state["deliveries"] = ""

	isSuccess, result := controllers.GetOwnedUserAddress(c)
	if !isSuccess {
		slackMsg := fmt.Sprintf("[front]HandleManageDelivery::controllers.GetOwnedUserAddress::%s", result.(string))
		utils.SendSlackMessage(utils.SlackChannel, slackMsg)
		state["errorMessage"] = "시스템 에러 발생"
		state["errorCode"] = "W0002"
		c.HTML(http.StatusBadRequest, "error.html", state)
		return
	}

	var basicAddress = *result.(*[]models.UserAddress)
	state["deliveries"] = basicAddress
	state["activeNav"] = "myPage"
	state["navTitle"] = "배송지 선택"

	c.HTML(http.StatusOK, "manage-delivery.html", state)
}
