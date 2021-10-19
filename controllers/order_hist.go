package controllers

import (
	"net/http"
	"yelloment-api/middleware/mixin"
	"yelloment-api/models"
	"yelloment-api/utils"

	"github.com/gin-gonic/gin"
)

// CreateOrderHist : 주문 기록 생성
func CreateOrderHist(c *gin.Context) {
	userNo := utils.GetUserNo(c)

	var orderHist models.OrderHist
	if err := c.ShouldBindJSON(&orderHist); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	isSuccess, message := CheckOrderNoValid(userNo, orderHist.OrderNo)
	if !isSuccess {
		c.JSON(http.StatusBadRequest, gin.H{"error": message})
		return
	}

	isSuccess, result := mixin.Create("orderHist", &orderHist, []string{"ExecUser", "ExecDate"}, true, []string{})
	utils.GetHTTPResponse(isSuccess, result, c)
}
