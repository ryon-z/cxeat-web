package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
	"yelloment-api/config"
	envutil "yelloment-api/env_util"
	auth "yelloment-api/middleware/auth"
	"yelloment-api/middleware/mixin"
	"yelloment-api/models"
	"yelloment-api/utils"

	jwt "github.com/appleboy/gin-jwt/v2"
	gin "github.com/gin-gonic/gin"
)

var identityKey = envutil.GetGoDotEnvVariable("IDENTITY_KEY")

// GetAuthMiddleware : authMiddleware 획득
func GetAuthMiddleware() (*jwt.GinJWTMiddleware, error) {
	// the jwt middleware
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:            "test zone",
		Key:              []byte(envutil.GetGoDotEnvVariable("AUTH_SECRET_KEY")),
		Timeout:          time.Hour * 24,
		MaxRefresh:       time.Hour * 24,
		IdentityKey:      identityKey,
		SigningAlgorithm: "HS256",
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*models.UserMst); ok {
				fmt.Println("auth.go, v.UserNo", v.UserNo)
				isNew := false
				if v.StatusCode == "first" {
					isNew = true
				}

				return jwt.MapClaims{
					identityKey: v.UserNo,
					"isNew":     isNew,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			fmt.Println("In IdentityHandler")
			claims := jwt.ExtractClaims(c)
			return &models.UserMst{
				UserNo: int(claims[identityKey].(float64)),
			}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			fmt.Println("In Authenticator")
			var isSuccess bool
			var result interface{}
			switch c.FullPath() {
			case "/user/kakao":
				fmt.Println("come in kakao")
				isSuccess, result = auth.GetUserInfoUsingKakao(c)
			case "/user/naver":
				fmt.Println("come in naver")
				isSuccess, result = auth.GetUserInfoUsingNaver(c)
			case "/login":
				fmt.Println("come in basic auth")
				isSuccess, result = auth.GetUserInfoUsingBasicAuth(c)
			default:
				panic(fmt.Sprintf("auth c.FullPath is not allowed, c.FullPath: %s", c.FullPath()))
			}

			if isSuccess {
				// 획득한 유저 정보
				userMst := result.(*models.UserMst)

				// 마지막 로그인 일자 업데이트
				var lastLoginDateUser models.UserMst
				loc, _ := time.LoadLocation(envutil.GetGoDotEnvVariable("TIMEZONE"))
				now := time.Now().In(loc)

				lastLoginDateString := fmt.Sprintf(
					"%04d-%02d-%02d %02d:%02d:%02d",
					now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(),
				)
				lastLoginDateUser.LastLoginDate = &lastLoginDateString
				wherePhrase := fmt.Sprintf("UserNo = %d", userMst.UserNo)
				updateSuccess, message := mixin.PartialUpdate("userMst", wherePhrase, &lastLoginDateUser, []string{}, []string{"BirthDay", "IsMktAgree"})
				if !updateSuccess {
					slackMsg := fmt.Sprintf("[front]GetAuthMiddleware::mixin.PartialUpdate::%s", message)
					utils.SendSlackMessage(utils.SlackChannel, slackMsg)
					_, errString := config.StringifyError(config.ErrorPrototype{YelloCode: 2003, Message: "updating lastLoginDate is failed"})
					return nil, errors.New(errString)
				}

				return userMst, nil
			}

			_, errString := config.StringifyError(result.(config.ErrorPrototype))
			return nil, errors.New(errString)
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			fmt.Println("In Authorizator")
			fmt.Println("c.Request.Header", c.Request.Header)
			result := true

			return result
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			fmt.Println("In Unauthorized")
			fmt.Println("code", code)
			fmt.Println("message", message)
			fmt.Println("c.Request.RequestURI", c.Request.RequestURI)

			apiAccess := true
			splited := strings.SplitN(c.Request.RequestURI, "/", 3)
			firstPath := splited[1]
			if strings.ToLower(firstPath) == "web" {
				apiAccess = false
			}

			// 로그인 실패 시 로그인 시 발생한 에러가 아닐 경우 로그인 페이지로 리다이렉트
			if !apiAccess && code == 401 && envutil.GetGoDotEnvVariable("USE_TEMPLATE_VIEW") == "TRUE" {
				expireSecString := envutil.GetGoDotEnvVariable("AUTH_TOKEN_EXPIRE_SEC")
				expireSec, err := strconv.Atoi(expireSecString)
				if err != nil {
					c.JSON(code, gin.H{
						"code":    code,
						"expire":  "2006-01-02T15:04:05Z",
						"token":   "",
						"message": err.Error(),
					})
					return
				}

				c.SetCookie("pageAfterLogin", c.Request.RequestURI, expireSec, "/", "", false, false)
				c.SetCookie("jwt", "", -1, "/", "", false, true)
				c.SetCookie("expire", "", -1, "/", "", false, true)
				c.Redirect(http.StatusTemporaryRedirect, "/web/login")
				return
			}

			// 그 외 경우
			fmt.Println("c.Request.Header", c.Request.Header)
			fmt.Println("c.Request.Cookies()", c.Request.Cookies())
			cookie, _ := c.Cookie("jwt")
			fmt.Println("cookie", cookie)
			c.JSON(code, gin.H{
				"code":    code,
				"expire":  "2006-01-02T15:04:05Z",
				"token":   "",
				"message": message,
			})
		},
		// TokenLookup is a string in the form of "<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>"
		// - "query:<name>"
		// - "cookie:<name>"
		// - "param:<name>"
		// TokenLookup: "header: Authorization, query: token, cookie: jwt",
		// TokenLookup: "header: Authorization, query: token, cookie: jwt",
		// TokenLookup: "header: Authorization",
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",
		TokenLookup: "cookie:jwt",

		// TokenHeadName is a string in the header. Default value is "Bearer"
		TokenHeadName: "Bearer",

		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,
	})

	return authMiddleware, err
}
