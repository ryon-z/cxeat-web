package viewcontrollers

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"yelloment-api/controllers"
	"yelloment-api/database"
	envutil "yelloment-api/env_util"
	"yelloment-api/middleware/mixin"
	"yelloment-api/models"
	"yelloment-api/utils"

	"github.com/gin-gonic/gin"
)

// ReviewItems : 리뷰할 주문 아이템
type ReviewItems struct {
	ItemNo      int      `db:"ItemNo"`
	ItemCnt     int      `db:"ItemCnt"`
	ReviewNo    *int     `db:"ReviewNo"`
	ReviewScore *float64 `db:"ReviewScore"`
	ReviewDesc  *string  `db:"ReviewDesc"`
	DpName      string   `db:"DpName"`
}

// ExposedOrderInfo : 뷰에서 그릴 때 필요한 주문 정보
type ExposedOrderInfo struct {
	OrderNo        int    `db:"OrderNo" json:"OrderNo"`
	CateType       string `db:"CateType" json:"CateType"`
	RcvName        string `db:"RcvName" json:"RcvName"`
	OrderRound     *int   `db:"OrderRound" json:"OrderRound"`
	PaymentDate    string `db:"PaymentDate" json:"PaymentDate"`
	CateTypeLabels []string
}

// GetReviewItems : 리뷰할 주문 아이템 획득
func GetReviewItems(userNo int, orderNo int) (isSuccess bool, result interface{}) {
	var reviewItems []ReviewItems
	sqlQuery := fmt.Sprintf(`
		SELECT A.ItemNo, A.ItemCnt, B.ReviewNo, B.ReviewScore, B.ReviewDesc, C.DpName
		FROM
			(SELECT OrderNo, ItemNo, ItemCnt FROM OrderItem WHERE OrderNo = %d) as A
		LEFT JOIN
			(SELECT OrderNo, ItemNo, ReviewNo, ReviewScore, ReviewDesc FROM ItemReview WHERE OrderNo = %d) as B
		ON A.OrderNo = B.OrderNo AND A.ItemNo = B.ItemNo
		JOIN (SELECT ItemNo, DpName FROM ItemMst) as C
		ON A.ItemNo = C.ItemNo
		JOIN (SELECT OrderNo FROM OrderMst WHERE OrderNo = %d AND userNo = %d) as D
		ON A.OrderNo = D.OrderNo
	;`, orderNo, orderNo, orderNo, userNo)
	fmt.Println(sqlQuery)

	err := database.DB.Select(&reviewItems, sqlQuery)
	if err != nil {
		return false, err.Error()
	}

	return true, reviewItems
}

// GetExposedOrderInfo : 쿠키로 부터 exposedOrderInfo 획득
func GetExposedOrderInfo(c *gin.Context) (isSuccess bool, result interface{}) {
	orderInfo, err := c.Request.Cookie("orderInfo")
	if err != nil {
		return false, "orderInfo cookie not exists"
	}

	// 뷰에서 cookie로 JSON을 넘길 때 JSON.stringify만 하면 uri decoding 시
	// 에러가 나기 때문에 JSON.stringify 결과를 한 번 더 uriEncoding해서 송부함
	// 그래서 decoding도 두 번 해야 함
	decoded, err := url.QueryUnescape(orderInfo.Value)
	if err != nil {
		return false, err.Error()
	}
	decoded, err = url.QueryUnescape(decoded)
	if err != nil {
		return false, err.Error()
	}

	var exposedOrderInfo ExposedOrderInfo

	json.Unmarshal([]byte(decoded), &exposedOrderInfo)

	return true, exposedOrderInfo
}

