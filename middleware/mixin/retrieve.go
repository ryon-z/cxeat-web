package mixin

import (
	"fmt"
	"yelloment-api/models"
	"yelloment-api/utils"
)

func getRetrieveIDWherePhrase(idFieldName string, id interface{}) (bool, string) {
	var wherePhrase string

	switch id.(type) {
	case int:
		wherePhrase = fmt.Sprintf("%s = %d", idFieldName, id)
	case float32, float64:
		wherePhrase = fmt.Sprintf("%s = %g", idFieldName, id)
	case string:
		isSuccess, message := utils.CheckSQLInjection(id.(string))
		if !isSuccess {
			return false, message
		}
		wherePhrase = fmt.Sprintf("%s = %s", idFieldName, id)

	default:
		return false, fmt.Sprintf("type of id is not allowed, type of id:%T\n", id)
	}

	return true, wherePhrase
}

// Retrieve : 조회
func Retrieve(modelName string, wherePhrase string, idFieldName string, id interface{}) (bool, interface{}) {
	isSuccess, message := getRetrieveIDWherePhrase(idFieldName, id)
	if !isSuccess {
		return false, message
	}

	completedWherePhrase := wherePhrase + " AND " + message
	return models.Select(modelName, completedWherePhrase, "")
}
