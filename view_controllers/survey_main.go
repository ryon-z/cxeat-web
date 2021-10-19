package viewcontrollers

import (
	"fmt"
	"net/http"
	"yelloment-api/utils"

	"github.com/gin-gonic/gin"
)

// SurveyType : 설문조사 유형
type SurveyType struct {
	Name string
	Type string
}

// SurveyInput : 설문조사 인풋 데이터
type SurveyInput struct {
	ExposedLabel string
	Label        string
	Value        string
	CustomClass  string
}

// HandleSurveyMain : 설문조사 본문
func HandleSurveyMain(c *gin.Context) {
	var state = gin.H{}
	utils.CombineCustomStateGlobalState(c, &state)

	isSuccess, state := utils.GetGlobalState(c, "newLineSymbol", state)
	if !isSuccess {
		state["errorMessage"] = "시스템 에러 발생"
		state["errorCode"] = "W0030"
		c.HTML(http.StatusBadRequest, "error.html", state)
		return
	}

	state["typeInfo"] = []SurveyType{
		{"numFamily", "FAMILY"},
		{"healthInfo", "HEALTH"},
		{"taste", "FRUIT_TASTE"},
		{"typeFruits", "FRUIT"},
		{"hateFruits", "FRUIT_HATE"},
		{"alergyFruits", "FRUIT_HATE"},
		{"preferenceFruits", "FRUIT_REQ"},
		{"typeMeal", "RICE"},
		{"numMeal", "MEAL_CNT"},
		{"preferenceMeal", "RICE_LIKE"},
		{"typeVegetable", "VEGET"},
		{"preferenceVegetable", "VEGET_REQ"},
		{"preferenceCategoryVegetable", "VEGET_LIKE"},
		{"hateVegetable", "VEGET_HATE"},
		{"alergyVegetable", "VEGET_HATE"},
		{"newAgriFood", "LIFE_STYLE"},
		{"lifeStyle", "LIFE_STYLE"},
	}
	state["healthInfo"] = []SurveyInput{
		{"체중관리", "체중관리", "", ""},
		{"고혈압 및 당뇨", "고혈압 및 당뇨", "", ""},
		{"소화", "소화", "", ""},
		{"배변", "배변", "", ""},
		{"탈모", "탈모", "", ""},
		{"갱년기", "갱년기", "", ""},
	}
	state["typeFruits"] = []SurveyInput{
		{"집에서 과일을 즐겨 먹어요", "모든과일", "", ""},
		{"과일을 좋아하지만 손질해서 먹는 건 귀찮아요", "간편과일", "", ""},
		{"과일을 별로 좋아하지 않아요", "과일X", "", ""},
	}
	state["preferenceFruits"] = []SurveyInput{
		{"사과", "사과", "", ""},
		{"바나나", "바나나", "", ""},
		{"방울토마토", "방울토마토", "", "survey_long_text"},
		{"토마토", "토마토", "", ""},
	}
	state["hateFruits"] = []SurveyInput{
		{"사과", "사과", "", ""},
		{"배", "배", "", ""},
		{"귤", "귤", "", ""},
		{"복숭아", "복숭아", "", ""},
		{"포도", "포도", "", ""},
		{"딸기", "딸기", "", ""},
		{"수박", "수박", "", ""},
		{"참외", "참외", "", ""},
		{"자두", "자두", "", ""},
		{"단감", "단감", "", ""},
		{"홍시", "홍시", "", ""},
		{"방울토마토", "방울토마토", "", "survey_long_text"},
		{"석류", "석류", "", ""},
		{"토마토", "토마토", "", ""},
		{"대추", "대추", "", ""},
		{"오렌지", "오렌지", "", ""},
		{"바나나", "바나나", "", ""},
		{"파인애플", "파인애플", "", ""},
		{"키위", "키위", "", ""},
		{"망고", "망고", "", ""},
		{"아보카도", "아보카도", "", ""},
		{"블루베리", "블루베리", "", ""},
		{"체리", "체리", "", ""},
		{"자몽", "자몽", "", ""},
		{"멜론", "멜론", "", ""},
	}
	state["typeMeal"] = []SurveyInput{
		{"대부분 밥을 지어 먹어요", "취식", "", ""},
		{"즉석밥을 선호해요", "즉석밥", "", ""},
		{"대부분 외식, 배달로 해결해요", "취식X", "", ""},
	}
	state["numMeal"] = []SurveyInput{
		{"1~3끼", "1~3끼", "2", ""},
		{"4~6끼", "4~6끼", "5", ""},
		{"7~9끼", "7~9끼", "8", ""},
		{"10끼 이상", "10끼 이상", "10", ""},
	}
	state["preferenceMeal"] = []SurveyInput{
		{"하얀 쌀밥이 좋아요", "백미", "", ""},
		{"현미밥이 좋아요", "현미", "", ""},
		{"잡곡밥이 좋아요", "잡곡", "", ""},
	}
	state["typeVegetable"] = []SurveyInput{
		{"큐잇이 매주 추천하는 제철 채소를 받고 싶어요", "채소", "", ""},
		{"채소는 필요하지 않아요", "채소X", "", ""},
	}
	state["preferenceVegetable"] = []SurveyInput{
		{"파", "파", "", ""},
		{"마늘", "마늘", "", ""},
		{"양파", "양파", "", ""},
		{"고추", "고추", "", ""},
		{"청양고추", "청양고추", "", ""},
		{"감자", "감자", "", ""},
	}
	state["preferenceCategoryVegetable"] = []SurveyInput{
		{fmt.Sprintf("생채소(브로컬리, 양배추, 두릅, 고추 등)%s특별한 양념이나 요리 없이 소스 등만 곁들여 생으로 섭취", state["newLineSymbol"]), "생채소", "", ""},
		{fmt.Sprintf("조리용 채소(나물류)%s삶고, 볶는 정도의 간단한 과정이 필요한 채소", state["newLineSymbol"]), "조리용채소", "", ""},
		{fmt.Sprintf("요리용 채소(오이, 당근, 호박, 감자 등)%s국, 찌개, 반찬 등 다양한 요리에 활용할 채소", state["newLineSymbol"]), "요리용채소", "", ""},
		{fmt.Sprintf("쌈채소 (상추, 깻잎, 호박잎 등)%s쌈 요리에 활용할 수 있는 채소", state["newLineSymbol"]), "쌈채소", "", ""},
		{fmt.Sprintf("고기 곁들임 채소 (버섯, 아스파라거스 등 )%s삼겹살, 스테이크 등에 어울리는 채소", state["newLineSymbol"]), "고기곁들임채소", "", ""},
		{fmt.Sprintf("특별 채소 (남해 보물초, 노루궁뎅이 버섯 등)%s평소에 흔히 접하는 채소가 아닌 큐잇 추천 상품으로%s만나볼 수 있는 지역 특산품이나 기획 상품", state["newLineSymbol"], state["newLineSymbol"]), "특별채소", "", ""},
	}
	state["hateVegetable"] = []SurveyInput{
		{"파", "파", "", ""},
		{"마늘", "마늘", "", ""},
		{"양파", "양파", "", ""},
		{"고추", "고추", "", ""},
		{"청양고추", "청양고추", "", ""},
		{"오이", "오이", "", ""},
		{"당근", "당근", "", ""},
		{"가지", "가지", "", ""},
		{"감자", "감자", "", ""},
		{"고구마", "고구마", "", ""},
		{"미나리", "미나리", "", ""},
		{"시금치", "시금치", "", ""},
		{"버섯류", "버섯류", "", ""},
		{"양배추", "양배추", "", ""},
		{"깻잎", "깻잎", "", ""},
	}
	state["newAgriFood"] = []SurveyInput{
		{"새로 나온 과일이나 식품은 꼭 구매해보는 편이에요", "신제품_꼭구매", "", ""},
		{"다른 사람들의 반응이 좋으면 사요", "신제품_반응체크구매", "", ""},
		{"항상 먹던 것만 먹는 편이에요", "신제품_먹던것구매", "", ""},
	}
	state["lifeStyle"] = []SurveyInput{
		{"직접 요리를 해요.", "직접요리", "", ""},
		{"요리를 해먹기보다는 간편하게 먹는 게 좋아요", "간편식", "", ""},
		{"씨가 있거나 깎아 먹는 과일은 손이 가지 않아요", "적정구매선호", "", ""},
		{"과일, 쌀, 채소 등 모든 식품은 단가가 비싸더라도 필요한 만큼만 사는 게 좋아요", "정량선호", "", ""},
		{"과일, 쌀, 채소 등 모든 식품은 남더라도 저렴하다면 많이 사는 게 좋아요", "저렴선호", "", ""},
		{"새로운 먹거리에 관심이 많아요", "신먹거리관심", "", ""},
		{"TV건강프로그램에 나오는 건강 식품은 먹어보고 싶어져요", "건강TV선호", "", ""},
		{"건강을 위해 다양한 즙을 즐겨먹어요", "즙선호", "", ""},
		{"장은 조금씩 자주 보는게 좋아요", "자주쇼핑", "", ""},
		{"견과류를 좋아해요", "견과류선호", "", ""},
		{"고구마, 옥수수, 단호박을 좋아해요", "고옥단선호", "", ""},
		{"아이가 잘 먹으면 가격이 비싸도 얼마든지 사줄 수 있어요", "아이먹거리", "", ""},
		{"조금 비싸더라도 유기농 또는 친환경 식품을 구매하는 편이에요", "친환경식품선호", "", ""},
		{"환경에 도움이 된다면 비싼 가격 또는  불편함을 감수할 수 있어요", "친환경소비", "", ""},
		{"고급스럽고 보기 좋은 포장이 좋아요", "고급포장", "", ""},
		{"소박한 친환경 포장이 좋아요", "친환경포장", "", ""},
	}

	state["navTitle"] = "설문조사"

	c.HTML(http.StatusOK, "survey/main.html", state)

}
