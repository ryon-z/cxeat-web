package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	envutil "yelloment-api/env_util"
	"yelloment-api/middleware/mixin"
	"yelloment-api/models"
	"yelloment-api/utils"

	"github.com/gin-gonic/gin"
)

// CardInput : 카드키 요청 시 사용하는 정보 구조체
type CardInput struct {
	CardNickName string `json:"CardNickName"`
	CardNumber   string `json:"CardNumber"`
	Birth        string `json:"Birth"`
	Expiry       string `json:"Expiry"`
	Pwd2Digit    string `json:"Pwd2Digit"`
	IsBasic      int    `json:"IsBasic"`
}

// IamportAuthResp : Iamport auth token 요청 결과 구조체
type IamportAuthResp struct {
	Code     int    `json:"code"`
	Message  string `json:"message"`
	Response struct {
		AccessToken string `json:"access_token"`
		Now         int    `json:"now"`
		ExpiredAt   int    `json:"expired_at"`
	} `json:"response"`
}

// IamportCardRegiResp : 카드 등록 요청 결과 구조체
type IamportCardRegiResp struct {
	Code     int         `json:"code"`
	Message  interface{} `json:"message"`
	Response struct {
		CardCode         string      `json:"card_code"`
		CardName         string      `json:"card_name"`
		CardNumber       string      `json:"card_number"`
		CardType         int         `json:"card_type"`
		CustomerAddr     interface{} `json:"customer_addr"`
		CustomerEmail    interface{} `json:"customer_email"`
		CustomerName     interface{} `json:"customer_name"`
		CustomerPostcode interface{} `json:"customer_postcode"`
		CustomerTel      interface{} `json:"customer_tel"`
		CustomerUID      string      `json:"customer_uid"`
		Inserted         int         `json:"inserted"`
		PgID             string      `json:"pg_id"`
		PgProvider       string      `json:"pg_provider"`
		Updated          int         `json:"updated"`
	} `json:"response"`
}

// checkCardInputValid : cardInput 유효성 체크
func checkCardInputValid(cardInput CardInput) (isSuccess bool, message string) {
	// Check cardNumber
	if len(cardInput.CardNumber) > 16 || len(cardInput.CardNumber) < 14 {
		return false, "CardNumber is not valid"
	}
	if _, err := strconv.Atoi(cardInput.CardNumber); err != nil {
		return false, err.Error()
	}

	// Check Birth
	if _, err := strconv.Atoi(cardInput.Birth); err != nil {
		return false, err.Error()
	}

	// Check Expiry
	splited := strings.Split(cardInput.Expiry, "-")
	if len(splited) != 2 {
		return false, "Expiry is not valid"
	}
	for i, value := range splited {
		if i == 0 && len(value) != 4 {
			return false, "Expiry is not valid"
		}
		if i == 1 && len(value) != 2 {
			return false, "Expiry is not valid"
		}
		if _, err := strconv.Atoi(value); err != nil {
			return false, err.Error()
		}
	}

	// Check IsBasic
	if cardInput.IsBasic != 0 && cardInput.IsBasic != 1 {
		return false, "IsBasic allowed only 0 or 1"
	}

	return true, "success"
}

// GetOwnedUserCard : 소유자 카드 리스트  획득
func GetOwnedUserCard(c *gin.Context) (bool, interface{}) {
	userNo := utils.GetUserNo(c)
	wherePhrase := fmt.Sprintf(`UserNo = %d AND StatusCode = "normal"`, userNo)
	return mixin.List("userCard", wherePhrase, "")
}

// GetOwnedBasicUserCard : 기본 소유자 카드 리스트 획득
func GetOwnedBasicUserCard(c *gin.Context) (bool, interface{}) {
	userNo := utils.GetUserNo(c)
	wherePhrase := fmt.Sprintf(`UserNo = %d AND StatusCode = "normal" AND IsBasic = 1`, userNo)
	return mixin.List("userCard", wherePhrase, "")
}

