package auth

import (
	"fmt"
	"yelloment-api/models"
	"yelloment-api/utils"

	jwt "github.com/appleboy/gin-jwt/v2"
	gin "github.com/gin-gonic/gin"
)

type login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

func getUserPayload(userID string, password string) (bool, interface{}) {
	wherePhrase := fmt.Sprintf(`UserID="%s"`, userID)
	isSuccess, selectResult := models.Select("userMst", wherePhrase, "")
	if !isSuccess {
		return false, selectResult.(string)
	}

	userInfo := *selectResult.(*[]models.UserMst)
	if len(userInfo) == 0 {
		return false, "User is not registered"
	}

	hash := userInfo[0].UserSecretKey
	isSuccess, message := utils.ComparePasswords(hash, password)
	resultUserInfo := models.UserMst{}
	if isSuccess {
		resultUserInfo.UserNo = userInfo[0].UserNo
		resultUserInfo.UserID = userInfo[0].UserID
		resultUserInfo.UserName = userInfo[0].UserName
		return true, &resultUserInfo
	}

	return false, message
}

// GetUserInfoUsingBasicAuth : basic auth를 이용하여 유저 정보 획득
func GetUserInfoUsingBasicAuth(c *gin.Context) (bool, interface{}) {
	userID, password, hasAuth := c.Request.BasicAuth()
	var loginVals login
	if !hasAuth {
		if err := c.ShouldBind(&loginVals); err != nil {
			return false, jwt.ErrMissingLoginValues
		}
		userID = loginVals.Username
		password = loginVals.Password
	}

	fmt.Println("userID", userID)
	fmt.Println("password", password)

	isSuccess, result := getUserPayload(userID, password)
	return isSuccess, result
}
