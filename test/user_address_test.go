package test

import (
	"testing"
	"yelloment-api/models"

	"github.com/go-playground/assert/v2"
)

func TestListAddressNormally(t *testing.T) {
	resp := RequestForTest("GET", Ts.URL+"/address/list", nil)
	respForTest := RespUnmarshal(resp)

	assert.Equal(t, resp.StatusCode, 200)
	assert.Equal(t, len(respForTest.Error), 0)
	dataFields := respForTest.Data.([]interface{})
	assert.Equal(t, len(dataFields), 1)
	dataFieldValue := dataFields[0].(map[string]interface{})
	assert.Equal(t, dataFieldValue["AddressNo"].(float64), float64(1))
}

func TestAddAddressNormally(t *testing.T) {
	addressLabel := "테스트집"
	roadAddress := "테스트집1길"
	lotAddress := "테스트집94-3"
	subAddress := "101호"
	postNo := "11113"
	contactNo := "010-1234-5678"
	testDataBytes := GetDataArg(models.UserAddress{
		UserNo:       1,
		AddressLabel: &addressLabel,
		RoadAddress:  &roadAddress,
		LotAddress:   &lotAddress,
		SubAddress:   &subAddress,
		PostNo:       &postNo,
		ContactNo:    &contactNo,
		IsBasic:      1,
		StatusCode:   "it is not allowed StatusCode",
	})

	resp := RequestForTest("POST", Ts.URL+"/address/add", testDataBytes)
	respForTest := RespUnmarshal(resp)

	assert.Equal(t, resp.StatusCode, 200)
	assert.Equal(t, respForTest.Data.(string), "success")

	resp = RequestForTest("GET", Ts.URL+"/address/list", nil)
	respForTest = RespUnmarshal(resp)

	assert.Equal(t, resp.StatusCode, 200)
	assert.Equal(t, len(respForTest.Error), 0)
	dataFields := respForTest.Data.([]interface{})
	assert.Equal(t, len(dataFields), 2)
	dataFieldValue := dataFields[1].(map[string]interface{})
	assert.Equal(t, dataFieldValue["AddressNo"].(float64), float64(2))
	assert.Equal(t, dataFieldValue["AddressLabel"].(string), "테스트집2")
	assert.Equal(t, dataFieldValue["StatusCode"].(string), "normal")
}

func TestRemoveAddressNormally(t *testing.T) {
	testDataBytes := GetDataArg(models.UserAddress{
		AddressNo: 1,
	})

	resp := RequestForTest("POST", Ts.URL+"/address/remove", testDataBytes)
	respForTest := RespUnmarshal(resp)

	assert.Equal(t, resp.StatusCode, 200)
	assert.Equal(t, respForTest.Data.(string), "success")

	resp = RequestForTest("GET", Ts.URL+"/address/list", nil)
	respForTest = RespUnmarshal(resp)

	assert.Equal(t, resp.StatusCode, 200)
	assert.Equal(t, len(respForTest.Error), 0)
	dataFields := respForTest.Data.([]interface{})
	assert.Equal(t, len(dataFields), 1)
	dataFieldValue := dataFields[0].(map[string]interface{})
	assert.Equal(t, dataFieldValue["AddressNo"].(float64), float64(2))
	assert.Equal(t, dataFieldValue["AddressLabel"].(string), "테스트집2")
	assert.Equal(t, dataFieldValue["StatusCode"].(string), "normal")
}
