package viewcontrollers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"yelloment-api/controllers"
	"yelloment-api/database"
	envutil "yelloment-api/env_util"
	"yelloment-api/models"
	"yelloment-api/utils"

	"github.com/gin-gonic/gin"
)

// ExposedSubs : 화면에 보여주기 위한 구독 목록
type ExposedSubs struct {
	models.SubsMst
	BoxLabel       string
	PeriodLabel    string
	DayOfWeek      string
	SubsStatus     string
	CateTypeLabels []string
}

// GetExposedSubs : 화면에 보여주기 위한 구독 목록 획득
func GetExposedSubs(c *gin.Context) (isSuccess bool, result interface{}) {
	userNo := utils.GetUserNo(c)

	var exposedSubs []ExposedSubs
	emptyModelArr := models.GetModelAddr("subsMst")
	fieldNamesString := "A." + strings.Join(envutil.GetStructFieldNames(emptyModelArr, []string{}), ", A.")

	sqlQuery := fmt.Sprintf(`
		SELECT %s
		FROM SubsMst AS A
		WHERE (A.UserNo = %d)
		ORDER BY A.StatusCode DESC, A.RegDate DESC
	;`, fieldNamesString, userNo)

	err := database.DB.Select(&exposedSubs, sqlQuery)
	if err != nil {
		return false, err.Error()
	}

	return true, exposedSubs
}

// HandleMyPageSubsList : my-page/subs/list 뷰
func HandleMyPageSubsList(c *gin.Context) {
	state := gin.H{}
	utils.CombineCustomStateGlobalState(c, &state)
	state["subsList"] = ""

	isSuccess, result := GetExposedSubs(c)
	if !isSuccess {
		slackMsg := fmt.Sprintf("[front]HandleMyPageSubsList::GetExposedSubs::%s", result.(string))
		utils.SendSlackMessage(utils.SlackChannel, slackMsg)
		state["errorMessage"] = "시스템 에러 발생"
		state["errorCode"] = "W0023"
		c.HTML(http.StatusBadRequest, "error.html", state)
		return
	}

	var subsList = result.([]ExposedSubs)
	state["subsList"] = subsList

	for i, subs := range subsList {
		// 한글 요일로 변환
		timeString := subs.FirstDate
		isSuccess, result := utils.GetKorDow(timeString)
		if !isSuccess {
			slackMsg := fmt.Sprintf("[front]HandleMyPageSubsList::utils.GetKorDow(%s)::%s", timeString, result)
			utils.SendSlackMessage(utils.SlackChannel, slackMsg)
			state["errorMessage"] = "시스템 에러 발생"
			state["errorCode"] = "W0022"
			c.HTML(http.StatusBadRequest, "error.html", state)
			return
		}
		subsList[i].DayOfWeek = result

		// 카테고리 타입 문자열
		isSuccess, cateTypeLabels := controllers.MakeCateTypeLabels(subsList[i].CateType, "|")
		if !isSuccess {
			slackMsg := fmt.Sprintf("[front]HandleMyPageSubsList::controllers.MakeCateTypeLabels(%s, '|')::%s",
				subsList[i].CateType, cateTypeLabels)
			utils.SendSlackMessage(utils.SlackChannel, slackMsg)
			state["errorMessage"] = "시스템 에러 발생"
			state["errorCode"] = "W0015"
			c.HTML(http.StatusBadRequest, "error.html", state)
			return
		}
		for _, row := range cateTypeLabels.([]string) {
			subsList[i].CateTypeLabels = append(subsList[i].CateTypeLabels, row)
		}

		// 박스 라벨
		isSuccess, boxLabel := controllers.GetCodeLabelFromCache("BOX_TYPE", subsList[i].BoxType)
		if !isSuccess {
			slackMsg := fmt.Sprintf("[front]HandleMyPageSubsList::controllers.GetCodeLabelFromCache('BOX_TYPE', %s)::%s",
				subsList[i].BoxType, boxLabel)
			utils.SendSlackMessage(utils.SlackChannel, slackMsg)
			state["errorMessage"] = "시스템 에러 발생"
			state["errorCode"] = "W0016"
			c.HTML(http.StatusBadRequest, "error.html", state)
			return
		}
		subsList[i].BoxLabel = boxLabel

		// 배송 주기
		isSuccess, periodLabel := controllers.GetCodeLabelFromCache("DELIVERY_PERIOD", strconv.Itoa(subsList[i].PeriodDay))
		if !isSuccess {
			slackMsg := fmt.Sprintf("[front]HandleMyPageSubsList::controllers.GetCodeLabelFromCache('DELIVERY_PERIOD', %d)::%s",
				subsList[i].PeriodDay, periodLabel)
			utils.SendSlackMessage(utils.SlackChannel, slackMsg)
			state["errorMessage"] = "시스템 에러 발생"
			state["errorCode"] = "W0016"
			c.HTML(http.StatusBadRequest, "error.html", state)
			return
		}
		subsList[i].PeriodLabel = periodLabel
	}
	state["activeNav"] = "myPage"
	state["navTitle"] = "구독 정보"

	c.HTML(http.StatusOK, "my-page/subs-list.html", state)
}
