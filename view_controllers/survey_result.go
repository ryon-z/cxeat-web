package viewcontrollers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"yelloment-api/controllers"
	"yelloment-api/database"
	envutil "yelloment-api/env_util"
	"yelloment-api/models"
	"yelloment-api/utils"

	"github.com/gin-gonic/gin"
)

// GetOwnedTags : 소유자 태그 단일 조회(태그ID는 쿠키)
func GetOwnedTags(c *gin.Context, tagGroupNo int) (isSuccess bool, result interface{}) {
	userNo := utils.GetUserNo(c)

	var tag []models.Tag
	emptyModelArr := models.GetModelAddr("tag")
	fieldNamesStringA := "A." + strings.Join(envutil.GetStructFieldNames(emptyModelArr, []string{}), ", A.")
	fieldNamesString := strings.Join(envutil.GetStructFieldNames(emptyModelArr, []string{}), ",")

	sqlQuery := fmt.Sprintf(`
		SELECT %s
		FROM (SELECT %s FROM Tag 
			WHERE IsUse = 1 AND TagGroupNo = %d) AS A
		JOIN (SELECT TagGroupNo FROM TagGroup
			WHERE UserNo = %d AND TagGroupNo = %d) AS B
		ON (A.TagGroupNo = B.TagGroupNo)
	;`, fieldNamesStringA, fieldNamesString, tagGroupNo, userNo, tagGroupNo)
	fmt.Println(sqlQuery)

	err := database.DB.Select(&tag, sqlQuery)
	if err != nil {
		return false, err.Error()
	}

	return true, tag
}

// getScoreSourceMap : 점수 계산에 필요한 값들 획득
func getScoreSourceMap(tags []models.Tag) (isSuccess bool, sourceMap interface{}) {
	const adultWeight, youthWeight, childWeight = 10, 8, 5
	fmt.Println("tags", tags)
	result := map[string]int{
		"scoreFamily":       0,
		"isWantedRice":      1,
		"isWantedFruit":     1,
		"isWantedVegetable": 1,
		"isWantedGrain":     0,
		"isWantedQRice":     0,
		"numMeals":          0,
		"numFruits":         0,
		"numVegetables":     0,
		"numFamily":         0,
		"isWantedBrownRice": 0,
		"haveYoung":         0,
	}

	for _, tag := range tags {
		if tag.TagType == "FAMILY" {
			value, err := strconv.Atoi(tag.TagValue)
			if err != nil {
				return false, "value of FAMILY is not number"
			}

			switch tag.TagLabel {
			case "성인":
				result["scoreFamily"] += value * adultWeight
				result["numFamily"] += value
			case "청소년":
				result["scoreFamily"] += value * youthWeight
				result["numFamily"] += value
				result["haveYoung"] = 1
			case "아동":
				result["scoreFamily"] += value * childWeight
				result["numFamily"] += value
				result["haveYoung"] = 1
			}
		}

		if tag.TagType == "RICE" {
			if tag.TagLabel == "취식X" {
				result["isWantedRice"] = 0
			} else if tag.TagLabel == "즉석밥" {
				result["isWantedRice"] = 0
				result["isWantedQRice"] = 1
			}
		}

		if tag.TagType == "FRUIT" && tag.TagLabel == "과일X" {
			result["isWantedFruit"] = 0
		}

		if tag.TagType == "VEGET" && tag.TagLabel == "채소X" {
			result["isWantedVegetable"] = 0
		}

		if tag.TagType == "RICE_LIKE" && tag.TagLabel == "잡곡" {
			result["isWantedGrain"] = 1
		}

		if tag.TagType == "RICE_LIKE" && tag.TagLabel == "현미" {
			result["isWantedBrownRice"] = 1
		}

		if tag.TagType == "MEAL_CNT" {
			value, err := strconv.Atoi(tag.TagValue)
			if err != nil {
				return false, "value of MEAL_CNT is not number"
			}
			result["numMeals"] += value
		}

		if tag.TagType == "FRUIT_REQ" {
			result["numFruits"] += 1
		}

		if tag.TagType == "VEGET_REQ" {
			result["numVegetables"] += 1
		}
	}

	return true, result
}

