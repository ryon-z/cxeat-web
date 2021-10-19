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

// ExposedSubsUnit : 단일 구독 정보
type ExposedSubsUnit struct {
	models.SubsMst
	CardName       *string `db:"CardName"`
	CardNickName   *string `db:"CardNickName"`
	CardNumber     *string `db:"CardNumber"`
	BoxLabel       string
	PeriodLabel    string
	DayOfWeek      string
	SubsStatus     string
	CateTypeLabels []string
}

// GetExposedSubsUnit : 화면에 보여주기 위한 구독 목록 획득
func GetExposedSubsUnit(c *gin.Context, subsNo int) (isSuccess bool, result interface{}) {
	userNo := utils.GetUserNo(c)

	var exposedSubsUnit []ExposedSubsUnit
	emptyModelArr := models.GetModelAddr("subsMst")
	fieldNamesStringA := "A." + strings.Join(envutil.GetStructFieldNames(emptyModelArr, []string{}), ", A.")
	fieldNamesString := strings.Join(envutil.GetStructFieldNames(emptyModelArr, []string{}), ", ")

	sqlQuery := fmt.Sprintf(`
		SELECT %s, C.CardName, C.CardNickName, C.CardNumber
		FROM 
			(SELECT %s FROM SubsMst
			WHERE SubsNo = %d AND UserNo = %d) AS A
		LEFT JOIN UserCard AS C
		ON (A.CardRegNo = C.CardRegNo)
	;`, fieldNamesStringA, fieldNamesString, subsNo, userNo)

	err := database.DB.Select(&exposedSubsUnit, sqlQuery)
	if err != nil {
		return false, err.Error()
	}

	return true, exposedSubsUnit
}

// HandleMyPageSubsInfo : my-page/subs/info 뷰
func HandleMyPageSubsInfo(c *gin.Context) {
	state := gin.H{}
	utils.CombineCustomStateGlobalState(c, &state)
	state["subsInfo"] = ""
	state["userAddresses"] = ""
	state["userCards"] = ""

	// subsNo 획득
	isSuccess, cookieResult := utils.GetIDFromCookie(c, "subscriptionNo")
	if !isSuccess {
		state["errorMessage"] = "선택한 구독이 없습니다. 구독관리로 이동하여 구독을 선택해주세요."
		c.HTML(http.StatusBadRequest, "error.html", state)
		return
	}
	subsNo := cookieResult.(int)

	isSuccess, result := GetExposedSubsUnit(c, subsNo)
	if !isSuccess {
		slackMsg := fmt.Sprintf("[front]HandleMyPageSubsInfo::GetExposedSubsUnit::%s", result.(string))
		utils.SendSlackMessage(utils.SlackChannel, slackMsg)
		state["errorMessage"] = "시스템 에러 발생"
		state["errorCode"] = "W0020"
		c.HTML(http.StatusBadRequest, "error.html", state)
		return
	}

	subsInfos := result.([]ExposedSubsUnit)
	if len(subsInfos) < 1 {
		state["errorMessage"] = "잘못된 구독 정보 요청"
		c.HTML(http.StatusBadRequest, "error.html", state)
		return
	}
	subsInfo := subsInfos[0]

	// 카테고리 타입 문자열
	isSuccess, cateTypeLabels := controllers.MakeCateTypeLabels(subsInfo.CateType, "|")
	if !isSuccess {
		slackMsg := fmt.Sprintf("[front]HandleMyPageSubsInfo::controllers.MakeCateTypeLabels(%s, '|')::%s",
			subsInfo.CateType, cateTypeLabels)
		utils.SendSlackMessage(utils.SlackChannel, slackMsg)
		state["errorMessage"] = "시스템 에러 발생"
		state["errorCode"] = "W0015"
		c.HTML(http.StatusBadRequest, "error.html", state)
		return
	}
	for _, row := range cateTypeLabels.([]string) {
		subsInfo.CateTypeLabels = append(subsInfo.CateTypeLabels, row)
	}

	// 박스 라벨
	isSuccess, boxLabel := controllers.GetCodeLabelFromCache("BOX_TYPE", subsInfo.BoxType)
	if !isSuccess {
		slackMsg := fmt.Sprintf("[front]HandleMyPageSubsInfo::controllers.GetCodeLabelFromCache('BOX_TYPE', %s)::%s",
			subsInfo.BoxType, boxLabel)
		utils.SendSlackMessage(utils.SlackChannel, slackMsg)
		state["errorMessage"] = "시스템 에러 발생"
		state["errorCode"] = "W0016"
		c.HTML(http.StatusBadRequest, "error.html", state)
		return
	}
	subsInfo.BoxLabel = boxLabel

	// 배송 주기
	isSuccess, periodLabel := controllers.GetCodeLabelFromCache("DELIVERY_PERIOD", strconv.Itoa(subsInfo.PeriodDay))
	if !isSuccess {
		slackMsg := fmt.Sprintf("[front]HandleMyPageSubsInfo::controllers.GetCodeLabelFromCache('DELIVERY_PERIOD', %d)::%s",
			subsInfo.PeriodDay, periodLabel)
		utils.SendSlackMessage(utils.SlackChannel, slackMsg)
		state["errorMessage"] = "시스템 에러 발생"
		state["errorCode"] = "W0016"
		c.HTML(http.StatusBadRequest, "error.html", state)
		return
	}
	subsInfo.PeriodLabel = periodLabel
	state["subsInfo"] = subsInfo

	// 구독 상태
	if subsInfo.StatusCode == "normal" {
		subsInfo.SubsStatus = "구독 중"
	} else {
		subsInfo.SubsStatus = "구독 취소"
	}

	// 배송요일
	isSuccess, dowResult := utils.GetKorDow(subsInfo.FirstDate)
	if !isSuccess {
		slackMsg := fmt.Sprintf("[front]HandleMyPageSubsInfo::utils.GetKorDow(%s)::%s", subsInfo.FirstDate, dowResult)
		utils.SendSlackMessage(utils.SlackChannel, slackMsg)
		state["errorMessage"] = "시스템 에러 발생"
		state["errorCode"] = "W0022"
		c.HTML(http.StatusBadRequest, "error.html", state)
		return
	}
	subsInfo.DayOfWeek = dowResult

	state["subsInfo"] = subsInfo

	// 유저 소유 주소지 목록
	isSuccess, result = controllers.GetOwnedUserAddress(c)
	if !isSuccess {
		slackMsg := fmt.Sprintf("[front]HandleMyPageSubsInfo::controllers.GetOwnedUserAddress::%s", result.(string))
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
		slackMsg := fmt.Sprintf("[front]HandleMyPageSubsInfo::controllers.GetOwnedUserCard::%s", result.(string))
		utils.SendSlackMessage(utils.SlackChannel, slackMsg)
		state["errorMessage"] = "시스템 에러 발생"
		state["errorCode"] = "W0001"
		c.HTML(http.StatusBadRequest, "error.html", state)
		return
	}

	var userCards = *result.(*[]models.UserCard)
	state["userCards"] = userCards
	state["activeNav"] = "myPage"
	state["navTitle"] = "구독 정보 상세"

	c.HTML(http.StatusOK, "my-page/subs-info.html", state)
}
