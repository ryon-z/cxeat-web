package viewcontrollers

import (
	"fmt"
	"net/http"
	"strings"
	"yelloment-api/controllers"
	"yelloment-api/database"
	"yelloment-api/utils"

	"github.com/gin-gonic/gin"
)

// AlrimTalkInfo : 알림톡정보
type AlrimTalkInfo struct {
	UserName            string `db:"UserName"`
	UserPhone           string `db:"UserPhone"`
	PeriodDay           int    `db:"PeriodDay"`
	CateType            string `db:"CateType"`
	BoxType             string `db:"BoxType"`
	DesiredDeliveryDate string `db:"DesiredDeliveryDate"`
	RcvName             string `db:"RcvName"`
	MainAddress         string `db:"MainAddress"`
	SubAddress          string `db:"SubAddress"`
	SubsPrice           *int   `db:"SubsPrice"`
	OrderPrice          *int   `db:"OrderPrice"`
	CateTypeLabels      []string
}

// GetAlrimTalkInfoSubs : 구독 알림톡정보 획득
func GetAlrimTalkInfoSubs(userNo int, subsNo int) (isSuccess bool, result interface{}) {
	var alrimTalkInfo []AlrimTalkInfo
	sqlQuery := fmt.Sprintf(`
		SELECT A.UserName, A.UserPhone, B.PeriodDay, 
			B.CateType, B.BoxType, 
			B.FirstDate AS "DesiredDeliveryDate", B.RcvName,
			B.MainAddress, B.SubAddress, B.SubsPrice 
		FROM (SELECT UserName, UserPhone, UserNo FROM UserMst WHERE UserNo = %d) AS A
		JOIN (SELECT PeriodDay, CateType, BoxType, FirstDate, SubsPrice, UserNo,
			RcvName, MainAddress, SubAddress
			FROM SubsMst WHERE SubsNo = %d) AS B
		ON A.UserNo = B.UserNo
	;`, userNo, subsNo)
	fmt.Println(sqlQuery)

	err := database.DB.Select(&alrimTalkInfo, sqlQuery)
	if err != nil {
		return false, err.Error()
	}

	return true, alrimTalkInfo
}

// GetAlrimTalkInfoOnce : 일회 알림톡정보 획득
func GetAlrimTalkInfoOnce(userNo int, orderNo int) (isSuccess bool, result interface{}) {
	var alrimTalkInfo []AlrimTalkInfo
	sqlQuery := fmt.Sprintf(`
		SELECT A.UserName, A.UserPhone, 0 as PeriodDay,
			B.CateType, B.BoxType, 
			B.ReqDelivDate AS "DesiredDeliveryDate", B.RcvName, 
			B.MainAddress, B.SubAddress, B.OrderPrice 
		FROM (SELECT UserName, UserPhone, UserNo FROM UserMst WHERE UserNo = %d) AS A
		JOIN (SELECT CateType, BoxType, ReqDelivDate, OrderPrice, UserNo,
			RcvName, MainAddress, SubAddress
			FROM OrderMst WHERE OrderNo = %d) AS B
		ON A.UserNo = B.UserNo
	;`, userNo, orderNo)
	fmt.Println(sqlQuery)

	err := database.DB.Select(&alrimTalkInfo, sqlQuery)
	if err != nil {
		return false, err.Error()
	}

	return true, alrimTalkInfo
}

