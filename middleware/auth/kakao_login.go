package auth

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
	"yelloment-api/config"
	envutil "yelloment-api/env_util"
	"yelloment-api/middleware/mixin"
	"yelloment-api/models"
	"yelloment-api/utils"

	"github.com/gin-gonic/gin"
)

// KakaoAuthTokenResp : kakao auth token resp
type KakaoAuthTokenResp struct {
	Error                 string `json:"error"`
	ErrorDescription      string `json:"error_description"`
	AccessToken           string `json:"access_token"`
	TokenType             string `json:"token_type"`
	RefreshToken          string `json:"refresh_token"`
	ExpiresIn             int    `json:"expires_in"`
	RefreshTokenExpiresIn int    `json:"refresh_token_expires_in"`
}

// KakaoUserInfoResp : kakao user info resp
type KakaoUserInfoResp struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
	ID               int    `json:"id"`
	ConnectedAt      string `json:"connected_at"`
	Properties       struct {
		Nickname string `json:"nickname"`
	} `json:"properties"`
	KakaoAccount struct {
		ProfileNeedsAgreement bool `json:"profile_needs_agreement"`
		Profile               struct {
			Nickname string `json:"nickname"`
		} `json:"profile"`
		HasEmail                  bool   `json:"has_email"`
		EmailNeedsAgreement       bool   `json:"email_needs_agreement"`
		IsEmailValid              bool   `json:"is_email_valid"`
		IsEmailVerified           bool   `json:"is_email_verified"`
		Email                     string `json:"email"`
		HasPhoneNumber            bool   `json:"has_phone_number"`
		PhoneNumberNeedsAgreement bool   `json:"phone_number_needs_agreement"`
		PhoneNumber               string `json:"phone_number"`
		HasAgeRange               bool   `json:"has_age_range"`
		AgeRangeNeedsAgreement    bool   `json:"age_range_needs_agreement"`
		AgeRange                  string `json:"age_range"`
		HasBirthyear              bool   `json:"has_birthyear"`
		BirthyearNeedsAgreement   bool   `json:"birthyear_needs_agreement"`
		Birthyear                 string `json:"birthyear"`
		HasBirthday               bool   `json:"has_birthday"`
		BirthdayNeedsAgreement    bool   `json:"birthday_needs_agreement"`
		Birthday                  string `json:"birthday"`
		BirthdayType              string `json:"birthday_type"`
		HasGender                 bool   `json:"has_gender"`
		GenderNeedsAgreement      bool   `json:"gender_needs_agreement"`
		Gender                    string `json:"gender"`
	} `json:"kakao_account"`
}

