package controllers

import (
	"fmt"
	"net/http"
	"yelloment-api/models"
	"yelloment-api/utils"

	"github.com/gin-gonic/gin"
)

// AlarmCancelSubs : /my-page/subs/cancel 뷰
func AlarmCancelSubs(c *gin.Context) {
	state := gin.H{}
	utils.CombineCustomStateGlobalState(c, &state)

	// 유저명과 유저 핸드폰 번호 획득
	userName := ""
	userPhone := ""
	isSuccess, userResult := GetOwnedUser(c)
	if !isSuccess {
		slackMsg := fmt.Sprintf("[front]AlarmCancelSubs::controllers.GetOwnedUser::%s", userResult.(string))
		utils.SendSlackMessage(utils.SlackChannel, slackMsg)
		c.JSON(http.StatusBadRequest, gin.H{"error": userResult.(string)})
		return
	}
	var user = *userResult.(*[]models.UserMst)
	if len(user) > 0 {
		userName = user[0].UserName
		userPhone = user[0].UserPhone
	}

	// 알림톡 전송
	utils.SendAlrimTalk(
		userPhone,
		"bizp_2021083018402810216575017",

		fmt.Sprintf(`%s님의 큐잇 정기구독 해지가 완료되었습니다. 

신선한 농산물과 건강한 식생활이 필요할 땐, 언제든 다시 큐잇을 이용하실 수 있어요.

도움이 필요하시면 고객센터에 문의해 주세요.

고객센터 070-4166-6077

다시 만나뵐 날을 기다릴게요.
큐잇 드림`, userName),
		[]utils.AlrimTalkButton{},
	)

	c.JSON(http.StatusOK, gin.H{"data": "success"})
	return
}
