package mixin

import (
	"fmt"
	"yelloment-api/models"
)

// Delete : 삭제
func Delete(modelName string, wherePhrase string) (bool, string) {
	return models.Delete(modelName, wherePhrase)
}

// DeleteOwned : 소유자 데이터 삭제
func DeleteOwned(modelName string, wherePhrase string, userNo int) (bool, string) {
	completedWherePhrase := fmt.Sprintf("UserNo = %d AND %s", userNo, wherePhrase)
	return models.Delete(modelName, completedWherePhrase)
}