// getScore : 점수 계산
func getScore(sourceMap map[string]int) int {
	const numMealsWeight, grainWeight, fruitsReqWeight, vegetableReqWeight = 1, 10, 1, 1

	riceScore := sourceMap["scoreFamily"]*sourceMap["isWantedRice"]*sourceMap["numMeals"]*numMealsWeight +
		sourceMap["isWantedGrain"]*grainWeight
	fruitsScore := sourceMap["scoreFamily"] * (sourceMap["isWantedFruit"] + sourceMap["numFruits"]*fruitsReqWeight)
	vegetableScore := sourceMap["scoreFamily"] * (sourceMap["isWantedVegetable"] + sourceMap["numVegetables"]*vegetableReqWeight)
	numRiceBows := sourceMap["numFamily"] * sourceMap["numMeals"]
	fmt.Println("riceScore", riceScore)
	fmt.Println("fruitsScore", fruitsScore)
	fmt.Println("vegetableScore", vegetableScore)
	fmt.Println("numRiceBows", numRiceBows)

	return riceScore + fruitsScore + vegetableScore + numRiceBows
}

// getSize : 크기등급 획득
func getSize(totalScore int) string {
	if totalScore < 100 {
		return "S"
	} else if totalScore < 200 {
		return "S+"
	} else if totalScore < 300 {
		return "R"
	} else if totalScore < 400 {
		return "R+"
	} else if totalScore < 500 {
		return "L"
	} else {
		return "L+"
	}
}

// getPhrase : 설명 문구
func getPhrase(size string, sourceMap map[string]int) (phrase string) {
	if size == "S" {
		if sourceMap["numFamily"] == 1 {
			return "혼자서도 딱 맞게 먹고 싶은<br /> 합리적인 당신을 위한 추천"
		} else if sourceMap["numFamily"] >= 2 && sourceMap["haveYoung"] == 1 {
			return "성장기 아이와 함께 즐기는 추천 구성"
		} else {
			return "꼭 필요한 것만 꼭 맞게 담은 추천 구성"
		}
	}

	if size == "S+" {
		if sourceMap["numFamily"] == 1 &&
			(sourceMap["isWantedBrownRice"] == 1 || sourceMap["isWantedGrain"] == 1) {
			return "한 끼를 먹어도 건강하게 즐기는 당신을 위한 추천"
		} else if sourceMap["numFamily"] == 1 {
			return "혼자서도 다양하게 즐기는 당신을 위한 추천"
		} else if sourceMap["numFamily"] >= 2 && sourceMap["haveYoung"] == 1 {
			return "성장기 아이와 함께 즐기는 추천 구성"
		} else {
			return "꼭 필요한 것만 쏙쏙 담은 추천 구성"
		}
	}

	if size == "R" {
		if sourceMap["haveYoung"] == 1 {
			return "성장기 아이와 함께 즐기는 추천 구성"
		} else if sourceMap["numFamily"] == 1 && sourceMap["numMeals"] >= 10 {
			return "집밥을 즐기는 미식가를 위한 추천"
		} else {
			return "우리 가족에게 딱 맞는 추천 구성"
		}
	}

	if size == "R+" {
		if sourceMap["haveYoung"] == 1 {
			return "성장기 아이와 함께 즐기는 추천 구성"
		} else if sourceMap["numFamily"] == 1 && sourceMap["numMeals"] >= 10 {
			return "집밥을 즐기는 미식가를 위한 추천"
		} else {
			return "꼭 필요한 것만 쏙쏙 담은 추천 구성"
		}
	}

	if size == "L" {
		return "온 가족 넉넉하게 즐기는 추천 구성"
	}

	return "온 가족 다양하고 풍성하게 즐기는 추천 구성"
}

