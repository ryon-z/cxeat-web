package controllers

import (
	"fmt"
	"net/http"
	"yelloment-api/middleware/mixin"
	"yelloment-api/models"
	"yelloment-api/utils"

	"github.com/gin-gonic/gin"
)

// OwnedItemReview : 소유주가 표시된 아이템 리뷰
type OwnedItemReview struct {
	ItemReviews []models.ItemReview `json:"ItemReviews"`
	UserNo      int                 `json:"UserNo"`
}

// ChekcItemReviewsValid : itemReviews 각 행의 orderNo와 reviewNo가 동일한지 검증
func ChekcItemReviewsValid(itemReviews []models.ItemReview) (isSuccess bool, message string) {
	var uniqueOrderNos []int
	var uniqueReviewNos []int
	for _, itemReview := range itemReviews {
		orderNo := itemReview.OrderNo
		reviewNo := itemReview.ReviewNo
		if !utils.IntInSlice(orderNo, uniqueOrderNos) {
			uniqueOrderNos = append(uniqueOrderNos, orderNo)
		}

		if !utils.IntInSlice(reviewNo, uniqueReviewNos) {
			uniqueReviewNos = append(uniqueReviewNos, reviewNo)
		}
	}
	if len(uniqueOrderNos) != 1 {
		return false, "orderNo in itemReviews is not one."
	}
	if len(uniqueReviewNos) != 1 {
		return false, "reviewNo in itemReviews is not one."
	}

	return true, "success"
}

// CreateItemReview : 소유자 아이템 리뷰 생성
func CreateItemReview(c *gin.Context) {
	var ownedItemReview OwnedItemReview
	if err := c.ShouldBindJSON(&ownedItemReview); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	itemReviews := ownedItemReview.ItemReviews
	userNo := ownedItemReview.UserNo

	// itemReviews 각 행의 orderNo와 reviewNo가 동일한지 검증
	isSuccess, message := ChekcItemReviewsValid(itemReviews)
	if !isSuccess {
		c.JSON(http.StatusBadRequest, gin.H{"error": message})
		return
	}

	// orderNo가 userNo의 orderNo인지 확인
	isSuccess, message = CheckOrderNoValid(userNo, itemReviews[0].OrderNo)
	if !isSuccess {
		c.JSON(http.StatusBadRequest, gin.H{"error": message})
		return
	}

	isSuccess, result := mixin.CreateMulti(
		"itemReview", &itemReviews,
		[]string{"RegDate", "UpdDate"}, true, []string{"ReviewDesc", "ReviewScore"})
	utils.GetHTTPResponse(isSuccess, result, c)
}

// PartialUpdateItemReview : 소유자 아이템 리뷰 수정
func PartialUpdateItemReview(c *gin.Context) {
	var ownedItemReview OwnedItemReview
	if err := c.ShouldBindJSON(&ownedItemReview); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	itemReviews := ownedItemReview.ItemReviews
	userNo := ownedItemReview.UserNo

	// itemReviews 각 행의 orderNo와 reviewNo가 동일한지 검증
	isSuccess, message := ChekcItemReviewsValid(itemReviews)
	if !isSuccess {
		c.JSON(http.StatusBadRequest, gin.H{"error": message})
		return
	}

	// orderNo가 userNo의 orderNo인지 확인
	isSuccess, message = CheckOrderNoValid(userNo, itemReviews[0].OrderNo)
	if !isSuccess {
		c.JSON(http.StatusBadRequest, gin.H{"error": message})
		return
	}

	// 일괄 삭제
	wherePhrase := fmt.Sprintf("OrderNo = %d AND ReviewNo = %d", itemReviews[0].OrderNo, itemReviews[0].ReviewNo)
	isSuccess, message = mixin.Delete("itemReview", wherePhrase)
	if !isSuccess {
		c.JSON(http.StatusBadRequest, gin.H{"error": message})
		return
	}

	// 일괄 삽입
	isSuccess, result := mixin.CreateMulti(
		"itemReview", &itemReviews,
		[]string{"RegDate", "UpdDate"}, true, []string{"ReviewDesc", "ReviewScore"})
	utils.GetHTTPResponse(isSuccess, result, c)
}
