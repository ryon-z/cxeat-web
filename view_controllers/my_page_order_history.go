package viewcontrollers

import (
	"fmt"
	"net/http"
	"strings"
	"time"
	"yelloment-api/controllers"
	"yelloment-api/database"
	envutil "yelloment-api/env_util"
	"yelloment-api/models"
	"yelloment-api/utils"

	"github.com/gin-gonic/gin"
)

// ExposedOrders : 화면에 보여주기 위한 주문 목록
type ExposedOrders struct {
	models.OrderMst
	BoxLabel       string
	OrderStatus    string
	CateTypeLabels []string
}

// GetExposedOrders : 화면에 보여주기 위한 주문 목록 획득
func GetExposedOrders(c *gin.Context) (isSuccess bool, result interface{}) {
	userNo := utils.GetUserNo(c)

	var exposedOrders []ExposedOrders
	emptyModelArr := models.GetModelAddr("orderMst")
	fieldNamesString := "A." + strings.Join(envutil.GetStructFieldNames(emptyModelArr, []string{}), ", A.")

	sqlQuery := fmt.Sprintf(`
		SELECT %s, RcvName
		FROM OrderMst AS A
		WHERE (A.UserNo = %d)
		ORDER BY RegDate DESC
	;`, fieldNamesString, userNo)
	fmt.Println(sqlQuery)

	err := database.DB.Select(&exposedOrders, sqlQuery)
	if err != nil {
		return false, err.Error()
	}

	return true, exposedOrders
}

// HandleMyPageOrderHistory : my-page/order/history 뷰
func HandleMyPageOrderHistory(c *gin.Context) {
	state := gin.H{}
	utils.CombineCustomStateGlobalState(c, &state)
	state["orderInfos"] = ""

	isSuccess, result := GetExposedOrders(c)
	if !isSuccess {
		slackMsg := fmt.Sprintf("[front]HandleMyPageOrderHistory::GetExposedOrders::%s", result.(string))
		utils.SendSlackMessage(utils.SlackChannel, slackMsg)
		state["errorMessage"] = "시스템 에러 발생"
		state["errorCode"] = "W0019"
		c.HTML(http.StatusBadRequest, "error.html", state)
		return
	}
	var orderInfos = result.([]ExposedOrders)
	for i, orderInfo := range orderInfos {
		// 배송 날짜 세팅
		raw, _ := time.Parse(time.RFC3339, orderInfo.ReqDelivDate)
		date := fmt.Sprintf("%04d-%02d-%02d", raw.Year(), raw.Month(), raw.Day())
		orderInfos[i].ReqDelivDate = date

		// 주문일
		raw, _ = time.Parse(time.RFC3339, orderInfo.RegDate)
		date = fmt.Sprintf("%04d-%02d-%02d", raw.Year(), raw.Month(), raw.Day())
		orderInfos[i].RegDate = date

		// 카테고리 타입 문자열
		isSuccess, cateTypeLabels := controllers.MakeCateTypeLabels(orderInfo.CateType, "|")
		if !isSuccess {
			slackMsg := fmt.Sprintf("[front]HandleMyPageOrderHistory::controllers.MakeCateTypeLabels(%s, '|')::%s",
				orderInfo.CateType, cateTypeLabels.(string))
			utils.SendSlackMessage(utils.SlackChannel, slackMsg)
			state["errorMessage"] = "시스템 에러 발생"
			state["errorCode"] = "W0015"
			c.HTML(http.StatusBadRequest, "error.html", state)
			return
		}
		for _, row := range cateTypeLabels.([]string) {
			orderInfos[i].CateTypeLabels = append(orderInfos[i].CateTypeLabels, row)
		}

		// 박스 라벨
		isSuccess, boxLabel := controllers.GetCodeLabelFromCache("BOX_TYPE", orderInfo.BoxType)
		if !isSuccess {
			slackMsg := fmt.Sprintf("[front]HandleMyPageOrderHistory::controllers.GetCodeLabelFromCache('BOX_TYPE', %s)::%s",
				orderInfo.BoxType, boxLabel)
			utils.SendSlackMessage(utils.SlackChannel, slackMsg)
			state["errorMessage"] = "시스템 에러 발생"
			state["errorCode"] = "W0016"
			c.HTML(http.StatusBadRequest, "error.html", state)
			return
		}
		orderInfos[i].BoxLabel = boxLabel

		// 주문 상태
		isSuccess, orderStatus := controllers.GetCodeLabelFromCache("ORDER_STATUS", orderInfo.StatusCode)
		if !isSuccess {
			slackMsg := fmt.Sprintf("[front]HandleMyPageOrderHistory::controllers.GetCodeLabelFromCache('ORDER_STATUS', %s)::%s",
				orderInfo.StatusCode, orderStatus)
			utils.SendSlackMessage(utils.SlackChannel, slackMsg)
			state["errorMessage"] = "시스템 에러 발생"
			state["errorCode"] = "W0016"
			c.HTML(http.StatusBadRequest, "error.html", state)
			return
		}
		orderInfos[i].OrderStatus = orderStatus
	}
	state["orderInfos"] = orderInfos
	state["activeNav"] = "myPage"
	state["navTitle"] = "주문 내역"

	c.HTML(http.StatusOK, "my-page/order-history.html", state)
}
