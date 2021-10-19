package test

import (
	"testing"
	"yelloment-api/models"

	"github.com/go-playground/assert/v2"
)

func TestListCardNormally(t *testing.T) {
	resp := RequestForTest("GET", Ts.URL+"/card/list", nil)
	respForTest := RespUnmarshal(resp)

	assert.Equal(t, resp.StatusCode, 200)
	assert.Equal(t, len(respForTest.Error), 0)
	dataFields := respForTest.Data.([]interface{})
	assert.Equal(t, len(dataFields), 1)
	dataFieldValue := dataFields[0].(map[string]interface{})
	assert.Equal(t, dataFieldValue["CardRegNo"].(float64), float64(1))
}

func TestAddCardNormally(t *testing.T) {
	testDataBytes := GetDataArg(models.UserCard{
		UserNo:     1,
		CardKey:    "기본카드key2",
		IsBasic:    0,
		StatusCode: "it is not allowed StatusCode",
	})

	resp := RequestForTest("POST", Ts.URL+"/card/add", testDataBytes)
	respForTest := RespUnmarshal(resp)

	assert.Equal(t, resp.StatusCode, 200)
	assert.Equal(t, respForTest.Data.(string), "success")

	resp = RequestForTest("GET", Ts.URL+"/card/list", nil)
	respForTest = RespUnmarshal(resp)

	assert.Equal(t, resp.StatusCode, 200)
	assert.Equal(t, len(respForTest.Error), 0)
	dataFields := respForTest.Data.([]interface{})
	assert.Equal(t, len(dataFields), 2)
	dataFieldValue := dataFields[1].(map[string]interface{})
	assert.Equal(t, dataFieldValue["CardRegNo"].(float64), float64(2))
	assert.Equal(t, dataFieldValue["CardKey"].(string), "기본카드key2")
	assert.Equal(t, dataFieldValue["StatusCode"].(string), "normal")
}

func TestRemoveCardNormally(t *testing.T) {
	testDataBytes := GetDataArg(models.UserCard{
		CardRegNo: 1,
	})

	resp := RequestForTest("POST", Ts.URL+"/card/remove", testDataBytes)
	respForTest := RespUnmarshal(resp)

	assert.Equal(t, resp.StatusCode, 200)
	assert.Equal(t, respForTest.Data.(string), "success")

	resp = RequestForTest("GET", Ts.URL+"/card/list", nil)
	respForTest = RespUnmarshal(resp)

	assert.Equal(t, resp.StatusCode, 200)
	assert.Equal(t, len(respForTest.Error), 0)
	dataFields := respForTest.Data.([]interface{})
	assert.Equal(t, len(dataFields), 1)
	dataFieldValue := dataFields[0].(map[string]interface{})
	assert.Equal(t, dataFieldValue["CardRegNo"].(float64), float64(2))
	assert.Equal(t, dataFieldValue["CardKey"].(string), "기본카드key2")
	assert.Equal(t, dataFieldValue["StatusCode"].(string), "normal")
}
