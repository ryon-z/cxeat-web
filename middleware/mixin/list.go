package mixin

import (
	"yelloment-api/models"
)

// List : 리스트 조회
func List(modelName string, wherePhrase string, orderPhrase string) (bool, interface{}) {
	return models.Select(modelName, wherePhrase, orderPhrase)
}
