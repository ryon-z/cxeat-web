package test

import (
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestListFaqNormally(t *testing.T) {
	resp := RequestForTest("GET", Ts.URL+"/faq/list", nil)
	respForTest := RespUnmarshal(resp)

	assert.Equal(t, resp.StatusCode, 200)
	assert.Equal(t, len(respForTest.Error), 0)
	dataFields := respForTest.Data.([]interface{})
	assert.Equal(t, len(dataFields), 1)
	dataFieldValue := dataFields[0].(map[string]interface{})
	assert.Equal(t, dataFieldValue["FaqNo"].(float64), float64(1))
}
