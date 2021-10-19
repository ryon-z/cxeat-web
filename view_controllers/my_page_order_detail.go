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

// ExposedOrderUnit : 화면에 보여주기 위한 주문 유닛
type ExposedOrderUnit struct {
	models.OrderMst
	CardName       *string `db:"CardName"`
	CardNickName   *string `db:"CardNickName"`
	CardNumber     *string `db:"CardNumber"`
	PaymentDate    *string `db:"PaymentDate"`
	PaymentPrice   *int    `db:"PaymentPrice"`
	BoxLabel       string
	OrderStatus    string
	CateTypeLabels []string
}

// ExposedOrderDetail : 화면에 보여주기 위한 주문 상세
type ExposedOrderDetail struct {
	models.OrderItem
	DpName string `db:"DpName"`
}

// GetExposedOrderUnit : 화면에 보여주기 위한 주문 상세 획득
func GetExposedOrderUnit(c *gin.Context, orderNo int) (isSuccess bool, result interface{}) {
	userNo := utils.GetUserNo(c)

	var exposedOrderUnit []ExposedOrderUnit
	emptyModelArr := models.GetModelAddr("orderMst")
	fieldNamesStringA := "A." + strings.Join(envutil.GetStructFieldNames(emptyModelArr, []string{}), ", A.")
	fieldNamesString := strings.Join(envutil.GetStructFieldNames(emptyModelArr, []string{}), ", ")

	sqlQuery := fmt.Sprintf(`
		SELECT %s, C.CardName, C.CardNickName, C.CardNumber, D.RegDate AS "PaymentDate", D.PaymentPrice
		FROM (SELECT %s FROM OrderMst 
			WHERE UserNo = %d AND OrderNo = %d) AS A
		LEFT JOIN UserCard AS C
		ON (A.CardRegNo = C.CardRegNo)
		LEFT JOIN (SELECT OrderNo, PaymentPrice, RegDate FROM OrderPayment
			WHERE StatusCode = "normal") AS D
		ON A.OrderNo = D.OrderNo
	;`, fieldNamesStringA, fieldNamesString, userNo, orderNo)
	fmt.Println(sqlQuery)

	err := database.DB.Select(&exposedOrderUnit, sqlQuery)
	if err != nil {
		return false, err.Error()
	}

	return true, exposedOrderUnit
}

// GetExposedOrderDetail : 화면에 보여주기 위한 주문 아이템 상세 목록 획득
func GetExposedOrderDetail(c *gin.Context, orderNo int) (isSuccess bool, result interface{}) {
	var exposedOrderDetail []ExposedOrderDetail
	emptyModelArr := models.GetModelAddr("orderItem")
	fieldNamesString := "B." + strings.Join(envutil.GetStructFieldNames(emptyModelArr, []string{}), ", B.")

	sqlQuery := fmt.Sprintf(`
		SELECT %s, A.DpName
		FROM ItemMst AS A
		JOIN OrderItem AS B
		ON. A.ItemNo = B.ItemNo
		where OrderNo = %d
	;`, fieldNamesString, orderNo)
	fmt.Println(sqlQuery)

	err := database.DB.Select(&exposedOrderDetail, sqlQuery)
	if err != nil {
		return false, err.Error()
	}

	return true, exposedOrderDetail
}