// getItemCategoriesString : 아이템 카테고리 문자열 획득
func getItemCategoriesString(sourceMap map[string]int) string {
	itemCategories := []string{}
	if sourceMap["isWantedRice"] == 1 {
		itemCategories = append(itemCategories, "RICE")
	}
	if sourceMap["isWantedFruit"] == 1 {
		itemCategories = append(itemCategories, "FRUIT")
	}
	if sourceMap["isWantedVegetable"] == 1 {
		itemCategories = append(itemCategories, "VEGET")
	}
	if sourceMap["isWantedGrain"] == 1 {
		itemCategories = append(itemCategories, "GRAIN")
	}
	if sourceMap["isWantedQRice"] == 1 {
		itemCategories = append(itemCategories, "QRICE")
	}

	return strings.Join(itemCategories, "|")
}

// HandleSurveyResult : 설문조사 결과
func HandleSurveyResult(c *gin.Context) {
	var state = gin.H{}
	utils.CombineCustomStateGlobalState(c, &state)

	// 큐레이션 시작
	// tagGroupNo 획득
	isSuccess, cookieResult := utils.GetIDFromCookie(c, "tagGroupNo")
	if !isSuccess {
		state["errorMessage"] = "설문조사 정보가 없습니다. 설문조사를 완료 후 진행해주세요."
		c.HTML(http.StatusBadRequest, "error.html", state)
		return
	}
	tagGroupNo := cookieResult.(int)
	state["tagGroupNo"] = tagGroupNo

	isSuccess, tagsResult := GetOwnedTags(c, tagGroupNo)
	if !isSuccess {
		slackMsg := fmt.Sprintf("[front]HandleSurveyCurate::GetOwnedTags::%s", tagsResult.(string))
		utils.SendSlackMessage(utils.SlackChannel, slackMsg)
		state["errorMessage"] = "시스템 에러 발생"
		state["errorCode"] = "W0024"
		c.HTML(http.StatusBadRequest, "error.html", state)
		return
	}

	tags := tagsResult.([]models.Tag)
	isSuccess, sourceMapResult := getScoreSourceMap(tags)
	if !isSuccess {
		slackMsg := fmt.Sprintf("[front]HandleSurveyCurate::getScoreSourceMap(tags)::%s", sourceMapResult.(string))
		utils.SendSlackMessage(utils.SlackChannel, slackMsg)
		state["errorMessage"] = "시스템 에러 발생"
		state["errorCode"] = "W0025"
		c.HTML(http.StatusBadRequest, "error.html", state)
		return
	}
	sourceMap := sourceMapResult.(map[string]int)
	fmt.Println("sourceMap", sourceMap)
	totalScore := getScore(sourceMap)
	size := getSize(totalScore)
	itemCategoryString := getItemCategoriesString(sourceMap)
	phrase := getPhrase(size, sourceMap)

	// 큐레이션 결과 획득
	state["size"] = size
	state["itemCategoryExp"] = itemCategoryString
	state["phrase"] = phrase

	isSuccess, itemState := controllers.GetSubsItemsState(4, -12)
	if !isSuccess {
		slackMsg := fmt.Sprintf("[front]HandleSurveyResult::controllers.GetSubsItemsState(4, -12)::%s", itemState["error"])
		utils.SendSlackMessage(utils.SlackChannel, slackMsg)
		state["errorMessage"] = "시스템 에러 발생"
		state["errorCode"] = "W0010"
		c.HTML(http.StatusBadRequest, "error.html", state)
		return
	}
	utils.CombineTwoGinH(&state, &itemState)

	state["UserName"] = ""
	isSuccess, userResult := controllers.GetOwnedUser(c)
	if !isSuccess {
		slackMsg := fmt.Sprintf("[front]HandleSurveyResult::controllers.GetOwnedUser::%s", userResult.(string))
		utils.SendSlackMessage(utils.SlackChannel, slackMsg)
		state["errorMessage"] = "시스템 에러 발생"
		state["errorCode"] = "W0013"
		c.HTML(http.StatusBadRequest, "error.html", state)
		return
	}
	var user = *userResult.(*[]models.UserMst)
	if len(user) > 0 {
		state["UserName"] = user[0].UserName
	}

	state["navTitle"] = "설문조사 결과"

	c.HTML(http.StatusOK, "survey/result.html", state)
	return

}
