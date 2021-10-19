package middleware

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"yelloment-api/middleware/mixin"
	"yelloment-api/models"
	"yelloment-api/utils"

	"github.com/gin-gonic/gin"
)

// Generic : CRUD Generic middleware
func Generic(modelName string, genericType string, idFieldName string, autoFields []string) gin.HandlerFunc {
	lowerGenericType := strings.ToLower(genericType)
	fn := func(c *gin.Context) {
		var isSuccess bool
		var data interface{}
		switch lowerGenericType {
		case "list":
			isSuccess, data = list(modelName)
		case "retrieve":
			isSuccess, data = retrieve(c, modelName, idFieldName)
		case "create":
			var skipFields []string
			skipFields = autoFields
			skipFields = append(skipFields, idFieldName)
			isSuccess, data = create(c, modelName, skipFields)
		case "createwithid":
			isSuccess, data = create(c, modelName, autoFields)
		case "partialupdate":
			var skipFields []string
			skipFields = autoFields
			skipFields = append(skipFields, idFieldName)
			isSuccess, data = partialUpdate(c, modelName, idFieldName, skipFields, []string{})
		case "delete":
			isSuccess, data = delete(c, modelName, idFieldName)
		case "listowned":
			isSuccess, data = listOwned(c, modelName)
		case "retrieveowned":
			isSuccess, data = retrieveOwned(c, modelName, idFieldName)
		case "createowned":
			var skipFields []string
			skipFields = autoFields
			skipFields = append(skipFields, idFieldName)
			isSuccess, data = createOwned(c, modelName, skipFields)
		case "createwihtoutidowned":
			isSuccess, data = createOwned(c, modelName, autoFields)
		case "partialupdateowned":
			isSuccess, data = partialUpdateOwned(c, modelName, idFieldName)
		case "deleteowned":
			isSuccess, data = deleteOwned(c, modelName, idFieldName)
		default:
			c.JSON(http.StatusBadRequest, fmt.Sprintf("%s is not allowed genericType", genericType))
		}

		if data == "EOF" {
			data = "data parameters are wrong"
		}

		// Set Response
		if isSuccess {
			c.JSON(http.StatusOK, gin.H{"data": data})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": data})
		}
	}

	return gin.HandlerFunc(fn)
}

// list : 리스트 조회
func list(modelName string) (bool, interface{}) {
	isSuccess, result := mixin.List(modelName, `1=1`, "")

	return isSuccess, result
}

// listOwned : 소유자 허용 리스트 조회
func listOwned(c *gin.Context, modelName string) (bool, interface{}) {
	userNo := utils.GetUserNo(c)
	wherePhrase := fmt.Sprintf(`UserNo = %d`, userNo)
	isSuccess, result := mixin.List(modelName, wherePhrase, "")

	return isSuccess, result
}

// retrieve : 검색
func retrieve(c *gin.Context, modelName string, idFieldName string) (bool, interface{}) {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return false, err.Error()
	}

	wherePhrase := fmt.Sprintf("%s = %d", idFieldName, idInt)
	isSuccess, result := mixin.Retrieve(modelName, wherePhrase, idFieldName, idInt)

	return isSuccess, result
}

// retrieveOwned : 소유자 허용 검색
func retrieveOwned(c *gin.Context, modelName string, idFieldName string) (bool, interface{}) {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return false, err.Error()
	}

	userNo := utils.GetUserNo(c)
	wherePhrase := fmt.Sprintf("UserNo = %d AND %s = %d", userNo, idFieldName, idInt)
	isSuccess, result := mixin.Retrieve(modelName, wherePhrase, idFieldName, idInt)

	return isSuccess, result
}

// create : 생성
func create(c *gin.Context, modelName string, skipFields []string) (bool, interface{}) {
	modelAddr := models.GetModelAddr(modelName)
	if err := c.ShouldBindJSON(modelAddr); err != nil {
		return false, err.Error()
	}

	isSuccess, result := mixin.Create(modelName, modelAddr, skipFields, true, []string{})

	return isSuccess, result
}

// createOwned : 소유자 허용 생성
func createOwned(c *gin.Context, modelName string, skipFields []string) (bool, interface{}) {
	modelAddr := models.GetModelAddr(modelName)
	if err := c.ShouldBindJSON(modelAddr); err != nil {
		return false, err.Error()
	}

	if !utils.IsContainUserNo(modelAddr) {
		return false, "UserNo is not contained"
	}

	userNo := utils.GetUserNo(c)
	isSuccess, result := mixin.CreateOwned(modelName, userNo, modelAddr, skipFields, true, []string{})

	return isSuccess, result
}

// partialUpdate : 업데이트
func partialUpdate(c *gin.Context, modelName string, idFieldName string, skipFields []string, allowedEmptyFields []string) (bool, interface{}) {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return false, err.Error()
	}

	modelAddr := models.GetModelAddr(modelName)
	if err := c.ShouldBindJSON(modelAddr); err != nil {
		return false, err.Error()
	}

	wherePhrase := fmt.Sprintf("%s = %d", idFieldName, idInt)
	isSuccess, message := mixin.PartialUpdate(modelName, wherePhrase, modelAddr, skipFields, allowedEmptyFields)

	return isSuccess, message
}

// partialUpdateOwned : 소유자 허용 업데이트
func partialUpdateOwned(c *gin.Context, modelName string, idFieldName string) (bool, interface{}) {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return false, err.Error()
	}
	modelAddr := models.GetModelAddr(modelName)
	if err := c.ShouldBindJSON(modelAddr); err != nil {
		return false, err.Error()
	}

	if !utils.IsContainUserNo(modelAddr) {
		return false, "UserNo is not contained"
	}

	userNo := utils.GetUserNo(c)
	wherePhrase := fmt.Sprintf("%s = %d", idFieldName, idInt)
	isSuccess, message := mixin.PartialUpdateOwned(modelName, wherePhrase, modelAddr, []string{}, []string{}, userNo)

	return isSuccess, message
}

// delete : 삭제
func delete(c *gin.Context, modelName string, idFieldName string) (bool, interface{}) {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return false, err.Error()
	}

	wherePhrase := fmt.Sprintf("%s = %d", idFieldName, idInt)
	isSuccess, message := mixin.Delete(modelName, wherePhrase)

	return isSuccess, message
}

// deleteOwned : 소유자 허용 삭제
func deleteOwned(c *gin.Context, modelName string, idFieldName string) (bool, interface{}) {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return false, err.Error()
	}

	userNo := utils.GetUserNo(c)
	wherePhrase := fmt.Sprintf("%s = %d", idFieldName, idInt)
	isSuccess, message := mixin.DeleteOwned(modelName, wherePhrase, userNo)

	return isSuccess, message
}
