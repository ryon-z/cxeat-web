package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"
	envutil "yelloment-api/env_util"
	"yelloment-api/utils"

	"github.com/gin-gonic/gin"
)

// AccessTokenResp : accessToken response
type AccessTokenResp struct {
	Code    int       `json:"code"`
	Expire  time.Time `json:"expire"`
	Token   string    `json:"token"`
	Message string    `json:"message"`
}

// Refresh : AccessToken 갱신
func Refresh(c *gin.Context) {
	state := gin.H{}
	utils.CombineCustomStateGlobalState(c, &state)

	url := "/auth/refresh"
	isSuccess, reqResult := utils.CustomRequest(url, "GET", [][]string{}, [][]string{}, c.Request.Cookies())
	if !isSuccess {
		c.JSON(http.StatusBadRequest, gin.H{"error": reqResult})
		return
	}

	var updatedTokenResp AccessTokenResp
	if err := json.Unmarshal([]byte(reqResult), &updatedTokenResp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	expireSecString := envutil.GetGoDotEnvVariable("AUTH_TOKEN_EXPIRE_SEC")
	expireSec, err := strconv.Atoi(expireSecString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.SetCookie("jwt", updatedTokenResp.Token, expireSec, "/", "", false, true)
	c.SetCookie("expire", updatedTokenResp.Expire.String(), expireSec, "/", "", false, false)

	c.JSON(http.StatusOK, gin.H{"data": "success"})
	return
}