func getKakaoAccessToken(c *gin.Context, code string) (bool, interface{}) {
	baseURL := "https://kauth.kakao.com/oauth/token"
	origin := envutil.GetGoDotEnvVariable("PROTOCOL") + "://" + c.Request.Host
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("client_id", envutil.GetGoDotEnvVariable("KAKAO_CLIENT_ID"))
	data.Set("redirect_uri", origin+"/web/login/kakao")
	data.Set("code", code)
	fmt.Println("code", code)

	client := &http.Client{}
	req, err := http.NewRequest("POST", baseURL, strings.NewReader(data.Encode()))
	if err != nil {
		errorResp := config.ErrorPrototype{YelloCode: 9001, Message: err.Error()}
		return false, errorResp
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	req.Header.Add("Host", "kauth.kakao.com")

	// Auth Token 획득
	resp, err := client.Do(req)
	if err != nil {
		return false, config.ErrorPrototype{YelloCode: 9001, Message: err.Error()}
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, config.ErrorPrototype{YelloCode: 9001, Message: err.Error()}
	}
	log.Println(string(body))

	var kakaoAuthTokenResp KakaoAuthTokenResp
	err = json.Unmarshal(body, &kakaoAuthTokenResp)
	if err != nil {
		return false, config.ErrorPrototype{YelloCode: 9002, Message: err.Error()}
	}

	if kakaoAuthTokenResp.Error == "" {
		return true, kakaoAuthTokenResp.AccessToken
	}

	return false, config.ErrorPrototype{YelloCode: 9002, Message: kakaoAuthTokenResp.ErrorDescription}
}

func getKakaoUserInfo(accessToken string) (bool, interface{}) {
	baseURL := "https://kapi.kakao.com/v2/user/me"

	client := &http.Client{}
	req, err := http.NewRequest("POST", baseURL, nil)
	if err != nil {
		return false, config.ErrorPrototype{YelloCode: 9002, Message: err.Error()}
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded;charset=utf-8")
	req.Header.Add("Host", "kauth.kakao.com")
	req.Header.Add("Authorization", "Bearer "+accessToken)

	// Auth Token 획득
	resp, err := client.Do(req)
	if err != nil {
		return false, config.ErrorPrototype{YelloCode: 9002, Message: err.Error()}
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, config.ErrorPrototype{YelloCode: 9002, Message: err.Error()}
	}
	log.Println(string(body))

	var kakaoUserInfoResp KakaoUserInfoResp
	err = json.Unmarshal(body, &kakaoUserInfoResp)
	if err != nil {
		return false, config.ErrorPrototype{YelloCode: 9002, Message: err.Error()}
	}

	if kakaoUserInfoResp.Error == "" {
		return true, kakaoUserInfoResp
	}

	return false, config.ErrorPrototype{YelloCode: 9002, Message: kakaoUserInfoResp.ErrorDescription}
}

// getConvertedKakaoPhone : 카카오 핸드폰 규격에 맞게 변환
func getConvertedKakaoPhone(phoneNumber string) (bool, string) {
	splited := strings.Split(phoneNumber, " ")
	if len(splited) != 2 {
		return false, "kakao phone number is not valid"
	}

	if splited[0] != "+82" {
		return false, "kakao phone number is not korean phone number"
	}

	converted := "0" + strings.ReplaceAll(splited[1], "-", "")
	if len(converted) != 11 && len(converted) != 10 {
		return false, "kakao phone number is not valid"
	}

	return true, converted
}

// getConvertedBirthDay : 카카오 생년월일 규격에 맞게 변환
func getConvertedBirthDay(ageRange string, birthYear string, birthDay string) (bool, string) {
	// 생년 또는 나이범주 체크
	if birthYear == "" && ageRange == "" {
		return false, ""
	}

	// 생일 체크
	if birthDay == "" {
		return false, ""
	}
	_, birthDayErr := strconv.Atoi(birthDay)
	if birthDayErr != nil {
		return false, ""
	}

	// 생년 추출
	var year int
	if birthYear != "" {
		birthYearInt, birthYearErr := strconv.Atoi(birthYear)
		if birthYearErr != nil {
			return false, ""
		}

		year = birthYearInt
	} else {
		splited := strings.Split(ageRange, "~")
		minAge, minAgeErr := strconv.Atoi(splited[0])
		if minAgeErr != nil {
			return false, ""
		}
		maxAge, maxAgeErr := strconv.Atoi(splited[1])
		if maxAgeErr != nil {
			return false, ""
		}
		meanAge := (minAge + maxAge) / 2.0

		loc, _ := time.LoadLocation(envutil.GetGoDotEnvVariable("TIMEZONE"))
		now := time.Now().In(loc)
		year = now.Year() - meanAge
	}

	return true, fmt.Sprintf("%d-%s-%s", year, birthDay[:2], birthDay[2:])

}

// createNewUser : 새 유저 등록
func createNewUser(userMst models.UserMst) (isSuccess bool, result interface{}) {
	var lastIndex interface{}
	userMst.StatusCode = "first"
	isSuccess, lastIndex = mixin.Create("userMst", &userMst, []string{"UserNo", "RegDate"}, true, []string{})
	if !isSuccess {
		return false, config.ErrorPrototype{YelloCode: 2001, Message: lastIndex.(string)}
	}
	userMst.UserNo = int(lastIndex.(int64))

	// 신규 회원가입 알림
	utils.SendAlrimTalk(
		userMst.UserPhone,
		"bizp_2021083018045626143856020",
		fmt.Sprintf(`%s님, 큐잇 회원 가입을 환영합니다!

큐잇은 식생활을 큐레이션 하는 농산물 큐레이션 정기구독서비스입니다.
		
큐잇은 데이터에 기반한 개인 라이프스타일과 취향을 반영하여 합리적이고 편리하게 농산물을 이용할 수 있도록 식생활을 연구합니다.
		
		
몇 가지 질문을 통해 고객님의 농산물 소비패턴과 건강정보를 확인하시고, 정기구독으로 다양한 혜택을 받으세요!`, userMst.UserName),
		[]utils.AlrimTalkButton{
			{
				BtnName:   "큐잇 정기구독 시작하기",
				BtnType:   "WL",
				UrlPC:     "https://cueat.kr/",
				UrlMobile: "https://cueat.kr/",
			},
		},
	)

	return true, &userMst
}

// GetUserInfoUsingKakao : 카카오 로그인을 이용하여 유저 정보 획득
func GetUserInfoUsingKakao(c *gin.Context) (bool, interface{}) {
	code := c.Query("code")

	isSuccess, result := getKakaoAccessToken(c, code)
	if !isSuccess {
		return false, result.(config.ErrorPrototype)
	}
	accessToken := result.(string)
	isSuccess, result = getKakaoUserInfo(accessToken)
	if !isSuccess {
		return false, result.(config.ErrorPrototype)
	}

	userInfo := result.(KakaoUserInfoResp)
	kakaoID := userInfo.ID

	// TEMPORARY TESTUSER
	if kakaoID == 1671613204 {
		isSuccess, result = mixin.List("userMst", fmt.Sprintf(`UserID = "%d"`, kakaoID), "")
		if !isSuccess {
			return false, config.ErrorPrototype{YelloCode: 2001, Message: result.(string)}
		}
		testUserMst := *result.(*[]models.UserMst)
		return true, &testUserMst[0]
	}

	// 닉네임에 포함된 이모지 치환
	nickName := userInfo.Properties.Nickname
	var emojiRx = regexp.MustCompile(`[\x{1F004}-\x{1F9E6}]|[\x{1F600}-\x{1F9D0}]`)
	nickName = emojiRx.ReplaceAllString(nickName, "(emoji)")

	email := ""
	if userInfo.KakaoAccount.HasEmail &&
		userInfo.KakaoAccount.IsEmailValid &&
		userInfo.KakaoAccount.IsEmailVerified {
		email = userInfo.KakaoAccount.Email
	}
	phone := ""
	if userInfo.KakaoAccount.HasPhoneNumber {
		phone = userInfo.KakaoAccount.PhoneNumber
	} else {
		return false, config.ErrorPrototype{YelloCode: 9004, Message: "Not exists phone number in kakao account info"}
	}
	isSuccess, phoneResult := getConvertedKakaoPhone(phone)
	if !isSuccess {
		return false, config.ErrorPrototype{YelloCode: 9006, Message: phoneResult}
	}
	convertedPhone := phoneResult

	gender := ""
	if userInfo.KakaoAccount.HasGender {
		gender = userInfo.KakaoAccount.Gender
	}

	birthYear := ""
	if userInfo.KakaoAccount.HasBirthyear {
		birthYear = userInfo.KakaoAccount.Birthyear
	}

	birthDay := ""
	if userInfo.KakaoAccount.HasBirthday {
		birthDay = userInfo.KakaoAccount.Birthday
	}

	ageRange := ""
	if userInfo.KakaoAccount.HasAgeRange {
		ageRange = userInfo.KakaoAccount.AgeRange
	}

	isSuccess, convertedBirthDay := getConvertedBirthDay(ageRange, birthYear, birthDay)

	isSuccess, result = mixin.List("userMst", fmt.Sprintf(`UserID = "%d"`, kakaoID), "")
	if !isSuccess {
		return false, config.ErrorPrototype{YelloCode: 2001, Message: result.(string)}
	}

	var userMst models.UserMst
	userMst.UserID = strconv.Itoa(kakaoID)
	userMst.UserType = "kakao"
	userMst.UserEmail = &email
	userMst.UserPhone = convertedPhone
	userMst.UserName = nickName
	userMst.UserGender = &gender
	userMst.UserSecretKey = "empty"
	if isSuccess {
		userMst.BirthDay = &convertedBirthDay
	}

	if len(*result.(*[]models.UserMst)) == 0 {
		return createNewUser(userMst)
	}

	if len(*result.(*[]models.UserMst)) != 1 {
		return false, config.ErrorPrototype{YelloCode: 2002, Message: "kakaoID is dupicated"}
	}

	// DB에 등록된 회원정보
	remainedUserMsts := *result.(*[]models.UserMst)
	fmt.Println(remainedUserMsts)

	// 현재 시각 획득
	if remainedUserMsts[0].StatusCode == "leave" {
		loc, _ := time.LoadLocation(envutil.GetGoDotEnvVariable("TIMEZONE"))
		now := time.Now().In(loc)

		if remainedUserMsts[0].LeaveDate == nil {
			return false, config.ErrorPrototype{YelloCode: 5000, Message: "leaveDate not exists"}
		}
		leaveDate, err := time.Parse(time.RFC3339, *(remainedUserMsts[0].LeaveDate))
		if err != nil {
			return false, config.ErrorPrototype{YelloCode: 5002, Message: "parsing leavedate is failed"}
		}

		// 현재 날짜가 탈퇴 후 3개월 이내면 재가입 금지
		nextSignupDate := leaveDate.AddDate(0, 3, 0)
		if now.Before(nextSignupDate) {
			return false, config.ErrorPrototype{
				YelloCode: 5001,
				Message:   fmt.Sprintf("%04d년 %02d월 %02d일", nextSignupDate.Year(), nextSignupDate.Month(), nextSignupDate.Day()),
			}
		}

		// 폐기된 유저 행 제거 후 새 유저 행 삽입
		var removedUser models.UserMst
		removedUser.UserID = "removed_" + remainedUserMsts[0].UserID
		wherePhrase := fmt.Sprintf("UserNo = %d", remainedUserMsts[0].UserNo)
		isSuccess, message := mixin.PartialUpdate("userMst", wherePhrase, &removedUser, []string{}, []string{"BirthDay", "IsMktAgree"})
		if !isSuccess {
			slackMsg := fmt.Sprintf("[front]GetUserInfoUsingKakao::mixin.PartialUpdate::%s", message)
			utils.SendSlackMessage(utils.SlackChannel, slackMsg)
			return false, config.ErrorPrototype{YelloCode: 2003, Message: "changing userID is failed"}
		}

		return createNewUser(userMst)
	}

	return true, &remainedUserMsts[0]
}
