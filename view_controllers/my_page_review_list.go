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

// ReviewsWithOrder : OrderMst와 ReviewMst를 LEFT JOIN
type ReviewsWithOrder struct {
	ReviewNo       int     `db:"ReviewNo"`
	ReviewDesc     *string `db:"ReviewDesc"`
	MeanScore      float64 `db:"MeanScore"`
	OrderNo        int     `db:"OrderNo"`
	CateType       string  `db:"CateType"`
	OrderRound     *int    `db:"OrderRound"`
	PaymentDate    string  `db:"PaymentDate"`
	RcvName        string  `db:"RcvName"`
	CateTypeLabels []string
}

// GetReviewsWithOrder : 주문 목록 표시
func GetReviewsWithOrder(userNo int) (isSuccess bool, result interface{}) {
	var reviewsWithOrder []ReviewsWithOrder
	sqlQuery := fmt.Sprintf(`
		SELECT B.ReviewNo, B.ReviewDesc, B.MeanScore, A.OrderNo, A.CateType, 
			A.OrderRound, PaymentDate, A.RcvName 
		FROM
			(SELECT OrderNo, AddressNo, CateType, OrderRound, RcvName FROM OrderMst
			WHERE UserNo = %d
			AND (StatusCode = "in-delivery" OR StatusCode = "done")) AS A
		JOIN (SELECT ReviewNo, OrderNo, ReviewDesc,
			(SELECT ROUND(AVG(ReviewScore), 1) FROM ItemReview AS ir WHERE ir.ReviewNo = rm.ReviewNo) AS MeanScore
			FROM ReviewMst AS rm
			) AS B
		ON A.OrderNo = B.OrderNo
		JOIN (SELECT OrderNo, RegDate AS PaymentDate FROM OrderPayment) AS D
		ON A.OrderNo = D.OrderNo
		ORDER BY PaymentDate DESC
	;`, userNo)
	fmt.Println(sqlQuery)

	err := database.DB.Select(&reviewsWithOrder, sqlQuery)
	if err != nil {
		return false, err.Error()
	}

	return true, reviewsWithOrder
}

// HandleMyPageReviewList : my-page/review/list 뷰
func HandleMyPageReviewList(c *gin.Context) {
	state := gin.H{}
	utils.CombineCustomStateGlobalState(c, &state)

	userNo := utils.GetUserNo(c)

	// OrderMst와 ReviewMst를 LEFT JOIN
	isSuccess, result := GetReviewsWithOrder(userNo)
	if !isSuccess {
		c.JSON(http.StatusBadRequest, gin.H{"error": result.(string)})
		return
	}

	// ReviewNo가 없는 OrderNo만 추출하여 각각 models.Review 객체에 담음
	reviewsWithOrder := result.([]ReviewsWithOrder)
	for i, review := range reviewsWithOrder {
		// 주문일
		reviewsWithOrder[i].PaymentDate = strings.Split(review.PaymentDate, "T")[0]

		// 카테고리 타입 문자열
		isSuccess, cateTypeLabels := controllers.MakeCateTypeLabels(review.CateType, "|")
		if !isSuccess {
			slackMsg := fmt.Sprintf("[front]HandleMyPageReviewList::controllers.MakeCateTypeLabels(%s, '|')::%s",
				review.CateType, cateTypeLabels)
			utils.SendSlackMessage(utils.SlackChannel, slackMsg)
			state["errorMessage"] = "시스템 에러 발생"
			state["errorCode"] = "W0015"
			c.HTML(http.StatusBadRequest, "error.html", state)
			return
		}
		for _, row := range cateTypeLabels.([]string) {
			reviewsWithOrder[i].CateTypeLabels = append(reviewsWithOrder[i].CateTypeLabels, row)
		}
	}
	state["reviewsWithOrder"] = reviewsWithOrder
	state["navTitle"] = "나의 리뷰"

	c.HTML(http.StatusOK, "my-page/review-list.html", state)
	return
}
