package controllers

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	"yelloment-api/database"
	envutil "yelloment-api/env_util"
	"yelloment-api/models"
	"yelloment-api/utils"

	"github.com/gin-gonic/gin"
	"github.com/muesli/cache2go"
)

// DesiredDeliveryDates : 희망배송일 구조체
type DesiredDeliveryDates struct {
	Dow   string
	Dates string
}

// ItemCategory : item category 구조체
type ItemCategory struct {
	Label string
	Value string
}

// getCodeTypeWherePhrase : codeType 조건절 생성, codeTypes가 empty면 모든 rows 조회
func getCodeTypeWherePhrase(codeTypes []string) string {
	var elems = []string{}
	if len(codeTypes) == 0 {
		elems = []string{"1 = 1"}
	}

	columnName := "codeType"
	for _, codeType := range codeTypes {
		elems = append(elems, fmt.Sprintf("%s = '%s'", columnName, codeType))
	}

	return strings.Join(elems, " OR ")
}

// GetCodes : 코드 획득
func GetCodes(codeTypes []string) (isSuccess bool, data interface{}) {
	codes := []models.CodeMst{}

	mstTableName := models.GetTableName("codeMst")
	typeTableName := models.GetTableName("codeType")

	codeTypeWherePhrase := getCodeTypeWherePhrase(codeTypes)
	emptyModelArr := models.GetModelAddr("codeMst")
	mstFieldNamesString := strings.Join(envutil.GetStructFieldNames(emptyModelArr, []string{}), ", ")
	mstFieldNamesStringA := "A." + strings.Join(envutil.GetStructFieldNames(emptyModelArr, []string{}), ", A.")

	sqlQuery := fmt.Sprintf(`
		SELECT %s
		FROM (
			SELECT %s FROM %s
			WHERE (%s)
			and IsUse = 1
		) AS A
		JOIN (SELECT CodeType FROM %s
			WHERE IsUse = 1
		) AS B
		ON A.CodeType = B.CodeType
		ORDER BY A.CodeOrder
	;`, mstFieldNamesStringA, mstFieldNamesString, mstTableName, codeTypeWherePhrase, typeTableName)
	fmt.Println(sqlQuery)

	err := database.DB.Select(&codes, sqlQuery)
	if err != nil {
		fmt.Println("err.Error()", err.Error())
		return false, err.Error()
	}

	return true, &codes
}

// getDesiredDeliveryDates : 오늘 + alphaDay:alphaHour시간 기준으로 deliveryStartLimit 안의 deliveryDows 별 배송가능일 획득
func getDesiredDeliveryDates(alphaDay int, alphaHour int, deliveryDows *[]models.CodeMst, deliveryStartLimit int) []DesiredDeliveryDates {
	desiredDeliveryDates := map[string][]string{}
	engDows := utils.EngDows

	loc, _ := time.LoadLocation(envutil.GetGoDotEnvVariable("TIMEZONE"))
	now := time.Now().In(loc)
	baseDate := now.Add(time.Hour*24*time.Duration(alphaDay) + time.Hour*time.Duration(alphaHour))
	baseDow := baseDate.Weekday().String()
	baseIndex := utils.StringIndexOf(baseDow, engDows)

	var selectedEngDow []string
	for _, code := range *deliveryDows {
		selectedEngDow = append(selectedEngDow, code.CodeKey)
	}

	for _, engDow := range selectedEngDow {
		var nextDate time.Time
		daysPast := 0
		desiredIndex := utils.StringIndexOf(engDow, engDows)
		if baseIndex <= desiredIndex {
			daysPast = desiredIndex - baseIndex
		} else {
			daysPast = 7 - (baseIndex - desiredIndex)
		}

		remainedDays := deliveryStartLimit - daysPast
		numDelivery := (remainedDays - (remainedDays % 7)) / 7
		nextDate = baseDate.AddDate(0, 0, daysPast)
		for i := 0; i < numDelivery; i++ {
			dateString := fmt.Sprintf("%04d-%02d-%02d", nextDate.Year(), nextDate.Month(), nextDate.Day())
			desiredDeliveryDates[engDow] = append(desiredDeliveryDates[engDow], dateString)
			nextDate = nextDate.AddDate(0, 0, 7)
		}
	}

	var result []DesiredDeliveryDates
	for key, value := range desiredDeliveryDates {
		var row DesiredDeliveryDates
		row.Dow = key

		// TEMPORARY 2021년 9월 23일 제외
		// --> start -->
		refinedValue := value
		for vI, v := range value {
			if v == "2021-09-23" {
				refinedValue = append(value[:vI], value[vI+1:]...)
				break
			}
		}
		row.Dates = strings.Join(refinedValue, ",")
		// <-- end < --

		// // below 1 line is Origin Code
		// row.Dates = strings.Join(value, ",")
		result = append(result, row)
	}
	fmt.Println("result", result)

	return result
}