// ListOwnedUserCard : 소유자 카드 리스트 조회
func ListOwnedUserCard(c *gin.Context) {
	isSuccess, result := GetOwnedUserCard(c)
	utils.GetHTTPResponse(isSuccess, result, c)
}

// RetrieveOwnedUserCard : 소유자 카드 단일 조회
func RetrieveOwnedUserCard(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userNo := utils.GetUserNo(c)

	wherePhrase := fmt.Sprintf("UserNo = %d AND cardRegNo = %d AND StatusCode = 'normal'", userNo, idInt)
	isSuccess, result := mixin.Retrieve("userCard", wherePhrase, "cardRegNo", idInt)
	// Client-side에서 ajax call을 허용하기 때문에 민감한 CardKey는 지우고 전송
	if isSuccess {
		var userCard = *result.(*[]models.UserCard)
		if len(userCard) >= 1 {
			userCard[0].CardKey = ""
		}
		c.JSON(http.StatusOK, gin.H{"data": userCard})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": result})
	}
}

// getIamportAuthToken : iamport auth token 획득
func getIamportAuthToken() (isSuccess bool, result string) {
	isSuccess, reqResult := utils.CustomRequest(
		"https://api.iamport.kr/users/getToken",
		"post",
		[][]string{
			{"imp_key", envutil.GetGoDotEnvVariable("IAMPORT_KEY")},
			{"imp_secret", envutil.GetGoDotEnvVariable("IAMPORT_SECRET")}},
		[][]string{{"Content-Type", "application/x-www-form-urlencoded"}},
		[]*http.Cookie{},
	)
	if !isSuccess {
		return false, reqResult
	}

	var iamportAuthResp IamportAuthResp
	if err := json.Unmarshal([]byte(reqResult), &iamportAuthResp); err != nil {
		return false, err.Error()
	}
	if iamportAuthResp.Code != 0 {
		return false, iamportAuthResp.Message
	}

	iamportAuthToken := iamportAuthResp.Response.AccessToken
	return true, iamportAuthToken
}

// isDuplicatedCard : 중복 카드 여부
func isDuplicatedCard(userNo int, cardNumber string) (isSuccess bool, isDuplicated bool, message string) {
	masked := cardNumber[0:8] + "****" + cardNumber[12:]

	wherePhrase := fmt.Sprintf("UserNo = %d AND CardNumber = '%s' AND StatusCode = 'normal'", userNo, masked)
	isSuccess, result := mixin.List("userCard", wherePhrase, "")
	if !isSuccess {
		return false, false, "faild looking for CardNumber"
	}

	userCards := *result.(*[]models.UserCard)
	if len(userCards) > 0 {
		return true, true, "this card is duplicated"
	}

	return true, false, "success"
}

