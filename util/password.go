package util

import "golang.org/x/crypto/bcrypt"

//生成密码
func GetPWD(str string) string {
	s, _ := bcrypt.GenerateFromPassword([]byte(str), bcrypt.DefaultCost)
	return string(s)
}

//检查密码
func CmpPWD(pwd string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(pwd), []byte(password))
	if err != nil {
		return false
	}
	return true
}