// HandleMyPageOrderDetail : my-page/order/detail 뷰
func HandleMyPageOrderDetail(c *gin.Context) {
	state := gin.H{}
	utils.CombineCustomStateGlobalState(c, &state)
	state["orderInfo"] = ""
	state["orderDetails"] = ""
	state["userAddress"] = ""
	state["userCard"] = ""

	// orderNo 획득
	isSuccess, cookieResult := utils.GetIDFromCookie(c, "orderNo")
	if !isSuccess {
		state["errorMessage"] = "선택된 주문이 없습니다. 주문내역으로 이동하여 주문을 선택해주세요."
		c.HTML(http.StatusBadRequest, "error.html", state)
		return
	}
	orderNo := cookieResult.(int)

	isSuccess, result := GetExposedOrderUnit(c, orderNo)
	if !isSuccess {
		slackMsg := fmt.Sprintf("[front]HandleMyPageOrderDetail::GetExposedOrderUnit::%s", result.(string))
		utils.SendSlackMessage(utils.SlackChannel, slackMsg)
		state["errorMessage"] = "시스템 에러 발생"
		state["errorCode"] = "W0014"
		c.HTML(http.StatusBadRequest, "error.html", state)
		return
	}

	var orderUnit = result.([]ExposedOrderUnit)
	if len(orderUnit) > 0 {
		raw, _ := time.Parse(time.RFC3339, orderUnit[0].ReqDelivDate)
		date := fmt.Sprintf("%04d-%02d-%02d", raw.Year(), raw.Month(), raw.Day())
		orderUnit[0].ReqDelivDate = date

		// 카테고리 타입 문자열
		isSuccess, cateTypeLabels := controllers.MakeCateTypeLabels(orderUnit[0].CateType, "|")
		if !isSuccess {
			slackMsg := fmt.Sprintf("[front]HandleMyPageOrderDetail::controllers.MakeCateTypeLabels(%s, '|')::%s",
				orderUnit[0].CateType, cateTypeLabels.(string))
			utils.SendSlackMessage(utils.SlackChannel, slackMsg)
			state["errorMessage"] = "시스템 에러 발생"
			state["errorCode"] = "W0015"
			c.HTML(http.StatusBadRequest, "error.html", state)
			return
		}
		for _, row := range cateTypeLabels.([]string) {
			orderUnit[0].CateTypeLabels = append(orderUnit[0].CateTypeLabels, row)
		}

		// 박스 라벨
		isSuccess, boxLabel := controllers.GetCodeLabelFromCache("BOX_TYPE", orderUnit[0].BoxType)
		if !isSuccess {
			slackMsg := fmt.Sprintf("[front]HandleMyPageOrderDetail::controllers.GetCodeLabelFromCache('BOX_TYPE', %s)::%s",
				orderUnit[0].BoxType, boxLabel)
			utils.SendSlackMessage(utils.SlackChannel, slackMsg)
			state["errorMessage"] = "시스템 에러 발생"
			state["errorCode"] = "W0016"
			c.HTML(http.StatusBadRequest, "error.html", state)
			return
		}
		orderUnit[0].BoxLabel = boxLabel

		// 주문 상태
		isSuccess, orderStatus := controllers.GetCodeLabelFromCache("ORDER_STATUS", orderUnit[0].StatusCode)
		if !isSuccess {
			slackMsg := fmt.Sprintf("[front]HandleMyPageOrderDetail::controllers.GetCodeLabelFromCache('ORDER_STATUS', %s)::%s",
				orderUnit[0].StatusCode, orderStatus)
			utils.SendSlackMessage(utils.SlackChannel, slackMsg)
			state["errorMessage"] = "시스템 에러 발생"
			state["errorCode"] = "W0016"
			c.HTML(http.StatusBadRequest, "error.html", state)
			return
		}
		orderUnit[0].StatusCode = orderStatus

		// 결제 일자
		// 결제 전일 경우 PaymentDate는 nil
		if orderUnit[0].PaymentDate != nil {
			paymentDate := strings.Split(*orderUnit[0].PaymentDate, "T")[0]
			orderUnit[0].PaymentDate = &paymentDate
		}

		state["orderInfo"] = orderUnit[0]
	}

	isSuccess, result = GetExposedOrderDetail(c, orderNo)
	if !isSuccess {
		slackMsg := fmt.Sprintf("[front]HandleMyPageOrderDetail::GetExposedOrderDetail::%s", result.(string))
		utils.SendSlackMessage(utils.SlackChannel, slackMsg)
		state["errorMessage"] = "시스템 에러 발생"
		state["errorCode"] = "W0017"
		c.HTML(http.StatusBadRequest, "error.html", state)
		return
	}
	var orderDetails = result.([]ExposedOrderDetail)
	state["orderDetails"] = orderDetails

	state["activeNav"] = "myPage"
	state["navTitle"] = "주문 내역 상세"

	// 택배사 링크
	var delivURL string
	if orderUnit[0].DelivCo != nil && orderUnit[0].DelivInvoiceNo != nil {
		switch s := *orderUnit[0].DelivCo; s {
		case "한진택배":
			delivURL = "//www.hanjin.co.kr/kor/CMS/DeliveryMgr/WaybillResult.do?mCode=MN038&schLang=KR&wblnum=" + *orderUnit[0].DelivInvoiceNo
		case "롯데택배":
			delivURL = "//www.lotteglogis.com/mobile/reservation/tracking/linkView?InvNo=" + *orderUnit[0].DelivInvoiceNo
		case "CJ대한통운":
			delivURL = "//www.cjlogistics.com/ko/tool/parcel/tracking?gnbInvcNo=" + *orderUnit[0].DelivInvoiceNo
		case "우체국":
			delivURL = "//service.epost.go.kr/trace.RetrieveDomRigiTraceList.comm?sid1=" + *orderUnit[0].DelivInvoiceNo
		case "로젠":
			delivURL = "//www.ilogen.com/m/personal/trace/" + *orderUnit[0].DelivInvoiceNo
		default:
			delivURL = ""
		}
	}
	state["delivURL"] = delivURL

	c.HTML(http.StatusOK, "my-page/order-detail.html", state)
}
