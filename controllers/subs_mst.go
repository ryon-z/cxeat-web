package controllers

import (
	"fmt"
	"net/http"
	"strings"
	"time"
	envutil "yelloment-api/env_util"
	"yelloment-api/middleware/mixin"
	"yelloment-api/models"
	"yelloment-api/utils"

	"github.com/gin-gonic/gin"
)

// GetOwnedSubsMstUnit : 소유자 유저 구독 행 조회(구독ID는 쿠키)
func GetOwnedSubsMstUnit(c *gin.Context, subsNo int) (isSuccess bool, result interface{}) {
	userNo := utils.GetUserNo(c)
	wherePhrase := fmt.Sprintf(`UserNo = %d AND SubsNo = %d`, userNo, subsNo)
	isSuccess, result = mixin.List("subsMst", wherePhrase, "")
	if !isSuccess {
		return false, "Loading API is failed"
	}

	return true, result
}

// GetOwnedSubsMst : 소유자 유저 구독 리스트 획득
func GetOwnedSubsMst(c *gin.Context) (isSuccess bool, result interface{}) {
	userNo := utils.GetUserNo(c)
	wherePhrase := fmt.Sprintf(`UserNo = %d AND (StatusCode = "normal" OR StatusCode = "pause")`, userNo)
	return mixin.List("subsMst", wherePhrase, "")
}

// ListOwnedSubsMst : 소유자 유저 구독 리스트 조회
func ListOwnedSubsMst(c *gin.Context) {
	isSuccess, result := GetOwnedSubsMst(c)
	utils.GetHTTPResponse(isSuccess, result, c)
}

// CreateOwnedSubsMst : 소유자 유저 구독 생성
func CreateOwnedSubsMst(c *gin.Context) {
	userNo := utils.GetUserNo(c)

	var subsMst models.SubsMst
	if err := c.ShouldBindJSON(&subsMst); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	isSuccess, message := CheckCardRegNoValid(userNo, subsMst.CardRegNo)
	if !isSuccess {
		c.JSON(http.StatusBadRequest, gin.H{"error": message})
		return
	}

	subsMst.StatusCode = "normal"
	subsMst.UserNo = userNo
	isSuccess, result := mixin.CreateOwned("subsMst", userNo, &subsMst, []string{"SubsNo", "RegDate"}, true, []string{"SubsType"})

	// 내부용 구독 알람
	hostname := strings.Split(c.Request.Host, ":")[0]
	if isSuccess && hostname != "test.cueat.kr" && hostname != "localhost" {
		subsAdminURL := fmt.Sprintf("http://admin.cueat.kr/subs/view/%v", result)
		message := fmt.Sprintf("구독(<%s|%v>)이 신청되었습니다", subsAdminURL, result)
		utils.SendSlackMessage("internal", message)
	}

	utils.GetHTTPResponse(isSuccess, result, c)
}

// PauseOwnedSubsMst : 소유자 유저 구독 일시정지
func PauseOwnedSubsMst(c *gin.Context) {
	userNo := utils.GetUserNo(c)

	var subsMst = models.SubsMst{}
	var subsMstReq models.SubsMst
	if err := c.ShouldBindJSON(&subsMst); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	subsMst.SubsNo = subsMstReq.SubsNo
	subsMst.StatusCode = "pause"
	wherePhrase := fmt.Sprintf(`UserNo = %d AND SubsNo = %d`, userNo, subsMst.SubsNo)
	isSuccess, message := mixin.PartialUpdate("subsMst", wherePhrase, &subsMst, []string{}, []string{})
	utils.GetHTTPResponse(isSuccess, message, c)
}

// DeactivateOwnedSubsMst : 소유자 유저 구독 비활성화
func DeactivateOwnedSubsMst(c *gin.Context) {
	userNo := utils.GetUserNo(c)

	var subsMst = models.SubsMst{}
	var subsMstReq models.SubsMst
	if err := c.ShouldBindJSON(&subsMstReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Println("subsMstReq", subsMstReq)

	// 업데이트일 및 취소일 갱신
	loc, _ := time.LoadLocation(envutil.GetGoDotEnvVariable("TIMEZONE"))
	now := time.Now().In(loc)
	currentTime := now.Format("2006-01-02 15:04:05")
	subsMst.UpdDate = &currentTime
	subsMst.CnlDate = &currentTime
	subsMst.CnlReason = subsMstReq.CnlReason

	subsMst.SubsNo = subsMstReq.SubsNo
	subsMst.StatusCode = "cancel"
	wherePhrase := fmt.Sprintf(`UserNo = %d AND SubsNo = %d`, userNo, subsMstReq.SubsNo)
	isSuccess, message := mixin.PartialUpdate("subsMst", wherePhrase, &subsMst, []string{}, []string{})

	// SubsHist에 행 삽입
	if isSuccess {
		var subsHist models.SubsHist
		statusCode := "cancel"
		histDesc := "구독 해지"

		subsHist.SubsNo = subsMst.SubsNo
		subsHist.StatusCode = &statusCode
		subsHist.HistDesc = &histDesc
		subsHist.ExecUser = "user"
		subsHist.ExecDate = currentTime
		isSuccessHist, result := mixin.Create("subsHist", &subsHist, []string{}, true, []string{})
		if !isSuccessHist {
			slackMsg := fmt.Sprintf("[front]DeactivateOwnedSubsMst::mixin.Create('subsHist'),::%v", result)
			utils.SendSlackMessage(utils.SlackChannel, slackMsg)
		}

		// 내부용 구독 알람
		hostname := strings.Split(c.Request.Host, ":")[0]
		if hostname != "test.cueat.kr" && hostname != "localhost" {
			subsAdminURL := fmt.Sprintf("http://admin.cueat.kr/subs/view/%d", subsMst.SubsNo)
			message := fmt.Sprintf("구독(<%s|%d>)이 해지되었습니다", subsAdminURL, subsMst.SubsNo)
			utils.SendSlackMessage("internal", message)
		}
	}

	utils.GetHTTPResponse(isSuccess, message, c)
}

// PartialUpdateOwnedSubsMst : 소유자 유저 구독 수정
func PartialUpdateOwnedSubsMst(c *gin.Context) {
	var isSuccess bool
	var message interface{}

	userNo := utils.GetUserNo(c)

	var subsMst models.SubsMst
	if err := c.ShouldBindJSON(&subsMst); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	subsNo := subsMst.SubsNo

	if subsMst.CardRegNo != 0 {
		isSuccess, message = CheckCardRegNoValid(userNo, subsMst.CardRegNo)
		if !isSuccess {
			c.JSON(http.StatusBadRequest, gin.H{"error": message})
			return
		}
	}

	wherePhrase := fmt.Sprintf(`UserNo = %d AND SubsNo = %d`, userNo, subsNo)
	isSuccess, message = mixin.PartialUpdate("subsMst", wherePhrase, &subsMst, []string{}, []string{"SubsType"})

	utils.GetHTTPResponse(isSuccess, message, c)
}
