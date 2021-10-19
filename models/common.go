package models

import (
	"fmt"
	"strings"
	"yelloment-api/database"
	envutil "yelloment-api/env_util"
)

// GetTableName : 모델의 테이블명 획득
func GetTableName(modelName string) string {
	return GetModelObject(modelName, "tableName").(string)
}

// GetModelAddr : 모델 구조체 주소값 획득
func GetModelAddr(modelName string) interface{} {
	return GetModelObject(modelName, "addr")
}

// GetModelArrAddr : 모델 구조체 arr 주소값 획득
func GetModelArrAddr(modelName string) interface{} {
	return GetModelObject(modelName, "arrAddr")
}

// GetModelObject : 모델 구조체 관련 객체
func GetModelObject(modelName string, mode string) interface{} {
	var result interface{}
	allowedModes := []string{"arrAddr", "addr", "tableName"}
	for i, allowedMode := range allowedModes {
		if mode == allowedMode {
			break
		}

		if i == len(allowedModes)-1 {
			panic(fmt.Sprintf("mode is not allowed, mode %s", mode))

		}
	}

	switch modelName {
	case "orderItem":
		switch mode {
		case "arrAddr":
			result = &[]OrderItem{}
		case "addr":
			result = &OrderItem{}
		case "tableName":
			result = OrderItem{}.TableName()
		}
	case "orderMst":
		switch mode {
		case "arrAddr":
			result = &[]OrderMst{}
		case "addr":
			result = &OrderMst{}
		case "tableName":
			result = OrderMst{}.TableName()
		}
	case "orderPayment":
		switch mode {
		case "arrAddr":
			result = &[]OrderPayment{}
		case "addr":
			result = &OrderPayment{}
		case "tableName":
			result = OrderPayment{}.TableName()
		}
	case "subsBundleItem":
		switch mode {
		case "arrAddr":
			result = &[]SubsBundleItem{}
		case "addr":
			result = &SubsBundleItem{}
		case "tableName":
			result = SubsBundleItem{}.TableName()
		}
	case "subsBundleMst":
		switch mode {
		case "arrAddr":
			result = &[]SubsBundleMst{}
		case "addr":
			result = &SubsBundleMst{}
		case "tableName":
			result = SubsBundleMst{}.TableName()
		}
	case "subsItemMst":
		switch mode {
		case "arrAddr":
			result = &[]SubsItemMst{}
		case "addr":
			result = &SubsItemMst{}
		case "tableName":
			result = SubsItemMst{}.TableName()
		}
	case "subsPlan":
		switch mode {
		case "arrAddr":
			result = &[]SubsPlan{}
		case "addr":
			result = &SubsPlan{}
		case "tableName":
			result = SubsPlan{}.TableName()
		}
	case "subsPlanOption":
		switch mode {
		case "arrAddr":
			result = &[]SubsPlanOption{}
		case "addr":
			result = &SubsPlanOption{}
		case "tableName":
			result = SubsPlanOption{}.TableName()
		}
	case "userAddress":
		switch mode {
		case "arrAddr":
			result = &[]UserAddress{}
		case "addr":
			result = &UserAddress{}
		case "tableName":
			result = UserAddress{}.TableName()
		}
	case "userCard":
		switch mode {
		case "arrAddr":
			result = &[]UserCard{}
		case "addr":
			result = &UserCard{}
		case "tableName":
			result = UserCard{}.TableName()
		}
	case "userExtraInfo":
		switch mode {
		case "arrAddr":
			result = &[]UserExtraInfo{}
		case "addr":
			result = &UserExtraInfo{}
		case "tableName":
			result = UserExtraInfo{}.TableName()
		}
	case "userMst":
		switch mode {
		case "arrAddr":
			result = &[]UserMst{}
		case "addr":
			result = &UserMst{}
		case "tableName":
			result = UserMst{}.TableName()
		}
	case "subsMst":
		switch mode {
		case "arrAddr":
			result = &[]SubsMst{}
		case "addr":
			result = &SubsMst{}
		case "tableName":
			result = SubsMst{}.TableName()
		}
	case "bannerMst":
		switch mode {
		case "arrAddr":
			result = &[]BannerMst{}
		case "addr":
			result = &BannerMst{}
		case "tableName":
			result = BannerMst{}.TableName()
		}
	case "agreementMst":
		switch mode {
		case "arrAddr":
			result = &[]AgreementMst{}
		case "addr":
			result = &AgreementMst{}
		case "tableName":
			result = AgreementMst{}.TableName()
		}
	case "faqMst":
		switch mode {
		case "arrAddr":
			result = &[]FaqMst{}
		case "addr":
			result = &FaqMst{}
		case "tableName":
			result = FaqMst{}.TableName()
		}
	case "codeMst":
		switch mode {
		case "arrAddr":
			result = &[]CodeMst{}
		case "addr":
			result = &CodeMst{}
		case "tableName":
			result = CodeMst{}.TableName()
		}
	case "codeType":
		switch mode {
		case "arrAddr":
			result = &[]CodeType{}
		case "addr":
			result = &CodeType{}
		case "tableName":
			result = CodeType{}.TableName()
		}
	case "tagGroup":
		switch mode {
		case "arrAddr":
			result = &[]TagGroup{}
		case "addr":
			result = &TagGroup{}
		case "tableName":
			result = TagGroup{}.TableName()
		}
	case "tag":
		switch mode {
		case "arrAddr":
			result = &[]Tag{}
		case "addr":
			result = &Tag{}
		case "tableName":
			result = Tag{}.TableName()
		}
	case "reviewMst":
		switch mode {
		case "arrAddr":
			result = &[]ReviewMst{}
		case "addr":
			result = &ReviewMst{}
		case "tableName":
			result = ReviewMst{}.TableName()
		}
	case "itemReview":
		switch mode {
		case "arrAddr":
			result = &[]ItemReview{}
		case "addr":
			result = &ItemReview{}
		case "tableName":
			result = ItemReview{}.TableName()
		}
	case "orderHist":
		switch mode {
		case "arrAddr":
			result = &[]OrderHist{}
		case "addr":
			result = &OrderHist{}
		case "tableName":
			result = OrderHist{}.TableName()
		}
	case "subsHist":
		switch mode {
		case "arrAddr":
			result = &[]SubsHist{}
		case "addr":
			result = &SubsHist{}
		case "tableName":
			result = SubsHist{}.TableName()
		}
	default:
		panic(fmt.Sprintf("modelName is not allowed, modelName %s", modelName))
	}

	return result
}

