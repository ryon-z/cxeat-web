package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"yelloment-api/config"
	envutil "yelloment-api/env_util"
	"yelloment-api/middleware/mixin"
	"yelloment-api/models"
	"yelloment-api/utils"

	"github.com/gin-gonic/gin"
)

// SearchAddressResp : 도로명 주소 검색 결과 구조체
type SearchAddressResp struct {
	Results struct {
		Common struct {
			ErrorMessage string `json:"errorMessage"`
			CountPerPage string `json:"countPerPage"`
			TotalCount   string `json:"totalCount"`
			ErrorCode    string `json:"errorCode"`
			CurrentPage  string `json:"currentPage"`
		} `json:"common"`
		Juso []struct {
			DetBdNmList   string `json:"detBdNmList"`
			EngAddr       string `json:"engAddr"`
			Rn            string `json:"rn"`
			EmdNm         string `json:"emdNm"`
			ZipNo         string `json:"zipNo"`
			RoadAddrPart2 string `json:"roadAddrPart2"`
			EmdNo         string `json:"emdNo"`
			SggNm         string `json:"sggNm"`
			JibunAddr     string `json:"jibunAddr"`
			SiNm          string `json:"siNm"`
			RoadAddrPart1 string `json:"roadAddrPart1"`
			BdNm          string `json:"bdNm"`
			AdmCd         string `json:"admCd"`
			UdrtYn        string `json:"udrtYn"`
			LnbrMnnm      string `json:"lnbrMnnm"`
			RoadAddr      string `json:"roadAddr"`
			LnbrSlno      string `json:"lnbrSlno"`
			BuldMnnm      string `json:"buldMnnm"`
			BdKdcd        string `json:"bdKdcd"`
			LiNm          string `json:"liNm"`
			RnMgtSn       string `json:"rnMgtSn"`
			MtYn          string `json:"mtYn"`
			BdMgtSn       string `json:"bdMgtSn"`
			BuldSlno      string `json:"buldSlno"`
		} `json:"juso"`
	} `json:"results"`
}

// SearchAddress : 주소 검색
func SearchAddress(c *gin.Context) {
	keyword := c.Query("keyword")
	if keyword == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "keyword is empty"})
		return
	}
	utils.CheckSQLInjection(keyword)

	baseURL := "https://www.juso.go.kr/addrlink/addrLinkApi.do"
	data := url.Values{}
	data.Set("confmKey", envutil.GetGoDotEnvVariable("ADDRESS_KEY"))
	data.Set("currentPage", "1")
	data.Set("countPerPage", "5")
	data.Set("keyword", keyword)
	data.Set("resultType", "json")

	client := &http.Client{}
	req, err := http.NewRequest("GET", baseURL, strings.NewReader(data.Encode()))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Println(string(body))

	var searchAddressResp SearchAddressResp
	err = json.Unmarshal(body, &searchAddressResp)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if searchAddressResp.Results.Common.ErrorCode != "0" {
		c.JSON(http.StatusBadRequest, gin.H{"error": searchAddressResp.Results.Common.ErrorMessage})
		return
	}

	type addrResultRow struct {
		JibunAddr string `json:"JibunAddr"`
		RoadAddr  string `json:"RoadAddr"`
		ZipNo     string `json:"ZipNo"`
	}
	result := []addrResultRow{}

	for _, addrStruct := range searchAddressResp.Results.Juso {
		row := addrResultRow{JibunAddr: addrStruct.JibunAddr, RoadAddr: addrStruct.RoadAddr, ZipNo: addrStruct.ZipNo}
		result = append(result, row)
	}

	c.JSON(http.StatusOK, gin.H{"data": result})
}

// GetOwnedUserAddressUnit : 소유자 유저 주소 행 조회(주소ID는 쿠키)
func GetOwnedUserAddressUnit(c *gin.Context, addressNo int) (bool, interface{}) {
	userNo := utils.GetUserNo(c)
	wherePhrase := fmt.Sprintf(`UserNo = %d AND AddressNo = %d AND StatusCode = "normal"`, userNo, addressNo)
	isSuccess, result := mixin.List("userAddress", wherePhrase, "")
	if !isSuccess {
		return false, result.(string)
	}

	return true, result
}

// GetOwnedUserAddress : 소유자 주소 리스트 획득
func GetOwnedUserAddress(c *gin.Context) (bool, interface{}) {
	userNo := utils.GetUserNo(c)
	wherePhrase := fmt.Sprintf(`UserNo = %d AND StatusCode = "normal"`, userNo)
	isSuccess, result := mixin.List("userAddress", wherePhrase, "RegDate DESC")
	return isSuccess, result
}

