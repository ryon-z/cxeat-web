package controllers

import (
	"net/http"
	"yelloment-api/middleware/mixin"
	"yelloment-api/models"
	"yelloment-api/utils"

	"github.com/gin-gonic/gin"
)

// CreateOwnedTag : 소유자 태그 그룹 생성
func CreateOwnedTag(c *gin.Context) {
	var tag []models.Tag
	if err := c.ShouldBindJSON(&tag); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	isSuccess, result := mixin.CreateMulti(
		"tag", &tag,
		[]string{"RegDate", "RegUser", "UpdDate", "UpdUser"}, true, []string{"TagLabel", "TagValue"})
	utils.GetHTTPResponse(isSuccess, result, c)
}
