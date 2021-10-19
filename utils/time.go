package utils

import (
	"reflect"
	"time"
)

// EngDows : 영문 요일
var EngDows = []string{
	"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday",
}

// KorDows : 한글 요일
var KorDows = []string{
	"월", "화", "수", "목", "금", "토", "일",
}

// DowMap : 요일 맵(key=한글요일, value=영문요일)
var DowMap = map[string]string{}

// ReverseDowMap : 리버스 요일 맵(key=영문요일, value=한글요일)
var ReverseDowMap = map[string]string{}

// InitAllDowMap : 모든 Dow Map 초기화
func InitAllDowMap() {
	emptyDowMap := map[string]string{}
	if reflect.DeepEqual(DowMap, emptyDowMap) {
		for i, engDow := range EngDows {
			DowMap[KorDows[i]] = engDow
			ReverseDowMap[engDow] = KorDows[i]
		}
	}
}

// GetEngDow : 시간을 받아 영문 요일 획득
func GetEngDow(timeString string) (isSuccess bool, engDow string) {
	t, err := time.Parse(time.RFC3339, timeString)
	if err != nil {
		return false, err.Error()
	}

	return true, t.Weekday().String()
}

// GetKorDow : 시간을 받아 한글 요일 획득
func GetKorDow(timeString string) (isSuccess bool, korDow string) {
	t, err := time.Parse(time.RFC3339, timeString)
	if err != nil {
		return false, err.Error()
	}
	korDow = ReverseDowMap[t.Weekday().String()]
	return true, korDow
}
