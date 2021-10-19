package test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"testing"
	"yelloment-api/models"

	"github.com/go-playground/assert/v2"
)

func TestInfoUserNormally(t *testing.T) {
	resp := RequestForTest("GET", Ts.URL+"/user/info", nil)
	body, err := ioutil.ReadAll(resp.Body)
	CheckError(err)

	var respForTest RespForTest
	err = json.Unmarshal(body, &respForTest)
	CheckError(err)

	assert.Equal(t, resp.StatusCode, 200)
	assert.Equal(t, len(respForTest.Error), 0)
	dataFields := respForTest.Data.([]interface{})
	assert.Equal(t, len(dataFields), 1)
	dataFieldValue := dataFields[0].(map[string]interface{})
	assert.Equal(t, dataFieldValue["UserNo"].(float64), float64(1))
}

func TestEditUserNormally(t *testing.T) {
	// 수정
	var RespForTest RespForTest
	testData := models.UserMst{
		UserName: "NewUserName",
	}
	testDataBytes, err := json.Marshal(testData)
	CheckError(err)

	resp := RequestForTest("POST", Ts.URL+"/user/edit", bytes.NewBuffer(testDataBytes))
	body, err := ioutil.ReadAll(resp.Body)
	CheckError(err)
	err = json.Unmarshal(body, &RespForTest)
	CheckError(err)

	// 수정 성공 여부(request 시 data에 들어있는 IdentityKey가 DB에 없는 경우, 무시되고 success return)
	assert.Equal(t, resp.StatusCode, 200)
	assert.Equal(t, RespForTest.Data.(string), "success")

	resp = RequestForTest("GET", Ts.URL+"/user/info", nil)
	body, err = ioutil.ReadAll(resp.Body)
	CheckError(err)

	err = json.Unmarshal(body, &RespForTest)
	CheckError(err)

	// 변경되었는지 확인
	assert.Equal(t, resp.StatusCode, 200)
	assert.Equal(t, len(RespForTest.Error), 0)
	dataFields := RespForTest.Data.([]interface{})
	dataFieldValue := dataFields[0].(map[string]interface{})
	assert.Equal(t, len(dataFields), 1)
	assert.Equal(t, dataFieldValue["UserName"].(string), "NewUserName")
	assert.Equal(t, dataFieldValue["UserID"].(string), "testbotID")
}
