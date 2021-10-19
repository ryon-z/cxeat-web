package viewcontrollers

import (
	"fmt"
	"net/http"
	"strings"
	"yelloment-api/database"
	envutil "yelloment-api/env_util"
	"yelloment-api/models"
	"yelloment-api/utils"

	"github.com/gin-gonic/gin"
)

// UserCardWithNumSubs : 구독 수와 함께 카드 정보
type UserCardWithNumSubs struct {
	models.UserCard
	NumSubs int `db:"NumSubs"`
}

// GetOwnedUserCardWithNumSubs : 구독 수와 함께 카드 정보 획득
func GetOwnedUserCardWithNumSubs(userNo int) (isSuccess bool, result interface{}) {
	var userCardWithNumSubs []UserCardWithNumSubs
	emptyModelArr := models.GetModelAddr("userCard")
	fieldNamesStringA := "A." + strings.Join(envutil.GetStructFieldNames(emptyModelArr, []string{}), ", A.")
	fieldNamesString := strings.Join(envutil.GetStructFieldNames(emptyModelArr, []string{}), ",")

	sqlQuery := fmt.Sprintf(`
		SELECT %s, count(B.SubsNo) AS "NumSubs" 
		FROM (
			SELECT %s FROM UserCard 
			WHERE (UserNo = %d AND StatusCode = "normal")) AS A
		LEFT JOIN (
			SELECT SubsNo, CardRegNo FROM SubsMst
			WHERE (StatusCode = "normal")) AS B
		ON A.CardRegNo = B.CardRegNo
		GROUP BY A.CardKey
	;`, fieldNamesStringA, fieldNamesString, userNo)

	err := database.DB.Select(&userCardWithNumSubs, sqlQuery)
	if err != nil {
		return false, err.Error()
	}

	return true, userCardWithNumSubs
}

// HandleManageCard : manage-card 뷰
func HandleManageCard(c *gin.Context) {
	state := gin.H{}
	utils.CombineCustomStateGlobalState(c, &state)
	state["cards"] = ""

	userNo := utils.GetUserNo(c)

	isSuccess, result := GetOwnedUserCardWithNumSubs(userNo)
	if !isSuccess {
		slackMsg := fmt.Sprintf("[front]HandleManageCard::GetOwnedUserCardWithNumSubs::%s", result.(string))
		utils.SendSlackMessage(utils.SlackChannel, slackMsg)
		state["errorMessage"] = "시스템 에러 발생"
		state["errorCode"] = "W0009"
		c.HTML(http.StatusBadRequest, "error.html", state)
		return
	}

	var cards = result.([]UserCardWithNumSubs)
	state["cards"] = cards
	state["activeNav"] = "myPage"
	state["navTitle"] = "카드관리"

	c.HTML(http.StatusOK, "manage-card.html", state)
}
