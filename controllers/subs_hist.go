package controllers

import (
	"net/http"
	"yelloment-api/middleware/mixin"
	"yelloment-api/models"
	"yelloment-api/utils"

	"github.com/gin-gonic/gin"
)

// CreateSubsHist : 주문 기록 생성
func CreateSubsHist(c *gin.Context) {
	userNo := utils.GetUserNo(c)

	var subsHist models.SubsHist
	if err := c.ShouldBindJSON(&subsHist); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	isSuccess, message := CheckSubsNoValid(userNo, subsHist.SubsNo)
	if !isSuccess {
		c.JSON(http.StatusBadRequest, gin.H{"error": message})
		return
	}

	isSuccess, result := mixin.Create("subsHist", &subsHist, []string{"ExecUser", "ExecDate"}, true, []string{})
	utils.GetHTTPResponse(isSuccess, result, c)
}
