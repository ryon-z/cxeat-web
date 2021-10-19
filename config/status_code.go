package config

import (
	"encoding/json"
	"fmt"
	envutil "yelloment-api/env_util"
)

// ErrorPrototype : 에러 원형
type ErrorPrototype struct {
	YelloCode int
	Message   string
}

var newLineSymbol = envutil.GetGoDotEnvVariable("NEW_LINE_SYMBOL")

// ExposedErrorMessageMap : 노출 에러 문자열 맵
var ExposedErrorMessageMap = map[int]string{
	1:    "시스템 에러 발생",
	2:    "DB 관련 문제 발생",
	5:    "탈퇴한 계정입니다.",
	9:    "로그인 관련 문제 발생",
	5001: fmt.Sprintf("탈퇴한 계정입니다. 재가입은 newSignupDate에 가능합니다. %s(고객센터: 070-4166-6077)", newLineSymbol),
	9002: fmt.Sprintf("카카오 인증 코드가 만료되었습니다. %s 로그인 버튼을 눌러 다시 로그인을 시도해주세요.", newLineSymbol),
	9004: "카카오 소셜 계정에 등록된 전화번호가 없어 회원가입이 불가능합니다.",
	9014: "네이버 소셜 계정에 등록된 전화번화가 없어 회원가입이 불가능합니다.",
}

// StringifyError : 에러 json 문자열 화
func StringifyError(err ErrorPrototype) (isSuccess bool, result string) {
	strout, marshalErr := json.Marshal(err)
	if marshalErr != nil {
		return false, "{'Code': 1001, 'Message': 'stringifying data is failed' }"
	}

	return true, string(strout)
}

// UnstringifyError : 에러 json 문자열 struct 화
func UnstringifyError(errString string) (isSuccess bool, result ErrorPrototype) {
	if err := json.Unmarshal([]byte(errString), &result); err != nil {
		return false, ErrorPrototype{YelloCode: 1002, Message: err.Error()}
	}

	return true, result
}

// UserStatusCodes : 유저상태코드
var UserStatusCodes []string

// SubsStatusCodes : 구독상태코드
var SubsStatusCodes []string

// AddressStatusCodes : 주소상태코드
var AddressStatusCodes []string

// OrderStatusCodes : 주문상태코드
var OrderStatusCodes []string

// InitStatusCodes : 상태코드 초기화
func InitStatusCodes() {
	UserStatusCodes = []string{"normal", "dormancy", "leave", "block", "first"}
	SubsStatusCodes = []string{"normal", "pause", "cancle"}
	AddressStatusCodes = []string{"normal", "pause", "cancle"}
	OrderStatusCodes = []string{"init", "ready-delivery", "standby-delivery", "in-delivery", "done"}
}
