package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"yelloment-api/models"
	"yelloment-api/router"
	"yelloment-api/utils"
)

// Ts : 테스트서버
var Ts *httptest.Server

// Bearer : Bearer 토큰
var Bearer string

// RespForTest : 테스트 응답
type RespForTest struct {
	Error string
	Data  interface{}
}

// TokenResp : 토큰 응답
type TokenResp struct {
	Code   int    `json:"code"`
	Expire string `json:"expire"`
	Token  string `json:"token"`
}

// CheckError : 에러 체크
func CheckError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// CheckCondition : 조건 체크
func CheckCondition(condition bool, errorMessage interface{}) {
	if !condition {
		log.Fatal(errorMessage)
	}
}

// RequestForTest : 테스트 유저로 요청
func RequestForTest(method string, url string, body io.Reader) *http.Response {
	req, err := http.NewRequest(method, url, body)
	CheckError(err)
	req.Header.Set("Authorization", "Bearer "+Bearer)
	req.Header.Add("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	CheckError(err)
	return resp
}

// RemoveTestDB : Test DB 삭제
func RemoveTestDB() {
	testDbPath := utils.GetWorkingDirPath() + "/test/test.db"

	if _, err := os.Stat(testDbPath); err == nil {
		err := os.Remove(testDbPath)
		if err != nil {
			log.Fatal(err)
		}
	}
}

type initTestArg struct {
	method  string
	url     string
	dataArg *bytes.Buffer
}

func GetDataArg(arg interface{}) *bytes.Buffer {
	testDataBytes, err := json.Marshal(arg)
	CheckError(err)

	return bytes.NewBuffer(testDataBytes)
}

func RespUnmarshal(resp *http.Response) RespForTest {
	var respForTest RespForTest
	body, err := ioutil.ReadAll(resp.Body)
	CheckError(err)
	err = json.Unmarshal(body, &respForTest)
	CheckError(err)

	return respForTest
}

func initTest() {
	addressLabel := "테스트집"
	roadAddress := "테스트집1길"
	lotAddress := "테스트집94-3"
	subAddress := "101호"
	postNo := "11113"
	contactNo := "010-1234-5678"
	args := []initTestArg{
		{
			"POST", Ts.URL + "/address/create", GetDataArg(models.UserAddress{
				UserNo:       1,
				AddressLabel: &addressLabel,
				RoadAddress:  &roadAddress,
				LotAddress:   &lotAddress,
				SubAddress:   &subAddress,
				PostNo:       &postNo,
				ContactNo:    &contactNo,
				IsBasic:      1,
				StatusCode:   "normal",
			})},
		{
			"POST", Ts.URL + "/card/create", GetDataArg(models.UserCard{
				UserNo:     1,
				CardKey:    "기본카드key",
				IsBasic:    1,
				StatusCode: "normal",
			})},
		{
			"POST", Ts.URL + "/subs/create", GetDataArg(models.SubsMst{
				UserNo:     1,
				CardRegNo:  1,
				SubsPrice:  1231,
				PeriodDay:  3,
				StatusCode: "normal",
			})},
		{
			"POST", Ts.URL + "/order/create", GetDataArg(models.OrderMst{
				OrderNo:    1,
				UserNo:     1,
				CardRegNo:  1,
				OrderType:  "테스트오더",
				OrderPrice: 123131,
				StatusCode: "init",
			})},
		{
			"POST", Ts.URL + "/banner/create", GetDataArg(models.BannerMst{
				BannerNo:     1,
				BannerType:   "테스트배너타입",
				BannerCode:   "테스트배너코드",
				BannerTitle:  "테스트",
				BannerDesc:   "testDesc",
				BannerImgURL: "testURL",
				StatusCode:   "ready",
			})},
		{
			"POST", Ts.URL + "/agreement/create", GetDataArg(models.AgreementMst{
				AgreementNo:    1,
				AgreementType:  "약관타입",
				AgreementTitle: "약관이름",
				AgreementDesc:  "약관설명",
				AttachFileURL:  "testURL",
				StatusCode:     "ready",
			})},
		{
			"POST", Ts.URL + "/faq/create", GetDataArg(models.FaqMst{
				FaqNo:      1,
				FaqType:    "FaqType",
				FaqTitle:   "FaqTitle",
				FaqDesc:    "FaqDesc",
				StatusCode: "ready",
			})},
	}

	for _, arg := range args {
		resp := RequestForTest(arg.method, arg.url, arg.dataArg)
		CheckCondition(resp.StatusCode == 200, fmt.Sprintf("CREATE %s is failed", arg.url))
	}
}

func TestMain(m *testing.M) {
	RemoveTestDB()
	os.Setenv("IS_TEST", "TRUE")

	r := router.SetupRouter("test")
	Ts = httptest.NewServer(r)
	defer Ts.Close()

	models.CreateAllTableIfNotExists()

	// 테스트 유저 생성 Post data 준비
	email := "testbot@gmail.com"
	testUser := models.UserMst{
		UserType:      "kakao",
		UserID:        "testbotID",
		UserSecretKey: "testbotPassword",
		UserName:      "testbotName",
		UserEmail:     &email,
		UserPhone:     "01012345678",
		IsMktAgree:    0,
		StatusCode:    "normal",
	}
	userBytes, err := json.Marshal(testUser)
	CheckError(err)

	// 테스트 유저 생성 Post
	resp, err := http.Post(Ts.URL+"/user/create", "application/json", bytes.NewBuffer(userBytes))
	CheckError(err)

	CheckCondition(resp.StatusCode == 200, resp)

	// 테스트 유저 JWT Auth Token 발급
	client := &http.Client{}
	req, err := http.NewRequest("POST", Ts.URL+"/login", nil)
	if err != nil {
		log.Fatal(err)
	}
	req.SetBasicAuth("testbotID", "testbotPassword")

	// Auth Token 획득
	resp, err = client.Do(req)
	CheckError(err)
	CheckCondition(resp.StatusCode == 200, resp)

	body, err := ioutil.ReadAll(resp.Body)
	CheckError(err)

	var tokenResp TokenResp
	err = json.Unmarshal(body, &tokenResp)
	CheckError(err)

	// Auth 토큰 저장
	Bearer = tokenResp.Token

	// 기초 데이터 삽입
	initTest()

	// 테스트 시작
	runTests := m.Run()
	os.Exit(runTests)
}
