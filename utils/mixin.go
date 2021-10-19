package utils

import (
	envutil "yelloment-api/env_util"
	"yelloment-api/models"

	"github.com/gin-gonic/gin"
)

// GetUserNo : 유저 번호 리턴
func GetUserNo(c *gin.Context) int {
	identityKey := envutil.GetGoDotEnvVariable("IDENTITY_KEY")
	user, _ := c.Get(identityKey)
	userNo := user.(*models.UserMst).UserNo

	return userNo
}