// GetSubsItemsState : 구독 항목 state 획득
func GetSubsItemsState(alphaNumDays int, alphaNumHours int) (bool, gin.H) {
	var state = gin.H{}
	isSuccess, codes := GetCodes([]string{
		"ITEM_CATEGORY_EXP", "DELIVERY_PERIOD", "DELIVERY_DOW", "DELIVERY_START_LIMIT", "ITEM_CATEGORY", "BOX_TYPE", "ORDER_TYPE_EXP"})
	if !isSuccess {
		state = gin.H{"error": "Failed loading api"}
		return false, state
	}
	var boxInfos []models.CodeMst
	var itemCategoryExps []models.CodeMst
	var deliveryPeriods []models.CodeMst
	var itemCategories []models.CodeMst
	var deliveryDows []models.CodeMst
	var orderTypeExps []models.CodeMst
	var deliveryStartLimit int

	for _, code := range *codes.(*[]models.CodeMst) {
		if code.CodeType == "BOX_TYPE" {
			boxInfos = append(boxInfos, code)
		} else if code.CodeType == "ITEM_CATEGORY_EXP" {
			itemCategoryExps = append(itemCategoryExps, code)
		} else if code.CodeType == "DELIVERY_PERIOD" {
			deliveryPeriods = append(deliveryPeriods, code)
		} else if code.CodeType == "DELIVERY_START_LIMIT" {
			limit, err := strconv.Atoi(code.CodeKey)
			if err != nil {
				return false, state
			}
			deliveryStartLimit = limit
		} else if code.CodeType == "ITEM_CATEGORY" {
			itemCategories = append(itemCategories, code)
		} else if code.CodeType == "ORDER_TYPE_EXP" {
			orderTypeExps = append(orderTypeExps, code)
		} else {
			deliveryDows = append(deliveryDows, code)
		}
	}

	state["boxInfos"] = boxInfos
	state["itemCategoryExps"] = itemCategoryExps
	state["deliveryPeriods"] = deliveryPeriods
	state["itemCategories"] = itemCategories
	state["deliveryDows"] = deliveryDows
	state["desiredDeliveryDates"] = getDesiredDeliveryDates(alphaNumDays, alphaNumHours, &deliveryDows, deliveryStartLimit)
	state["reverseDowMap"] = utils.ReverseDowMap
	state["orderTypeExps"] = orderTypeExps

	return true, state
}

// GetItemCategoriesState : itemCategory의 value와 label을 가지고 string 구조체를 만들어 리턴
func GetItemCategoriesState(useCache bool) (isSuccess bool, state gin.H) {
	state = gin.H{}
	var codesResult interface{}

	if useCache {
		isSuccess, codesResult = GetCodesFromCache([]string{"ITEM_CATEGORY"})
		if !isSuccess {
			state = gin.H{"error": "Failed loading cache"}
			return false, state
		}
	} else {
		isSuccess, codesResult = GetCodes([]string{"ITEM_CATEGORY"})
		if !isSuccess {
			state = gin.H{"error": "Failed loading api"}
			return false, state
		}
	}

	var itemCategories = *codesResult.(*[]models.CodeMst)
	var result []ItemCategory

	for _, category := range itemCategories {
		result = append(result, ItemCategory{category.CodeLabel, category.CodeKey})
	}

	state["itemCategories"] = result

	return true, state
}

// MakeCateTypeLabels : 화면에 표시할 구독상품구성 라벨 문자열 리턴
func MakeCateTypeLabels(cateType string, cateTypeSep string) (bool, interface{}) {
	var labels []string

	isSuccess, state := GetItemCategoriesState(true)
	if !isSuccess {
		return false, state["error"].(string)
	}

	splited := strings.Split(cateType, cateTypeSep)
	for _, cateTypeValue := range splited {
		for _, category := range state["itemCategories"].([]ItemCategory) {
			if cateTypeValue == category.Value {
				labels = append(labels, category.Label)
				continue
			}
		}
	}

	return true, labels
}

// getCodeCache : 코드 캐시 획득
func getCodeCache() (isSuccess bool, result interface{}) {
	cache := cache2go.Cache("myCache")
	var codelist []models.CodeMst

	for {
		res, err := cache.Value("codelist")

		if err != nil {
			isSuccess, result := GetCodes([]string{})
			if !isSuccess {
				return false, result.(string)
			}
			cache.Add("codelist", 3*time.Minute, *result.(*[]models.CodeMst))
		} else {
			codelist = res.Data().([]models.CodeMst)
			break
		}
	}

	return true, codelist
}

// GetCodesFromCache : 같은 codeType을 가진 모든 코드 조회
func GetCodesFromCache(codeTypes []string) (isSuccess bool, result interface{}) {
	var codes []models.CodeMst
	isSuccess, result = getCodeCache()
	if !isSuccess {
		return false, result.(string)
	}

	codeList := result.([]models.CodeMst)
	for _, code := range codeList {
		if utils.StringInSlice(code.CodeType, codeTypes) {
			codes = append(codes, code)
		}
	}

	return true, &codes
}

// GetCodeLabelFromCache : 캐시에서 코드 라벨 획득
func GetCodeLabelFromCache(codeType string, codeKey string) (isSuccess bool, result string) {
	isSuccess, cacheResult := getCodeCache()
	if !isSuccess {
		return false, cacheResult.(string)
	}

	for _, code := range cacheResult.([]models.CodeMst) {
		if code.CodeType == codeType && code.CodeKey == codeKey {
			return true, code.CodeLabel
		}
	}

	return true, ""
}
