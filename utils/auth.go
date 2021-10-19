package utils

import "golang.org/x/crypto/bcrypt"

// HashAndSalt : 입력 받은 password를 암호화
func HashAndSalt(password string) (bool, string) {
	cost := 10
	hash, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return false, err.Error()
	}

	return true, string(hash)
}

// ComparePasswords : 암호호된 패스워드와 평문 패스워드 비교
func ComparePasswords(hashedPassword string, plainPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	if err != nil {
		return false, err.Error()
	}

	return true, "success"
}
