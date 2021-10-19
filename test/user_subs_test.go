package test

import (
	"testing"
	"yelloment-api/models"

	"github.com/go-playground/assert/v2"
)

func TestListSubsNormally(t *testing.T) {
	resp := RequestForTest("GET", Ts.URL+"/subs/list", nil)
	respForTest := RespUnmarshal(resp)

	assert.Equal(t, resp.StatusCode, 200)
	assert.Equal(t, len(respForTest.Error), 0)
	dataFields := respForTest.Data.([]interface{})
	assert.Equal(t, len(dataFields), 1)
	dataFieldValue := dataFields[0].(map[string]interface{})
	assert.Equal(t, dataFieldValue["SubsNo"].(float64), float64(1))
}

func TestRequestSubsNormally(t *testing.T) {
	testDataBytes := GetDataArg(models.SubsMst{
		UserNo:     1,
		CardRegNo:  1,
		SubsPrice:  1231,
		PeriodDay:  3,
		StatusCode: "it is not allowed StatusCode",
	})

	resp := RequestForTest("POST", Ts.URL+"/subs/request", testDataBytes)
	respForTest := RespUnmarshal(resp)

	assert.Equal(t, resp.StatusCode, 200)
	assert.Equal(t, respForTest.Data.(string), "success")

	resp = RequestForTest("GET", Ts.URL+"/subs/list", nil)
	respForTest = RespUnmarshal(resp)

	assert.Equal(t, resp.StatusCode, 200)
	assert.Equal(t, len(respForTest.Error), 0)
	dataFields := respForTest.Data.([]interface{})
	assert.Equal(t, len(dataFields), 2)
	dataFieldValue := dataFields[1].(map[string]interface{})
	assert.Equal(t, dataFieldValue["SubsNo"].(float64), float64(2))
	assert.Equal(t, dataFieldValue["SubsPlanCode"].(string), "테스트구독코드2")
	assert.Equal(t, dataFieldValue["StatusCode"].(string), "normal")
}

func TestPauseSubsNormally(t *testing.T) {
	testDataBytes := GetDataArg(models.SubsMst{
		SubsNo: 1,
	})

	resp := RequestForTest("POST", Ts.URL+"/subs/pause", testDataBytes)
	respForTest := RespUnmarshal(resp)

	assert.Equal(t, resp.StatusCode, 200)
	assert.Equal(t, respForTest.Data.(string), "success")

	resp = RequestForTest("GET", Ts.URL+"/subs/list", nil)
	respForTest = RespUnmarshal(resp)

	assert.Equal(t, resp.StatusCode, 200)
	assert.Equal(t, len(respForTest.Error), 0)
	dataFields := respForTest.Data.([]interface{})
	assert.Equal(t, len(dataFields), 2)
	dataFieldValue := dataFields[0].(map[string]interface{})
	assert.Equal(t, dataFieldValue["SubsNo"].(float64), float64(1))
	assert.Equal(t, dataFieldValue["SubsPlanCode"].(string), "테스트구독코드")
	assert.Equal(t, dataFieldValue["StatusCode"].(string), "pause")
}
func TestCancelSubsNormally(t *testing.T) {
	testDataBytes := GetDataArg(models.SubsMst{
		SubsNo: 2,
	})

	resp := RequestForTest("POST", Ts.URL+"/subs/cancel", testDataBytes)
	respForTest := RespUnmarshal(resp)

	assert.Equal(t, resp.StatusCode, 200)
	assert.Equal(t, respForTest.Data.(string), "success")

	resp = RequestForTest("GET", Ts.URL+"/subs/list", nil)
	respForTest = RespUnmarshal(resp)

	assert.Equal(t, resp.StatusCode, 200)
	assert.Equal(t, len(respForTest.Error), 0)
	dataFields := respForTest.Data.([]interface{})
	assert.Equal(t, len(dataFields), 1)
	dataFieldValue := dataFields[0].(map[string]interface{})
	assert.Equal(t, dataFieldValue["SubsNo"].(float64), float64(1))
	assert.Equal(t, dataFieldValue["SubsPlanCode"].(string), "테스트구독코드")
	assert.Equal(t, dataFieldValue["StatusCode"].(string), "pause")
}

func TestEditSubsNormally(t *testing.T) {
	testDataBytes := GetDataArg(models.SubsMst{
		SubsNo: 1,
	})

	resp := RequestForTest("POST", Ts.URL+"/subs/edit", testDataBytes)
	respForTest := RespUnmarshal(resp)

	assert.Equal(t, resp.StatusCode, 200)
	assert.Equal(t, respForTest.Data.(string), "success")

	resp = RequestForTest("GET", Ts.URL+"/subs/list", nil)
	respForTest = RespUnmarshal(resp)

	assert.Equal(t, resp.StatusCode, 200)
	assert.Equal(t, len(respForTest.Error), 0)
	dataFields := respForTest.Data.([]interface{})
	assert.Equal(t, len(dataFields), 1)
	dataFieldValue := dataFields[0].(map[string]interface{})
	assert.Equal(t, dataFieldValue["SubsNo"].(float64), float64(1))
	assert.Equal(t, dataFieldValue["StartDate"].(string), "2011-11-11")
}
