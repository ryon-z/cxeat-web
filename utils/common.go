package utils

import (
	"fmt"
	"net/http"
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"yelloment-api/config"
	envutil "yelloment-api/env_util"
	"yelloment-api/global"

	"github.com/gin-gonic/gin"
)

// GetStructColumnAndValuePhrase : GetStructMap의 StructSlice 버전
func GetStructColumnAndValuePhrase(
	inputStructSliceAddr interface{},
	skipFields []string,
	withoutEmptyFields bool,
	allowedEmptyFields []string) (
	isSuccess bool, namePhrase string, valuePhrase string) {
	namePhrase = ""
	valuePhrase = ""
	var names = []string{}
	var valuePhrases = []string{}

	switch reflect.TypeOf(inputStructSliceAddr).Elem().Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(inputStructSliceAddr).Elem()
		for i := 0; i < s.Len(); i++ {
			structMap := GetStructMap(s.Index(i), skipFields, withoutEmptyFields, allowedEmptyFields)
			var unitValuePhrase string
			var values = []string{}

			// 초기화
			if len(names) == 0 {
				for key, value := range structMap {
					names = append(names, key)
					values = append(values, value)
				}
				namePhrase = strings.Join(names, ",")
				unitValuePhrase = strings.Join(values, ",")
			} else {
				for _, key := range names {
					values = append(values, structMap[key])
				}
				unitValuePhrase = strings.Join(values, ",")
			}

			valuePhrases = append(valuePhrases, fmt.Sprintf("(%s)", unitValuePhrase))
		}
		isSuccess = true
		valuePhrase = strings.Join(valuePhrases, ", ")
	default:
		isSuccess = false
		namePhrase = "type of inputStructSliceAddr is wrong"
	}
	return isSuccess, namePhrase, valuePhrase
}

// GetStructMap : 구조체를 key가 fieldName이고 value가 fieldValue인 map으로 변환하여 리턴
func GetStructMap(inputStructAddr interface{}, skipFields []string, withoutEmptyFields bool, allowedEmptyFields []string) map[string]string {
	// skipFields - 필드값을 아예 저장하지 않음
	// withoutEmptyFields - 빈 문자열("")이나 빈 숫자(0 or 0.0)은 저장하지 않음
	// allowedEmptyFields - withoutEmptyFields가 true임에도 예외적으로 빈 문자열을 저장해야하는 필드
	var result = map[string]string{}
	var val reflect.Value

	switch inputStructAddr.(type) {
	case reflect.Value:
		val = inputStructAddr.(reflect.Value)
	default:
		val = reflect.ValueOf(inputStructAddr).Elem()
	}
Exit:
	for i := 0; i < val.NumField(); i++ {
		value := val.Field(i)
		fieldType := val.Type().Field(i).Type
		fieldName := val.Type().Field(i).Name
		for _, skipField := range skipFields {
			if skipField == fieldName {
				continue Exit
			}
		}

		switch fieldType.String() {
		case "string":
			strVal := value.Interface().(string)
			if withoutEmptyFields && !StringInSlice(fieldName, allowedEmptyFields) && strVal == "" {
				continue Exit
			}
			result[fieldName] = fmt.Sprintf("'%s'", value)
		case "*string":
			strVal := value.Interface().(*string)
			if strVal == nil {
				continue Exit
			}
			if withoutEmptyFields && !StringInSlice(fieldName, allowedEmptyFields) && *strVal == "" {
				continue Exit
			}
			result[fieldName] = fmt.Sprintf("'%s'", *strVal)

		case "int":
			intVal := value.Interface().(int)
			if withoutEmptyFields && !StringInSlice(fieldName, allowedEmptyFields) && intVal == 0 {
				continue Exit
			}
			strVal := strconv.FormatInt(int64(intVal), 10)
			result[fieldName] = strVal
		case "*int":
			intVal := value.Interface().(*int)
			if intVal == nil {
				continue Exit
			}
			if withoutEmptyFields && !StringInSlice(fieldName, allowedEmptyFields) && *intVal == 0 {
				continue Exit
			}
			strVal := strconv.FormatInt(int64(*intVal), 10)
			result[fieldName] = strVal
		case "float64":
			floatVal := value.Interface().(float64)
			if withoutEmptyFields && !StringInSlice(fieldName, allowedEmptyFields) && floatVal == 0.0 {
				continue Exit
			}
			strVal := strconv.FormatFloat(float64(floatVal), 'E', -1, 64)
			result[fieldName] = strVal
		case "*float64":
			floatVal := value.Interface().(*float64)
			if floatVal == nil {
				continue Exit
			}
			if withoutEmptyFields && !StringInSlice(fieldName, allowedEmptyFields) && *floatVal == 0.0 {
				continue Exit
			}
			strVal := strconv.FormatFloat(float64(*floatVal), 'E', -1, 64)
			result[fieldName] = strVal
		default:
			message := fmt.Sprintf("GetStructMap :: fieldType is wrong %s\n", fieldType.String())
			panic(message)
		}
	}

	return result
}

