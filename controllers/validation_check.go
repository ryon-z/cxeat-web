package controllers

import (
	"fmt"
	"yelloment-api/middleware/mixin"
	"yelloment-api/models"
)

// CheckAddressNoValid : 주소 번호 유효성 체크
func CheckAddressNoValid(userNo int, addressNo int) (bool, string) {
	wherePhrase := fmt.Sprintf(`UserNo = %d AND AddressNo = %d`, userNo, addressNo)
	isSuccess, result := mixin.List("userAddress", wherePhrase, "")
	if !isSuccess {
		return false, result.(string)
	}

	if len(*result.(*[]models.UserAddress)) < 1 {
		return false, fmt.Sprintf("AddressNo is not exists, AddressNo: %d", userNo)
	}

	return true, "success"
}

// CheckCardRegNoValid : 카드 동록 번호 유효성 체크
func CheckCardRegNoValid(userNo int, cardRegNo int) (bool, string) {
	wherePhrase := fmt.Sprintf(`UserNo = %d AND CardRegNo = %d`, userNo, cardRegNo)
	isSuccess, result := mixin.List("userCard", wherePhrase, "")
	if !isSuccess {
		return false, result.(string)
	}

	if len(*result.(*[]models.UserCard)) < 1 {
		return false, fmt.Sprintf("UserCard is not exists, UserCard: %d", userNo)
	}

	return true, "success"
}

// CheckOrderNoValid : 주문 번호 유효성 체크
func CheckOrderNoValid(userNo int, orderNo int) (bool, string) {
	if orderNo == 0 {
		return false, "orderNo is empty"
	}

	wherePhrase := fmt.Sprintf("UserNo = %d AND OrderNo = %d", userNo, orderNo)
	isSuccess, result := mixin.List("orderMst", wherePhrase, "")
	if !isSuccess {
		return false, result.(error).Error()
	}

	orderMst := *result.(*[]models.OrderMst)
	if len(orderMst) < 1 {
		return false, "orderNo is wrong"
	}

	return true, "success"
}

// CheckSubsNoValid : 구독 번호 유효성 체크
func CheckSubsNoValid(userNo int, subsNo int) (bool, string) {
	if subsNo == 0 {
		return false, "subsNo is empty"
	}

	wherePhrase := fmt.Sprintf("UserNo = %d AND SubsNo = %d", userNo, subsNo)
	isSuccess, result := mixin.List("subsMst", wherePhrase, "")
	if !isSuccess {
		return false, result.(error).Error()
	}

	subsMst := *result.(*[]models.SubsMst)
	if len(subsMst) < 1 {
		return false, "subsMst is wrong"
	}

	return true, "success"
}