// GetReviewDetailState : /web/my-page/review/detail 뷰 state 획득
func GetReviewDetailState(c *gin.Context, state gin.H) (isSuccess bool, result interface{}) {
	// ExposedOrderInfo 획득
	isSuccess, result = GetExposedOrderInfo(c)
	if !isSuccess {
		if result.(string) == "orderInfo cookie not exists" {
			state["errorMessage"] = "선택한 리뷰가 없습니다. 나의 리뷰로 이동하여 리뷰를 선택해주세요."
			return false, state
		}

		slackMsg := fmt.Sprintf("[front]HandleMyPageReviewDetail::GetReviewDetailState::%s", result.(string))
		utils.SendSlackMessage(utils.SlackChannel, slackMsg)
		state["errorMessage"] = "시스템 에러 발생"
		state["errorCode"] = "W0028"
		return false, state
	}
	var exposedOrderInfo = result.(ExposedOrderInfo)

	// 카테고리 타입 문자열
	isSuccess, cateTypeLabels := controllers.MakeCateTypeLabels(exposedOrderInfo.CateType, "|")
	if !isSuccess {
		slackMsg := fmt.Sprintf("[front]HandleMyPageReviewList::controllers.MakeCateTypeLabels(%s, '|')::%s",
			exposedOrderInfo.CateType, cateTypeLabels.(string))
		utils.SendSlackMessage("system", slackMsg)
		state["errorMessage"] = "시스템 에러 발생"
		state["errorCode"] = "W0015"

		c.Redirect(http.StatusTemporaryRedirect, "/web/main")
		return
	}
	for _, row := range cateTypeLabels.([]string) {
		exposedOrderInfo.CateTypeLabels = append(exposedOrderInfo.CateTypeLabels, row)
	}

	// userNo 획득
	userNo := utils.GetUserNo(c)

	isSuccess, result = GetReviewItems(userNo, exposedOrderInfo.OrderNo)
	if !isSuccess {
		slackMsg := fmt.Sprintf(
			"[front]HandleMyPageReviewDetail::GetReviewItems(%d, %d)::%s",
			userNo, exposedOrderInfo.OrderNo, result.(string),
		)
		utils.SendSlackMessage(utils.SlackChannel, slackMsg)
		state["errorMessage"] = "시스템 에러 발생"
		state["errorCode"] = "W0029"
		return false, state
	}
	reviewItems := result.([]ReviewItems)

	isSuccess, cookieResult := utils.GetIDFromCookie(c, "reviewNo")
	if !isSuccess {
		state["errorMessage"] = "선택한 리뷰가 없습니다. 나의 리뷰로 이동하여 리뷰를 선택해주세요."
		return false, state
	}
	reviewNo := cookieResult.(int)

	// review.ReviewDesc 획득
	isSuccess, result = controllers.GetReviewMstUnit(c, reviewNo)
	if !isSuccess {
		slackMsg := fmt.Sprintf(
			"[front]HandleMyPageReviewDetail::GetReviewMstUnit(c, %d)::%s",
			reviewNo, result.(string),
		)
		utils.SendSlackMessage(utils.SlackChannel, slackMsg)
		state["errorMessage"] = "시스템 에러 발생"
		state["errorCode"] = "W0030"
		return false, state
	}
	review := *result.(*[]models.ReviewMst)
	if len(review) != 1 {
		slackMsg := fmt.Sprintf(
			"[front]HandleMyPageReviewDetail::GetReviewMstUnit(c, %d)::%s",
			reviewNo, result.(string),
		)
		utils.SendSlackMessage(utils.SlackChannel, slackMsg)
		state["errorMessage"] = "시스템 에러 발생"
		state["errorCode"] = "W0031"
		return false, state
	}

	state["reviewDesc"] = review[0].ReviewDesc
	state["reviewItems"] = reviewItems
	state["orderInfo"] = exposedOrderInfo
	state["reviewNo"] = reviewNo
	state["navTitle"] = "나의 리뷰 상세"

	return true, state
}

// HandleMyPageReviewDetail : /web/my-page/review/detail 뷰
func HandleMyPageReviewDetail(c *gin.Context) {
	state := gin.H{}
	utils.CombineCustomStateGlobalState(c, &state)

	isSuccess, result := GetReviewDetailState(c, state)
	if !isSuccess {
		c.HTML(http.StatusBadRequest, "error.html", result.(gin.H))
		return
	}
	resultState := result.(gin.H)
	utils.CombineTwoGinH(&state, &resultState)

	c.HTML(http.StatusOK, "my-page/review-detail.html", state)
	return
}

// HandleModifyReview : 리뷰 수정
func HandleModifyReview(c *gin.Context) {
	state := gin.H{}
	utils.CombineCustomStateGlobalState(c, &state)

	isSuccess, result := GetReviewDetailState(c, state)
	if !isSuccess {
		c.HTML(http.StatusBadRequest, "error.html", result.(gin.H))
		return
	}
	resultState := result.(gin.H)
	utils.CombineTwoGinH(&state, &resultState)
	state["userNo"] = utils.GetUserNo(c)
	state["navTitle"] = "리뷰 수정"

	c.HTML(http.StatusOK, "my-page/edit-review-detail.html", state)
	return
}

// getCorrectOrderMst : orderID와 cipher로 검증하여 올바른 orderMst 확보
func getCorrectOrderMst(orderID string, cipher string, state gin.H) (isSuccess bool, result interface{}) {
	if orderID == "" {
		state["errorMessage"] = "필수 파라미터 누락(order-id)"
		return false, state
	}
	if cipher == "" {
		state["errorMessage"] = "필수 파라미터 누락(cipher)"
		return false, state
	}

	orderNo, err := strconv.Atoi(orderID)
	if err != nil {
		state["errorMessage"] = "order-id가 숫자가 아닙니다."
		return false, state
	}

	wherePhrase := fmt.Sprintf(`OrderNo = %d`, orderNo)
	isSuccess, result = mixin.List("orderMst", wherePhrase, "")
	if !isSuccess {
		slackMsg := fmt.Sprintf("[front]HandleInsertReviewOnce::mixin.List('orderMst', %s, '')::%s", wherePhrase, result.(string))
		utils.SendSlackMessage(utils.SlackChannel, slackMsg)
		state["errorMessage"] = "시스템 에러 발생"
		state["errorCode"] = "W0008"
		return false, state
	}

	orderMst := *result.(*[]models.OrderMst)
	if len(orderMst) != 1 {
		state["errorMessage"] = "잘못된 order-id 입력"
		return false, state
	}

	userNo := orderMst[0].UserNo
	reviewSecret := envutil.GetGoDotEnvVariable("REVIEW_SECRET")
	phrase := fmt.Sprintf("%d|%s", userNo, reviewSecret)
	hash := sha256.Sum256([]byte(phrase))
	encoded := hex.EncodeToString(hash[:])
	fmt.Println("encoded", encoded)
	if cipher != encoded {
		state["errorMessage"] = "잘못된 cipher 입력"
		return false, state
	}

	return true, orderMst
}

