package controllers

import (
	"fmt"
	"net/http"
	"regexp"
	"time"
	"yelloment-api/config"
	envutil "yelloment-api/env_util"
	"yelloment-api/middleware/mixin"
	"yelloment-api/models"
	"yelloment-api/utils"

	"github.com/gin-gonic/gin"
)

// checkUserValid : 유저 구조체 유효성 검사
func checkUserValid(user models.UserMst) (bool, string) {
	if user.UserName == "" {
		return false, "UserName is empty"
	}
	if user.UserSecretKey == "" {
		return false, "Password is empty"
	}

	if len([]byte(user.UserSecretKey)) >= 70 {
		return false, "Password is too long"
	}

	if !utils.StringInSlice(user.StatusCode, config.UserStatusCodes) {
		return false, "StatusCode is not allowed"
	}

	return true, "success"
}

// CreateUser : 유저 생성
func CreateUser(c *gin.Context) {
	var user models.UserMst
	if err := c.ShouldBindJSON(&user); err != nil {
		fmt.Println("user", user)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Println("come in here?")

	isSuccess, message := checkUserValid(user)
	if !isSuccess {
		c.JSON(http.StatusBadRequest, gin.H{"error": message})
		return
	}

	// Encryt UserSecretKey
	isSuccess, result := utils.HashAndSalt(user.UserSecretKey)
	if isSuccess {
		user.UserSecretKey = result
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": result})
		return
	}

	user.StatusCode = "normal"
	isSuccess, createResult := mixin.Create("userMst", &user, []string{"UserNo", "RegDate"}, true, []string{})
	utils.GetHTTPResponse(isSuccess, createResult, c)
}

// GetOwnedUser : 소유자 유저 정보 획득
func GetOwnedUser(c *gin.Context) (isSuccess bool, result interface{}) {
	userNo := utils.GetUserNo(c)
	wherePhrase := fmt.Sprintf(`UserNo = %d AND (StatusCode = "normal" OR StatusCode = "first")`, userNo)
	return mixin.List("userMst", wherePhrase, "")
}

// ListOwnedUser : 소유자 유저 정보 리스트 조회
func ListOwnedUser(c *gin.Context) {
	isSuccess, result := GetOwnedUser(c)
	utils.GetHTTPResponse(isSuccess, result, c)
}

// PartialUpdateOwnedUser : 소유자 유저 정보 수정
func PartialUpdateOwnedUser(c *gin.Context) {
	userNo := utils.GetUserNo(c)

	var user models.UserMst
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if user.StatusCode != "" && !utils.StringInSlice(user.StatusCode, config.UserStatusCodes) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "StatusCode is not allowed"})
		return
	}

	var emojiRx = regexp.MustCompile(`[\x{1F004}-\x{1F9E6}]|[\x{1F600}-\x{1F9D0}]`)
	user.UserName = emojiRx.ReplaceAllString(user.UserName, "(emoji)")

	wherePhrase := fmt.Sprintf("UserNo = %d", userNo)
	isSuccess, message := mixin.PartialUpdate("userMst", wherePhrase, &user, []string{}, []string{"BirthDay", "IsMktAgree"})
	utils.GetHTTPResponse(isSuccess, message, c)
}

// WithdrawalOwnedUser : 고객 탈퇴
func WithdrawalOwnedUser(c *gin.Context) {
	userNo := utils.GetUserNo(c)

	var user models.UserMst
	user.UserNo = userNo
	user.UserName = ""
	user.UserPhone = ""
	userEmail := ""
	user.UserEmail = &userEmail
	birthDay := "1111-01-01"
	user.BirthDay = &birthDay
	user.StatusCode = "leave"

	loc, _ := time.LoadLocation(envutil.GetGoDotEnvVariable("TIMEZONE"))
	now := time.Now().In(loc)
	leaveDate := now.Format("2006-01-02 15:04:05")
	user.LeaveDate = &leaveDate

	wherePhrase := fmt.Sprintf("UserNo = %d", userNo)
	isSuccess, message := mixin.PartialUpdate("userMst", wherePhrase, &user, []string{}, []string{"BirthDay", "IsMktAgree"})
	utils.GetHTTPResponse(isSuccess, message, c)
}