// GetColumnAndValuePhrase : key가 컬럼명이고 value가 컬럼값인 map을 받아 namePhrase와 valuePhrase를 리턴
func GetColumnAndValuePhrase(structMap map[string]string) (string, string) {
	var keys []string
	var values []string
	for key, value := range structMap {
		keys = append(keys, key)
		values = append(values, value)
	}

	return strings.Join(keys, ","), strings.Join(values, ",")
}

// IntInSlice : int target 값이 Slice 안에 있나 체크
func IntInSlice(target int, list []int) bool {
	for _, elem := range list {
		if elem == target {
			return true
		}
	}
	return false
}

// StringInSlice : string target 값이 Slice 안에 있나 체크
func StringInSlice(target string, list []string) bool {
	for _, elem := range list {
		if elem == target {
			return true
		}
	}
	return false
}

// StringIndexOf : string target 값의 Array index를 리턴
func StringIndexOf(target string, list []string) int {
	for i, elem := range list {
		if elem == target {
			return i
		}
	}
	return -1
}

// GetHTTPResponse : HTTP response 리턴
func GetHTTPResponse(isSuccess bool, message interface{}, c *gin.Context) {
	if isSuccess {
		c.JSON(http.StatusOK, gin.H{"data": message})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": message})
	}
	return
}

// EmptyArgument : argument가 비었는지 체크
func EmptyArgument(argument string) bool {
	return argument == ""
}

// GetWorkingDirPath : Working Directory 경로 얻기
func GetWorkingDirPath() string {
	var (
		_, b, _, _ = runtime.Caller(0)
		basepath   = filepath.Dir(filepath.Dir(b))
	)

	return basepath
}

// IsContainUserNo : 입력 받은 구조체에 UserNo가 포함 되었는지 확인
func IsContainUserNo(inputStruct interface{}) bool {
	for _, fieldName := range envutil.GetStructFieldNames(inputStruct, []string{}) {
		if fieldName == "UserNo" {
			return true
		}
	}

	return false
}

// CheckSQLInjection : Where 절에 injectionWords가 포함되었는지 체크
func CheckSQLInjection(wherePhrase string) (bool, string) {
	for _, injectionWord := range config.InjectionWords {
		if strings.Contains(wherePhrase, injectionWord) {
			return false, "'where' phares contained injection word"
		}
	}

	return true, "success"
}

// CombineTwoGinH : 두 개의 gin.H를 합성(동일한 키가 존재할 시 덮어쓰기 됨)
func CombineTwoGinH(passiveGinH *gin.H, activateGinH *gin.H) {
	for key, val := range *activateGinH {
		(*passiveGinH)[key] = val
	}
}

// GetIDFromCookie : 쿠키에 담긴 ID 값 획득
func GetIDFromCookie(c *gin.Context, cookieKey string) (isSuccess bool, result interface{}) {
	idStr, err := c.Cookie(cookieKey)
	if err != nil {
		return false, fmt.Sprintf("%s cookie not exists", cookieKey)

	}
	idInt, err := strconv.Atoi(idStr)
	if err != nil {
		return false, fmt.Sprintf("%s cookie is wrong", cookieKey)
	}

	return true, idInt
}

// GetGlobalState : 전역 변수 값 꺼내기
func GetGlobalState(c *gin.Context, key string, resultState gin.H) (isSuccess bool, state gin.H) {
	globalVar, exists := c.Get(key)
	if !exists {
		state = gin.H{}
		CombineCustomStateGlobalState(c, &state)
		slackMsg := fmt.Sprintf("[front]Refresh::c.Get('%s')", key)
		SendSlackMessage(SlackChannel, slackMsg)
		state["errorMessage"] = "시스템 에러 발생"
		state["errorCode"] = "W0030"
		return false, state
	}

	resultState[key] = globalVar

	return true, resultState
}

// CombineCustomStateGlobalState : custom state와 GlobalState를 합성(동일한 키가 존재할 시 customState에 덮어쓰기 됨)
func CombineCustomStateGlobalState(c *gin.Context, customState *gin.H) {
	for _, key := range global.GlobalStateKeys {
		if value, exists := c.Get(key); exists {
			(*customState)[key] = value
		}
	}
}