// Select : Select 실행
func Select(modelName string, wherePhrase string, orderPhrase string) (bool, interface{}) {
	modelArrAddr := GetModelArrAddr(modelName)
	tableName := GetTableName(modelName)
	emptyModelArr := GetModelAddr(modelName)
	fieldNamesString := strings.Join(envutil.GetStructFieldNames(emptyModelArr, []string{}), ", ")

	sqlQuery := fmt.Sprintf(`
		SELECT %s FROM %s
		WHERE (%s)
	;`, fieldNamesString, tableName, wherePhrase)

	if orderPhrase != "" {
		sqlQuery = strings.Replace(sqlQuery, ";", "", 1)
		sqlQuery = fmt.Sprintf("%s ORDER BY %s;", sqlQuery, orderPhrase)
	}
	fmt.Println(sqlQuery)

	err := database.DB.Select(modelArrAddr, sqlQuery)
	if err != nil {
		return false, err.Error()
	}

	return true, modelArrAddr
}

// Insert : insert 실행
func Insert(namePhrase string, valuePhrase string, tableName string, remainValuesParentheses bool) (bool, interface{}) {
	var wrapper string
	if remainValuesParentheses {
		wrapper = fmt.Sprintf("VALUES (%s)", valuePhrase)
	} else {
		wrapper = fmt.Sprintf("VALUES %s", valuePhrase)
	}
	sqlQuery := fmt.Sprintf(`
		INSERT INTO %s (%s)
		%s
	;`, tableName, namePhrase, wrapper)
	fmt.Println(sqlQuery)

	instance, err := database.DB.Exec(sqlQuery)
	if err != nil {
		return false, err.Error()
	}

	lastInsertID, err := instance.LastInsertId()
	if err != nil {
		return false, err.Error()
	}

	return true, lastInsertID
}

// Delete : delete 실행
func Delete(modelName string, wherePhrase string) (bool, string) {
	tableName := GetTableName(modelName)
	sqlQuery := fmt.Sprintf(`
		DELETE FROM %s
		WHERE (%s)
	;`, tableName, wherePhrase)
	fmt.Println(sqlQuery)

	_, err := database.DB.Exec(sqlQuery)
	if err != nil {
		return false, err.Error()
	}

	return true, "success"
}

// Update : update 실행
func Update(structMap map[string]string, tableName string, wherePhrase string) (bool, string) {
	var updatePhraseSlice []string
	var updatePhrase string

	for key, value := range structMap {
		elem := fmt.Sprintf("%s = %s", key, value)
		updatePhraseSlice = append(updatePhraseSlice, elem)
	}

	updatePhrase = strings.Join(updatePhraseSlice, ",")

	sqlQuery := fmt.Sprintf(`
		UPDATE %s SET %s WHERE (%s)
	;`, tableName, updatePhrase, wherePhrase)
	fmt.Println(sqlQuery)

	_, err := database.DB.Exec(sqlQuery)
	if err != nil {
		return false, err.Error()
	}

	return true, "success"
}

// CreateAllTableIfNotExists : 등록된 모든 모델의 테이블이 존재하지 않을 시 생성(오직 테스트 시에만 사용)
func CreateAllTableIfNotExists() {
	// database.DB.AutoMigrate(&OrderItem{})
	// database.DB.AutoMigrate(&OrderMst{})
	// database.DB.AutoMigrate(&OrderPayment{})
	// database.DB.AutoMigrate(&SubsBundleItem{})
	// database.DB.AutoMigrate(&SubsBundleMst{})
	// database.DB.AutoMigrate(&SubsItemMst{})
	// database.DB.AutoMigrate(&SubsPlan{})
	// database.DB.AutoMigrate(&SubsPlanOption{})
	// database.DB.AutoMigrate(&UserAddress{})
	// database.DB.AutoMigrate(&UserCard{})
	// database.DB.AutoMigrate(&UserExtraInfo{})
	// database.DB.AutoMigrate(&UserMst{})
	// database.DB.AutoMigrate(&UserSubs{})
	// database.DB.AutoMigrate(&BannerMst{})
	// database.DB.AutoMigrate(&AgreementMst{})
	// database.DB.AutoMigrate(&FaqMst{})
	// database.DB.AutoMigrate(&CodeMst{})
	// database.DB.AutoMigrate(&CodeType{})
}
