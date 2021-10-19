package viewcontrollers

import (
	"fmt"
	"net/http"
	"yelloment-api/database"
	"yelloment-api/models"
	"yelloment-api/utils"

	"github.com/gin-gonic/gin"
)

// PhotoReview : 포토 리뷰
type PhotoReview struct {
	ImagePath  string
	Title      string
	Round      int
	BoxSize    string
	Products   []string
	TextReview string
}

// GetUserCount : 유저 수 획득
func GetUserCount() (isSuccess bool, result interface{}) {
	var queryResult int
	tablename := models.GetTableName("userMst")
	wherePhrase := `1 = 1`
	sqlQuery := fmt.Sprintf(`
		SELECT count(*) AS CNT
		FROM %s
		WHERE (%s)
		LIMIT 1
	`, tablename, wherePhrase)

	err := database.DB.Get(&queryResult, sqlQuery)
	if err != nil {
		return false, err.Error()
	}

	return true, queryResult
}

// HandleMain : main 뷰
func HandleMain(c *gin.Context) {
	state := gin.H{}
	utils.CombineCustomStateGlobalState(c, &state)

	// 유저 수 확보
	state["userCount"] = 4000
	isSuccess, result := GetUserCount()
	if isSuccess {
		state["userCount"] = result.(int)
	} else {
		slackMsg := fmt.Sprintf("[front]HandleMain::webcontrollers.GetUserCount::%s", result.(string))
		utils.SendSlackMessage(utils.SlackChannel, slackMsg)
	}

	state["isMain"] = "yes"
	state["photoReview"] = []PhotoReview{
		{"/contents/images/mainReviewImg_01.jpg", "30대, 여, 3인가구, 매주", 2, "레귤러+", []string{"#쌀", "#과일", "#채소"}, "일하랴 애들 키우랴, 몸이 두개여도 모자란 워킹맘이에요. 매주 장보는 것도 저한테는 일이었는데 우리가족에게 필요한 쌀, 과일, 야채들을 알아서 보내주니 너무 편하네요."},
		{"/contents/images/mainReviewImg_02.jpg", "40대, 여, 3인가구, 매주", 2, "레귤러", []string{"#과일", "#채소"}, "항상 먹던 것만 먹고 살았는데 큐잇 덕분에 새로운 입맛을 찾은 것 같아요. 남편도 아들도 처음 맛보는 과일들이 신기한지 다음주에 뭐가 오냐고 계속 물어보네요."},
		{"/contents/images/mainReviewImg_03.jpg", "60대, 여, 1인가구, 격주", 4, "싱글", []string{"#쌀", "#잡곡", "#과일", "#야채"}, "지방에 혼자 계신 어머니가 잘 챙겨 드시는지 항상 마음에 걸렸는데 큐잇 덕분에 마음이 좀 놓이네요. 특히 과일을 좋아하시는데 전부 맛있다고 하셔서 안심이에요."},
		{"/contents/images/mainReviewImg_04.jpg", "30대, 남, 3인가구, 격주", 7, "싱글", []string{"#즉석밥", "#과일", "#채소"}, "항상 마트에 가면 대용량만 팔아서 몇 가지 과일을 오래 두고 먹었는데 더 저렴한 비용으로 다양하게 맛볼 수 있어서 좋아요."},
		{"/contents/images/mainReviewImg_05.jpg", "40대, 여, 4인가구, 매주", 3, "싱글", []string{"#과일", "#채소"}, "정기구독한 이후로 오히려 장보는 비용이 줄었어요. 그 동안 다 먹지도 못할 식재료들을 습관처럼 얼마나 냉장고에 쟁여 두고 있었는지 후회가 되네요. 감사합니다."},
		{"/contents/images/mainReviewImg_06.jpg", "30대, 여, 3인가구, 매주", 1, "싱글+", []string{"#쌀", "#과일", "#채소"}, "과일을 정말 좋아하지만 맛없는 건 먹지 않는 까다로운 28개월 아이가 집에 있습니다ㅋㅋㅋ 그런데 전반적으로 과일을 정말 잘 먹어서 놀랐어요. 체리, 블루베리, 심지어 바나 나도 맛없다 느끼면 안 먹는데 끝까지 다 먹더라구요~ 제 입맛엔 스테비아 방울 토마토도 굉장히 맛있었어요! 아쉬운 점은 애플수박 안쪽에 약간 과질이 무른 부분이 있었다는 거에요. 처음으로 받아봤는데 만족스러웠어요~"},
		{"/contents/images/mainReviewImg_07.jpg", "30대, 여, 4인가구, 매주", 4, "싱글+", []string{"#쌀", "#과일", "#채소"}, "갓난아이 키우면서 초등학생도 키우려니 장보기 힘든 요즘. 다양한 과일을 배송해주니 과일편식 없이 다양하게 먹을 수 있다. 이번에 양구 노란 멜론이나 미니 수박 같은건 집 앞 과일가게에선 못 본 거라 큐잇 구독으로 먹어보니 좋음.특히 애플수박은 맛도 있고 양도 딱 한번 먹기적당해서 보관 걱정, 상할 걱정 같은 거 안 하는 것도 완전 👍"},
		{"/contents/images/mainReviewImg_08.jpg", "30대, 여, 4인가구, 매주", 3, "싱글+", []string{"#쌀", "#과일", "#채소"}, "쌀, 과일, 채소 신청했는데 정해진 요일에 딱 일주일치가 오네요. 신기하게 은근히 양이 딱 맞아요 ㅎㅎ매주 알아서 오는 게 생각보다 편하더라구요. 먹을 만큼만 오니까 음식물 쓰레기 양도 줄고 장 보는 번거로움도 줄고,, 맞벌이 부부에게 강추에요. 스테비아 방토는 처음 먹어봤는데 진짜 달고 맛있어요! 와 진짜사탕 먹는 줄~~~~~ 과일 취향을 새콤보다는 달콤으로 조절했더니 진짜 그렇게 왔어요. 신기방기~ 멀리 계신 부모님께도 보내드리는 것도 좋을 것 같아요. "},
		{"/contents/images/mainReviewImg_09.jpg", "30대, 여, 2인가구, 매주", 2, "싱글", []string{"#과일"}, "매주 싱싱한 과일들이 이것저것 와서 좋아요. 와인을 즐겨 마시는데 와인이랑 어울리는 과일들이 자꾸 와서 와인을 더 사게 되네요………ㅋㅋ 다음 박스 오기 전까지 아껴 먹어야 됩니다 ㅎㅎ"},
		{"/contents/images/mainReviewImg_10.jpg", "40대, 여, 4인가구, 매주", 1, "레귤러", []string{"#쌀", "#과일", "#채소"}, "최근 써 본 정기구독 서비스 중 은근히 유용해요! 라이프스타일에 맞춰서 적정량을 보내주는데 저희집식구 여자 4명이 한번씩 먹기에 딱좋네요~ 근데 다 그렇다 쳐도 진짜맛있어요. 초당옥수수랑 체리, 스테비아 토마토 다 엄청 맛있어서 놀랐다는 ~~~~~ 우리집 여자들 과일킬러라서 과일 많이 먹는데 전부 다 합격이네요ㅎㅎ"},
		{"/contents/images/mainReviewImg_11.jpg", "30대, 여, 4인가구, 매주", 1, "레귤러+", []string{"쌀", "#과일"}, "저희 아이도 밥을 잘 안 먹어서 늘 걱정인데 큐잇에 들어있는 과일은 다 배송 온 첫날 순삭이에요. 아이가 좋아하는 과일 위주로 들어있는 것 같아요! 애플수박도 보더니 아기 수박이냐면서 넘 좋아하면서 먹네요. 아이 때문에 계속 잘 이용하게 될 것 같아요~ "},
		{"/contents/images/mainReviewImg_12.jpg", "30대, 여, 4인가구, 매주", 1, "싱글+", []string{"쌀", "#과일"}, "같은 배속에서 나왔지만 어쩜 이렇게 먹는 취향이 다른지..애들부터 남편까지 식성이 가지각색인데 뭐 하나 사려면 양이 많아서 결국 나혼자 다 먹거나 아쉽게도 곰팡이 펴서 버리는 게 반이었는데 우리 가족이 다양하게 먹을 수 있게 여러가지로 오니까 정말 좋아요. 요즘 마트 가는 것도 부담스러운데 정기적으로 오니까 은근 제 시간이 많이 남네요?? 생각보다 생활이 편해져서 앞으로 쭉 애용할 느낌~❤"},
		{"/contents/images/mainReviewImg_13.jpg", "20대, 여, 1인가구, 매주", 1, "싱글", []string{"#쌀", "#과일"}, "과일 좋아해서 랜덤박스 많이 이용해 봤는데 다른 데보다 훨씬 실속 있어요. 내가 필요했던 것들, 시도해 보고 싶었던 것들만 소량으로 들어있어서 좋음#가성비템 #1인가구추천"},
	}

	c.HTML(http.StatusOK, "main.html", state)
}
