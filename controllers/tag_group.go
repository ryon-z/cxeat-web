package controllers

import (
	"net/http"
	"yelloment-api/middleware/mixin"
	"yelloment-api/models"
	"yelloment-api/utils"

	"github.com/gin-gonic/gin"
)

// CreateOwnedTagGroup : 소유자 태그 그룹 생성
func CreateOwnedTagGroup(c *gin.Context) {
	userNo := utils.GetUserNo(c)

	var tagGroup models.TagGroup
	if err := c.ShouldBindJSON(&tagGroup); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tagGroup.UserNo = userNo
	isSuccess, result := mixin.CreateOwned("tagGroup", userNo, &tagGroup, []string{"RegDate"}, true, []string{"TagGroupType"})
	utils.GetHTTPResponse(isSuccess, result, c)
}
