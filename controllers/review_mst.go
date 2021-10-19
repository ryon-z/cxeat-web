package controllers

import (
	"fmt"
	"net/http"
	"yelloment-api/middleware/mixin"
	"yelloment-api/models"
	"yelloment-api/utils"

	"github.com/gin-gonic/gin"
)

// OwnedReview : 소유주가 표시된 아이템 리뷰
type OwnedReview struct {
	Review models.ReviewMst `json:"Review"`
	UserNo int              `json:"UserNo"`
}

// CreateReviewMst : 소유자 리뷰 생성
func CreateReviewMst(c *gin.Context) {
	var ownedReview OwnedReview
	if err := c.ShouldBindJSON(&ownedReview); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	reviewMst := ownedReview.Review
	userNo := ownedReview.UserNo

	// 주문 번호 유효성 검사
	isSuccess, message := CheckOrderNoValid(userNo, reviewMst.OrderNo)
	if !isSuccess {
		c.JSON(http.StatusBadRequest, gin.H{"error": message})
		return
	}

	// 중복 리뷰 체크
	wherePhrase := fmt.Sprintf("OrderNo = %d", reviewMst.OrderNo)
	isSuccess, result := mixin.List("reviewMst", wherePhrase, "")
	if !isSuccess {
		c.JSON(http.StatusBadRequest, gin.H{"error": result})
		return
	}
	reviewMstForDuliCheck := *result.(*[]models.ReviewMst)
	fmt.Println("reviewMstForDuliCheck", reviewMstForDuliCheck)
	if len(reviewMstForDuliCheck) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "this review is already registrated"})
		return
	}

	isSuccess, result = mixin.Create("reviewMst", &reviewMst, []string{"RegDate", "UpdDate"}, true, []string{"ReviewDesc"})
	utils.GetHTTPResponse(isSuccess, result, c)
}

// PartialUpdateReviewMst : 소유자 리뷰 수정
func PartialUpdateReviewMst(c *gin.Context) {
	var ownedReview OwnedReview
	if err := c.ShouldBindJSON(&ownedReview); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	reviewMst := ownedReview.Review
	userNo := ownedReview.UserNo

	isSuccess, message := CheckOrderNoValid(userNo, reviewMst.OrderNo)
	if !isSuccess {
		c.JSON(http.StatusBadRequest, gin.H{"error": message})
		return
	}

	wherePhrase := fmt.Sprintf("OrderNo = %d", reviewMst.OrderNo)
	isSuccess, message = mixin.PartialUpdate("reviewMst", wherePhrase, &reviewMst, []string{"RegDate", "UpdDate"}, []string{"ReviewDesc"})
	utils.GetHTTPResponse(isSuccess, message, c)
}

// GetReviewMstUnit : 리뷰 조회
func GetReviewMstUnit(c *gin.Context, reviewNo int) (bool, interface{}) {
	wherePhrase := fmt.Sprintf(`ReviewNo = %d`, reviewNo)
	isSuccess, result := mixin.List("reviewMst", wherePhrase, "RegDate DESC")
	return isSuccess, result
}