func HandleSurveyCompleted(c *gin.Context) {
	var state = gin.H{}
	utils.CombineCustomStateGlobalState(c, &state)

	userNo := utils.GetUserNo(c)

	orderType, err := c.Cookie("orderType")
	if err != nil {
		state["errorMessage"] = "주문 정보가 없습니다."
		c.HTML(http.StatusBadRequest, "error.html", state)
		return
	}

	isSuccess, infoIDStr := utils.GetIDFromCookie(c, "infoID")
	if !isSuccess {
		state["errorMessage"] = "주문 정보가 없습니다."
		c.HTML(http.StatusBadRequest, "error.html", state)
		return
	}
	infoID := infoIDStr.(int)

	var result interface{}
	if orderType == "SUBS" {
		isSuccess, result = GetAlrimTalkInfoSubs(userNo, infoID)
	} else {
		isSuccess, result = GetAlrimTalkInfoOnce(userNo, infoID)
	}
	if !isSuccess {
		slackMsg := fmt.Sprintf("[front]HandleSurveyCompleted::GetAlrimTalkInfoOnce::%s", result.(string))
		utils.SendSlackMessage(utils.SlackChannel, slackMsg)
		state["errorMessage"] = "잘못된 주문 정보."
		c.HTML(http.StatusBadRequest, "error.html", state)
		return
	}

	alrimTalkInfo := result.([]AlrimTalkInfo)
	if len(alrimTalkInfo) != 1 {
		slackMsg := fmt.Sprintf("[front]HandleSurveyCompleted::GetAlrimTalkInfoOnce::%s", "구독 완료 알림톡 행 수가 1개가 아님")
		utils.SendSlackMessage(utils.SlackChannel, slackMsg)
		state["errorMessage"] = "잘못된 주문 정보."
		c.HTML(http.StatusBadRequest, "error.html", state)
		return
	}

	// 카테고리 타입 문자열
	isSuccess, cateTypeLabels := controllers.MakeCateTypeLabels(alrimTalkInfo[0].CateType, "|")
	if !isSuccess {
		slackMsg := fmt.Sprintf("[front]HandleSurveyCompleted::controllers.MakeCateTypeLabels(%s, '|')::%s",
			alrimTalkInfo[0].CateType, cateTypeLabels)
		utils.SendSlackMessage(utils.SlackChannel, slackMsg)
		state["errorMessage"] = "시스템 에러 발생"
		state["errorCode"] = "W0015"
		c.HTML(http.StatusBadRequest, "error.html", state)
		return
	}
	cateTypeString := strings.Join(cateTypeLabels.([]string), ",")

	// 배송주기
	periodLabel := "1회"
	if alrimTalkInfo[0].PeriodDay == 7 {
		periodLabel = "매주"
	} else if alrimTalkInfo[0].PeriodDay != 0 && alrimTalkInfo[0].PeriodDay%7 == 0 {
		periodLabel = fmt.Sprintf("%d주에 한 번", alrimTalkInfo[0].PeriodDay/7)
	}

	// 박스 라벨
	isSuccess, boxLabel := controllers.GetCodeLabelFromCache("BOX_TYPE", alrimTalkInfo[0].BoxType)
	if !isSuccess {
		slackMsg := fmt.Sprintf("[front]HandleSurveyCompleted::controllers.GetCodeLabelFromCache('BOX_TYPE', %s)::%s",
			alrimTalkInfo[0].BoxType, boxLabel)
		utils.SendSlackMessage(utils.SlackChannel, slackMsg)
		state["errorMessage"] = "시스템 에러 발생"
		state["errorCode"] = "W0016"
		c.HTML(http.StatusBadRequest, "error.html", state)
		return
	}

	// 구독정보
	subsInfo := fmt.Sprintf("%s(%s)", boxLabel, cateTypeString)

	// 요일
	isSuccess, dayOfWeek := utils.GetKorDow(alrimTalkInfo[0].DesiredDeliveryDate)
	if !isSuccess {
		slackMsg := fmt.Sprintf("[front]HandleSurveyCompleted::utils.GetKorDow(%s)::%s", alrimTalkInfo[0].DesiredDeliveryDate, dayOfWeek)
		utils.SendSlackMessage(utils.SlackChannel, slackMsg)
		state["errorMessage"] = "시스템 에러 발생"
		state["errorCode"] = "W0022"
		c.HTML(http.StatusBadRequest, "error.html", state)
		return
	}
	dayOfWeek = dayOfWeek + "요일"

	// 배송희망일
	desiredDeliveryDate := strings.Split(alrimTalkInfo[0].DesiredDeliveryDate, "T")[0]

	// 배송지 주소
	address := alrimTalkInfo[0].MainAddress + " " + alrimTalkInfo[0].SubAddress

	// 알림톡 전송
	utils.SendAlrimTalk(
		alrimTalkInfo[0].UserPhone,
		"bizp_2021040611494225563208072",

		fmt.Sprintf(`%s님, 큐잇 정기구독 신청이 완료되었습니다.
%s %s, 일주일을 위한 신선한 농산물이 찾아갑니다.

수취인 : %s
배송지 : %s
배송예정일 : %s
구독정보 : %s

배송예정일에 큐잇이 도착하지 않을 경우, 고객센터로 연락주세요!
친절하게 곧바로 응대해 드리겠습니다.
고객센터 070-4166-6077`,
			alrimTalkInfo[0].UserName, periodLabel, dayOfWeek, alrimTalkInfo[0].RcvName,
			address, desiredDeliveryDate, subsInfo),
		[]utils.AlrimTalkButton{},
	)

	state["UserName"] = alrimTalkInfo[0].UserName
	state["navTitle"] = "구독완료"

	if alrimTalkInfo[0].SubsPrice != nil {
		state["priceForTracking"] = alrimTalkInfo[0].SubsPrice
	} else if alrimTalkInfo[0].OrderPrice != nil {
		state["priceForTracking"] = alrimTalkInfo[0].OrderPrice
	}

	c.HTML(http.StatusOK, "survey/completed.html", state)
	return
}