// CreateOwnedUserCard : 소유자 카드 생성
func CreateOwnedUserCard(c *gin.Context) {
	userNo := utils.GetUserNo(c)

	var cardInput CardInput
	if err := c.ShouldBindJSON(&cardInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cardInput.CardNumber = strings.ReplaceAll(cardInput.CardNumber, "-", "")
	if isSuccess, message := checkCardInputValid(cardInput); !isSuccess {
		c.JSON(http.StatusBadRequest, gin.H{"error": message})
		return
	}
	fmt.Println("cardInput.CardNumber", cardInput.CardNumber)

	// Registrate Card
	isSuccess, result := getIamportAuthToken()
	if !isSuccess {
		c.JSON(http.StatusBadRequest, gin.H{"error": result})
		return
	}

	// Check Card is duplicated
	isSuccess, isDuplicated, message := isDuplicatedCard(userNo, cardInput.CardNumber)
	if !isSuccess {
		c.JSON(http.StatusBadRequest, gin.H{"error": message})
		return
	}
	if isDuplicated {
		c.JSON(http.StatusBadRequest, gin.H{"data": message})
		return
	}

	iamportAuthToken := result
	customerUID := fmt.Sprintf("%s_%d", utils.GetRandString(10), userNo)
	isSuccess, reqResult := utils.CustomRequest(
		"https://api.iamport.kr/subscribe/customers/"+customerUID,
		"post",
		[][]string{
			{"card_number", cardInput.CardNumber},
			{"expiry", cardInput.Expiry},
			{"birth", cardInput.Birth},
			{"pwd_2digit", cardInput.Pwd2Digit},
			{"customer_uid", customerUID}},
		[][]string{
			{"Authorization", "Bearer " + iamportAuthToken},
			{"Accept", "application/json"},
			{"Content-Type", "application/x-www-form-urlencoded"}},
		[]*http.Cookie{},
	)
	if !isSuccess {
		c.JSON(http.StatusBadRequest, gin.H{"error": "It is failed to registrate card."})
		return
	}

	var iamportCardRegiResp IamportCardRegiResp
	if err := json.Unmarshal([]byte(reqResult), &iamportCardRegiResp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if iamportCardRegiResp.Code != 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": iamportCardRegiResp.Message})
		return

	}

	// nullable Fields
	cardName := iamportCardRegiResp.Response.CardName
	cardNumber := iamportCardRegiResp.Response.CardNumber
	cardNickName := cardInput.CardNickName

	// Insert Card information to DB
	var card models.UserCard
	fmt.Println("iamportCardRegiResp", iamportCardRegiResp)
	card.UserNo = userNo
	card.CardName = &cardName
	card.CardNickName = &cardNickName
	card.CardNumber = &cardNumber
	card.CardKey = customerUID
	card.StatusCode = "normal"
	card.IsBasic = cardInput.IsBasic
	isSuccess, createResult := mixin.CreateOwned("userCard", userNo, &card, []string{"CardRegNo", "RegDate"}, true, []string{"IsBasic"})
	utils.GetHTTPResponse(isSuccess, createResult, c)
}

// DeactivateOwndUserCard : 소유자 카드 비활성화
func DeactivateOwndUserCard(c *gin.Context) {
	userNo := utils.GetUserNo(c)

	var inputCard models.UserCard
	if err := c.ShouldBindJSON(&inputCard); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Println("inputCard", inputCard)

	if inputCard.CardRegNo == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "CardRegNo is not contained"})
		return
	}

	// 만약 IsBasic이 1이고 또 다른 주소가 존재한다면 또 다른 주소의 IsBasic 값을 1로 변경
	if inputCard.IsBasic == 1 {
		isSuccess, cardResult := GetOwnedUserCard(c)
		if !isSuccess {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed loading API"})
			return
		}

		var userCards = *cardResult.(*[]models.UserCard)
		if len(userCards) >= 2 {
			for _, userCard := range userCards {
				if userCard.CardRegNo != inputCard.CardRegNo {
					var otherCard models.UserCard
					otherCard.UserNo = userNo
					otherCard.CardRegNo = userCard.CardRegNo
					otherCard.IsBasic = 1
					wherePhrase := fmt.Sprintf("UserNo = %d AND CardRegNo = %d", userNo, userCard.CardRegNo)
					isSuccess, message := mixin.PartialUpdate("userCard", wherePhrase, &otherCard, []string{"RegDate"}, []string{"IsBasic"})
					if !isSuccess {
						c.JSON(http.StatusBadRequest, gin.H{"error": message})
						return
					}
					break
				}
			}
		}
	}

	var targetCard models.UserCard
	targetCard.StatusCode = "cancel"
	targetCard.IsBasic = 0

	wherePhrase := fmt.Sprintf(`UserNo = %d AND CardRegNo = %d`, userNo, inputCard.CardRegNo)
	isSuccess, message := mixin.PartialUpdate("userCard", wherePhrase, &targetCard, []string{"RegDate"}, []string{"IsBasic"})
	utils.GetHTTPResponse(isSuccess, message, c)
}
