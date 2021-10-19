package test

import (
	"testing"
	"yelloment-api/models"

	"github.com/go-playground/assert/v2"
)

func TestListOrderNormally(t *testing.T) {
	resp := RequestForTest("GET", Ts.URL+"/order/list", nil)
	respForTest := RespUnmarshal(resp)

	assert.Equal(t, resp.StatusCode, 200)
	assert.Equal(t, len(respForTest.Error), 0)
	dataFields := respForTest.Data.([]interface{})
	assert.Equal(t, len(dataFields), 1)
	dataFieldValue := dataFields[0].(map[string]interface{})
	assert.Equal(t, dataFieldValue["OrderNo"].(float64), float64(1))
}

func TestCancelOrderNormally(t *testing.T) {
	testDataBytes := GetDataArg(models.OrderMst{
		OrderNo: 1,
	})

	resp := RequestForTest("POST", Ts.URL+"/order/cancel", testDataBytes)
	respForTest := RespUnmarshal(resp)

	assert.Equal(t, resp.StatusCode, 200)
	assert.Equal(t, respForTest.Data.(string), "success")

	resp = RequestForTest("GET", Ts.URL+"/order/list", nil)
	respForTest = RespUnmarshal(resp)

	assert.Equal(t, resp.StatusCode, 200)
	assert.Equal(t, len(respForTest.Error), 0)
	dataFields := respForTest.Data.([]interface{})
	assert.Equal(t, len(dataFields), 0)
}
