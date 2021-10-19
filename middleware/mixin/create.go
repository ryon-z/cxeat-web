package mixin

import (
	"strconv"
	"yelloment-api/models"
	"yelloment-api/utils"
)

// Create : 생성
func Create(modelName string, modelAddr interface{}, skipFields []string, withoutEmptyFields bool, allowedEmptyFields []string) (bool, interface{}) {
	tableName := models.GetTableName(modelName)
	structMap := utils.GetStructMap(modelAddr, skipFields, withoutEmptyFields, allowedEmptyFields)
	namePhrase, valuePhrase := utils.GetColumnAndValuePhrase(structMap)
	isSuccess, message := models.Insert(namePhrase, valuePhrase, tableName, true)

	return isSuccess, message
}

// CreateOwned : 소유자 데이터 생성
func CreateOwned(modelName string, userNo int, modelAddr interface{}, skipFields []string, withoutEmptyFields bool, allowedEmptyFields []string) (bool, interface{}) {
	tableName := models.GetTableName(modelName)
	structMap := utils.GetStructMap(modelAddr, skipFields, withoutEmptyFields, allowedEmptyFields)
	structMap["UserNo"] = strconv.Itoa(userNo)
	namePhrase, valuePhrase := utils.GetColumnAndValuePhrase(structMap)
	isSuccess, message := models.Insert(namePhrase, valuePhrase, tableName, true)

	return isSuccess, message
}

// CreateMulti : 여러 데이터 생성
func CreateMulti(modelName string, modelSliceAddr interface{}, skipFields []string, withoutEmptyFields bool, allowedEmptyFields []string) (bool, interface{}) {
	tableName := models.GetTableName(modelName)
	isSuccess, result, valuePhrase := utils.GetStructColumnAndValuePhrase(modelSliceAddr, skipFields, withoutEmptyFields, allowedEmptyFields)
	if !isSuccess {
		return isSuccess, result
	}
	namePhrase := result
	isSuccess, message := models.Insert(namePhrase, valuePhrase, tableName, false)
	return isSuccess, message
}
