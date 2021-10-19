package controllers

import (
	"fmt"
	"net/http"
	"strings"
	"yelloment-api/middleware/mixin"
	"yelloment-api/models"
	"yelloment-api/utils"

	"github.com/gin-gonic/gin"
)

// CreateOnceOwnedOrderMst : 소유자 단건 유저 주문 생성
func CreateOnceOwnedOrderMst(c *gin.Context) {
	userNo := utils.GetUserNo(c)

	var orderMst models.OrderMst
	if err := c.ShouldBindJSON(&orderMst); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	isSuccess, message := CheckCardRegNoValid(userNo, orderMst.CardRegNo)
	if !isSuccess {
		c.JSON(http.StatusBadRequest, gin.H{"error": message})
		return
	}

	orderMst.StatusCode = "init"
	orderMst.UserNo = userNo
	subsNo := 0
	orderMst.SubsNo = &subsNo
	orderRound := 0
	orderMst.OrderRound = &orderRound
	orderMst.OrderType = "ONCE"
	orderMst.UserNo = userNo
	isSuccess, result := mixin.CreateOwned("orderMst", userNo, &orderMst, []string{"OrderNo", "RegDate"}, true, []string{})

	// 내부용 1회 주문 알람
	hostname := strings.Split(c.Request.Host, ":")[0]
	if isSuccess && hostname != "test.cueat.kr" && hostname != "localhost" {
		orderAdminURL := fmt.Sprintf("https://admin.cueat.kr/order/view/%v", result)
		message := fmt.Sprintf("1회 주문(<%s|%d>)이 신청되었습니다.", orderAdminURL, result)
		utils.SendSlackMessage("internal", message)
	}

	utils.GetHTTPResponse(isSuccess, result, c)
}

// GetOwnedOrder : 소유자 주문 획득
func GetOwnedOrder(c *gin.Context) (isSuccess bool, result interface{}) {
	userNo := utils.GetUserNo(c)
	wherePhrase := fmt.Sprintf(`UserNo = %d AND (StatusCode != "cancel")`, userNo)
	return mixin.List("orderMst", wherePhrase, "")
}

// ListOwnedOrder : 소유자 주문 리스트 조회
func ListOwnedOrder(c *gin.Context) {
	isSuccess, result := GetOwnedOrder(c)
	utils.GetHTTPResponse(isSuccess, result, c)
}

// Not used
// GetOwnedOrderUnit : 소유자 주문 단일 조회(주문ID는 쿠키)
func GetOwnedOrderUnit(c *gin.Context, orderNo int) (isSuccess bool, result interface{}) {
	// isSuccess, result = utils.GetIDFromCookie(c, "orderNo")
	// if !isSuccess {
	// 	return false, result.(string)
	// }
	// orderNo := result.(int)

	userNo := utils.GetUserNo(c)
	wherePhrase := fmt.Sprintf(`UserNo = %d AND OrderNo = %d`, userNo, orderNo)
	isSuccess, result = mixin.List("orderMst", wherePhrase, "")
	if !isSuccess {
		return false, result.(string)
	}

	return true, result
}

// DeactivateOwnedOrder : 소유자 주문 취소
func DeactivateOwnedOrder(c *gin.Context) {
	userNo := utils.GetUserNo(c)

	var orderMst = models.OrderMst{}
	var orderMstReq models.OrderMst
	if err := c.ShouldBindJSON(&orderMstReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	orderMst.OrderNo = orderMstReq.OrderNo
	orderMst.StatusCode = "cancel"
	wherePhrase := fmt.Sprintf(`UserNo = %d AND OrderNo = %d`, userNo, orderMst.OrderNo)
	isSuccess, message := mixin.PartialUpdate("orderMst", wherePhrase, &orderMst, []string{}, []string{})
	utils.GetHTTPResponse(isSuccess, message, c)
}