// GetOwnedBasicUserAddress : 소유자 기본 주소 획득
func GetOwnedBasicUserAddress(c *gin.Context) (bool, interface{}) {
	userNo := utils.GetUserNo(c)
	wherePhrase := fmt.Sprintf(`UserNo = %d AND IsBasic = 1 AND StatusCode = "normal"`, userNo)
	isSuccess, result := mixin.List("userAddress", wherePhrase, "")
	return isSuccess, result
}

// ListOwnedUserAddress : 소유자 주소 리스트 조회
func ListOwnedUserAddress(c *gin.Context) {
	isSuccess, result := GetOwnedUserAddress(c)
	utils.GetHTTPResponse(isSuccess, result, c)
}

// RetrieveOwnedUserAddress : 소유자 주소 단일 조회
func RetrieveOwnedUserAddress(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userNo := utils.GetUserNo(c)

	wherePhrase := fmt.Sprintf("UserNo = %d AND AddressNo = %d AND StatusCode = 'normal'", userNo, idInt)
	isSuccess, result := mixin.Retrieve("userAddress", wherePhrase, "addressNo", idInt)
	utils.GetHTTPResponse(isSuccess, result, c)
}

// CreateOwnedUserAddress : 소유자 주소 생성
func CreateOwnedUserAddress(c *gin.Context) {
	userNo := utils.GetUserNo(c)

	var address models.UserAddress
	if err := c.ShouldBindJSON(&address); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	address.StatusCode = "normal"
	isSuccess, result := mixin.CreateOwned("userAddress",
		userNo, &address, []string{"AddressNo", "RegDate"}, true, []string{})
	utils.GetHTTPResponse(isSuccess, result, c)
}

// PartialUpdateOwnedUserAddress : 소유자 주소 수정
func PartialUpdateOwnedUserAddress(c *gin.Context) {
	userNo := utils.GetUserNo(c)

	var addr models.UserAddress
	if err := c.ShouldBindJSON(&addr); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if addr.StatusCode != "" && !utils.StringInSlice(addr.StatusCode, config.AddressStatusCodes) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "StatusCode is not allowed"})
		return
	}

	wherePhrase := fmt.Sprintf("UserNo = %d AND addressNo = %d", userNo, addr.AddressNo)
	isSuccess, message := mixin.PartialUpdate("userAddress", wherePhrase, &addr, []string{"AddressNo", "RegDate"}, []string{"IsBasic"})
	utils.GetHTTPResponse(isSuccess, message, c)
}

// DeactivateOwnedUserAddress : 주소 비활성화
func DeactivateOwnedUserAddress(c *gin.Context) {
	userNo := utils.GetUserNo(c)

	var inputAddress models.UserAddress
	if err := c.ShouldBindJSON(&inputAddress); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 필수 파라미터 체크
	if inputAddress.AddressNo == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "AddressNo is not contained"})
		return
	}

	// 만약 IsBasic이 1이고 또 다른 주소가 존재한다면 또 다른 주소의 IsBasic 값을 1로 변경
	if inputAddress.IsBasic == 1 {
		isSuccess, addressResult := GetOwnedUserAddress(c)
		if !isSuccess {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed loading API"})
			return
		}

		var userAddresses = *addressResult.(*[]models.UserAddress)
		if len(userAddresses) >= 2 {
			for _, userAddress := range userAddresses {
				if userAddress.AddressNo != inputAddress.AddressNo {
					var otherAddress models.UserAddress
					otherAddress.UserNo = userNo
					otherAddress.AddressNo = userAddress.AddressNo
					otherAddress.IsBasic = 1
					wherePhrase := fmt.Sprintf("UserNo = %d AND addressNo = %d", userNo, userAddress.AddressNo)
					isSuccess, message := mixin.PartialUpdate("userAddress", wherePhrase, &otherAddress, []string{"RegDate"}, []string{"IsBasic"})
					if !isSuccess {
						c.JSON(http.StatusBadRequest, gin.H{"error": message})
						return
					}
					break
				}
			}
		}
	}

	var targetAddress models.UserAddress
	targetAddress.StatusCode = "cancel"
	targetAddress.IsBasic = 0

	wherePhrase := fmt.Sprintf(`UserNo = %d AND AddressNo = %d`, userNo, inputAddress.AddressNo)
	isSuccess, message := mixin.PartialUpdate("userAddress", wherePhrase, &targetAddress, []string{"RegDate"}, []string{"IsBasic"})
	utils.GetHTTPResponse(isSuccess, message, c)
}
