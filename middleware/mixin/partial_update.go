package mixin

import (
	"fmt"
	"yelloment-api/models"
	"yelloment-api/utils"
)

// PartialUpdate : 부분 업데이트
func PartialUpdate(modelName string, wherePhrase string, modelAddr interface{}, skipFields []string, allowedEmptyFields []string) (bool, string) {
	structMap := utils.GetStructMap(modelAddr, skipFields, true, allowedEmptyFields)
	tableName := models.GetTableName(modelName)
	return models.Update(structMap, tableName, wherePhrase)
}

// PartialUpdateOwned : 소유자 데이터 부분 업데이트
func PartialUpdateOwned(modelName string, wherePhrase string, modelAddr interface{}, skipFields []string, allowedEmptyFields []string, userNo int) (bool, string) {
	structMap := utils.GetStructMap(modelAddr, []string{}, true, allowedEmptyFields)
	tableName := models.GetTableName(modelName)
	completedWherePhrase := fmt.Sprintf("UserNo = %d AND %s", userNo, wherePhrase)
	return models.Update(structMap, tableName, completedWherePhrase)
}
