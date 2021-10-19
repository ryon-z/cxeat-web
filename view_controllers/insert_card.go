package viewcontrollers

import (
	"fmt"
	"net/http"
	"yelloment-api/controllers"
	"yelloment-api/models"
	"yelloment-api/utils"

	"github.com/gin-gonic/gin"
)

// HandleInsertCard : insert-card 뷰
func HandleInsertCard(c *gin.Context) {
	state := gin.H{}
	utils.CombineCustomStateGlobalState(c, &state)
	state["basicCard"] = ""

	isSuccess, cardResult := controllers.GetOwnedBasicUserCard(c)
	if !isSuccess {
		slackMsg := fmt.Sprintf("[front]HandleInsertCard::controllers.GetOwnedBasicUserCard::%s", cardResult.(string))
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
	state["activeNav"] = "myPage"
	state["navTitle"] = "카드추가"

	c.HTML(http.StatusOK, "insert-card.html", state)
}