// GetExposedOrderInfoFromQuery : 쿼리로부터 ExposedOrderInfo 획득
func GetExposedOrderInfoFromQuery(userNo int, orderNo int) (isSuccess bool, result interface{}) {
	var exposedOrderInfo []ExposedOrderInfo
	sqlQuery := fmt.Sprintf(`
		SELECT A.OrderNo, A.CateType AS CateType, A.OrderRound, PaymentDate, A.RcvName 
		FROM
			(SELECT OrderNo, AddressNo, CateType, OrderRound, RcvName FROM OrderMst
			WHERE UserNo = %d AND OrderNo = %d
			AND (StatusCode = "in-delivery" OR StatusCode = "done")) AS A
		JOIN (SELECT OrderNo, RegDate AS PaymentDate FROM OrderPayment) AS D
		ON A.OrderNo = D.OrderNo
	;`, userNo, orderNo)
	fmt.Println(sqlQuery)

	err := database.DB.Select(&exposedOrderInfo, sqlQuery)
	if err != nil {
		return false, err.Error()
	}

	return true, exposedOrderInfo
}

// HandleInsertReviewOnce : review/insert/once 뷰
func HandleInsertReviewOnce(c *gin.Context) {
	state := gin.H{}
	utils.CombineCustomStateGlobalState(c, &state)

	// 인자 유효성 체크
	orderID := c.Query("order-id")
	cipher := c.Query("cipher")
	isSuccess, result := getCorrectOrderMst(orderID, cipher, state)

	// 최초 리뷰 등록 시 문제가 발생하면 메인으로 리다이렉트
	// 하나라도 실패하면 무조건 리다이렉트
	if !isSuccess {
		c.Redirect(http.StatusTemporaryRedirect, "/web/main")
		return
	}
	orderMst := result.([]models.OrderMst)

	// exposedOrderInfo 확보
	isSuccess, result = GetExposedOrderInfoFromQuery(orderMst[0].UserNo, orderMst[0].OrderNo)
	if !isSuccess {
		c.Redirect(http.StatusTemporaryRedirect, "/web/main")
		return
	}
	exposedOrderInfo := result.([]ExposedOrderInfo)
	if len(exposedOrderInfo) < 1 {
		c.Redirect(http.StatusTemporaryRedirect, "/web/main")
		return
	}

	// 주문일
	exposedOrderInfo[0].PaymentDate = strings.Split(exposedOrderInfo[0].PaymentDate, "T")[0]

	// 카테고리 타입 문자열
	isSuccess, cateTypeLabels := controllers.MakeCateTypeLabels(exposedOrderInfo[0].CateType, "|")
	if !isSuccess {
		slackMsg := fmt.Sprintf("[front]HandleMyPageReviewList::controllers.MakeCateTypeLabels(%s, '|')::%s",
			exposedOrderInfo[0].CateType, cateTypeLabels.(string))
		utils.SendSlackMessage(utils.SlackChannel, slackMsg)
		state["errorMessage"] = "시스템 에러 발생"
		state["errorCode"] = "W0015"

		c.Redirect(http.StatusTemporaryRedirect, "/web/main")
		return
	}
	for _, row := range cateTypeLabels.([]string) {
		exposedOrderInfo[0].CateTypeLabels = append(exposedOrderInfo[0].CateTypeLabels, row)
	}

	// 리뷰 할 item 획득, itemScore는 nil 임
	isSuccess, result = GetReviewItems(orderMst[0].UserNo, orderMst[0].OrderNo)
	if !isSuccess {
		c.Redirect(http.StatusTemporaryRedirect, "/web/main")
		return
	}
	reviewItems := result.([]ReviewItems)

	state["reviewItems"] = reviewItems
	state["orderInfo"] = exposedOrderInfo[0]
	state["userNo"] = orderMst[0].UserNo

	state["isDuplicated"] = "yes"
	fmt.Println("reviewItems", reviewItems)
	if len(reviewItems) >= 1 && reviewItems[0].ReviewNo == nil {
		state["isDuplicated"] = "no"
	}
	state["navTitle"] = "리뷰 작성"

	// 브라우저 캐시 금지
	c.Header("Cache-Control", "no-cache, no-store, must-revalidate")

	c.HTML(http.StatusOK, "my-page/insert-review-detail.html", state)
}
