package auth

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"yelloment-api/config"
	envutil "yelloment-api/env_util"
	"yelloment-api/middleware/mixin"
	"yelloment-api/models"

	"github.com/gin-gonic/gin"
)

// NaverAuthTokenResp : naver auth token resp
type NaverAuthTokenResp struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
	AccessToken      string `json:"access_token"`
	TokenType        string `json:"token_type"`
	RefreshToken     string `json:"refresh_token"`
	ExpiresIn        string `json:"expires_in"`
}

// NaverUserInfoResp : naver user info resp
type NaverUserInfoResp struct {
	Code       int    `json:"code"`
	Msg        string `json:"msg"`
	Resultcode string `json:"resultcode"`
	Message    string `json:"message"`
	Response   struct {
		ID         string `json:"id"`
		Age        string `json:"age"`
		Gender     string `json:"gender"`
		Email      string `json:"email"`
		Mobile     string `json:"mobile"`
		MobileE164 string `json:"mobile_e164"`
		Name       string `json:"name"`
		Birthday   string `json:"birthday"`
	} `json:"response"`
}

func getNaverAccessToken(code string) (bool, interface{}) {
	baseURL := "https://nid.naver.com/oauth2.0/token"
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("client_id", envutil.GetGoDotEnvVariable("NAVER_CLIENT_ID"))
	data.Set("client_secret", envutil.GetGoDotEnvVariable("NAVER_CLIENT_SECRET"))
	data.Set("code", code)
	data.Set("state", "개발중")
	fmt.Println("code", code)

	client := &http.Client{}
	req, err := http.NewRequest("POST", baseURL, strings.NewReader(data.Encode()))
	if err != nil {
		return false, config.ErrorPrototype{YelloCode: 9011, Message: err.Error()}
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	// Auth Token 획득
	resp, err := client.Do(req)
	if err != nil {
		return false, config.ErrorPrototype{YelloCode: 9011, Message: err.Error()}
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, config.ErrorPrototype{YelloCode: 9011, Message: err.Error()}
	}
	log.Println(string(body))

	var naverAuthTokenResp NaverAuthTokenResp
	err = json.Unmarshal(body, &naverAuthTokenResp)
	if err != nil {
		return false, config.ErrorPrototype{YelloCode: 9012, Message: err.Error()}
	}

	if naverAuthTokenResp.Error == "" {
		return true, naverAuthTokenResp.AccessToken
	}

	return false, config.ErrorPrototype{YelloCode: 9012, Message: naverAuthTokenResp.ErrorDescription}
}

func getNaverUserInfo(accessToken string) (bool, interface{}) {
	baseURL := "https://openapi.naver.com/v1/nid/me"

	client := &http.Client{}
	req, err := http.NewRequest("GET", baseURL, nil)
	if err != nil {
		return false, config.ErrorPrototype{YelloCode: 9013, Message: err.Error()}
	}
	req.Header.Add("Authorization", "Bearer "+accessToken)
	fmt.Println("accessToken", accessToken)

	// Auth Token 획득
	resp, err := client.Do(req)
	if err != nil {
		return false, config.ErrorPrototype{YelloCode: 9013, Message: err.Error()}
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, config.ErrorPrototype{YelloCode: 9013, Message: err.Error()}
	}
	log.Println(string(body))

	var naverUserInfoResp NaverUserInfoResp
	err = json.Unmarshal(body, &naverUserInfoResp)
	if err != nil {
		return false, config.ErrorPrototype{YelloCode: 9013, Message: err.Error()}
	}

	if naverUserInfoResp.Resultcode == "00" {
		return true, naverUserInfoResp
	}

	return false, config.ErrorPrototype{YelloCode: 9013, Message: naverUserInfoResp.Msg}
}

// getConvertedNaverPhone : 네이버 핸드폰 규격에 맞게 변환
func getConvertedNaverPhone(phoneNumber string) (bool, string) {
	converted := strings.ReplaceAll(phoneNumber, "-", "")
	if len(converted) != 11 && len(converted) != 10 {
		return false, "naver phone number is not valid"
	}

	return true, converted
}

// GetUserInfoUsingNaver : 네이버 로그인을 이용하여 유저 정보 획득
func GetUserInfoUsingNaver(c *gin.Context) (bool, interface{}) {
	code := c.Query("code")

	isSuccess, result := getNaverAccessToken(code)
	if !isSuccess {
		return false, result.(config.ErrorPrototype)
	}

	accessToken := result.(string)
	isSuccess, result = getNaverUserInfo(accessToken)
	if !isSuccess {
		return false, result.(config.ErrorPrototype)
	}

	userInfo := result.(NaverUserInfoResp)
	fmt.Println("userInfo", userInfo)
	naverID := userInfo.Response.ID
	email := userInfo.Response.Email
	mobile := userInfo.Response.Mobile
	gender := userInfo.Response.Gender
	name := userInfo.Response.Name
	if mobile == "" {
		return false, config.ErrorPrototype{YelloCode: 9014, Message: "Not exists phone number in naver account info"}
	}
	isSuccess, phoneResult := getConvertedNaverPhone(mobile)
	if !isSuccess {
		return false, config.ErrorPrototype{YelloCode: 9016, Message: phoneResult}
	}
	convertedPhone := phoneResult

	isSuccess, result = mixin.List("userMst", fmt.Sprintf(`UserID = "%s"`, naverID), "")
	if !isSuccess {
		return false, config.ErrorPrototype{YelloCode: 2001, Message: result.(string)}
	}

	var userMst models.UserMst
	userMst.UserID = naverID
	userMst.UserType = "naver"
	userMst.UserEmail = &email
	userMst.UserPhone = convertedPhone
	userMst.UserName = name
	userMst.UserGender = &gender
	userMst.UserSecretKey = "empty"

	var lastIndex interface{}
	if len(*result.(*[]models.UserMst)) == 0 {
		isSuccess, lastIndex = mixin.Create("userMst", &userMst, []string{"UserNo", "RegDate", "StatusCode"}, true, []string{})
		if !isSuccess {
			return false, config.ErrorPrototype{YelloCode: 2001, Message: lastIndex.(string)}
		}
		userMst.UserNo = int(lastIndex.(int64))

		return true, &userMst
	}

	if len(*lastIndex.(*[]models.UserMst)) != 1 {
		return false, config.ErrorPrototype{YelloCode: 2002, Message: "naverID is dupicated"}
	}

	remainedUserMsts := *lastIndex.(*[]models.UserMst)
	fmt.Println(remainedUserMsts)
	return true, &remainedUserMsts[0]
}
