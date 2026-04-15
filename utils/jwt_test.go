package utils

import "testing"

// 实现解析Token测试用例
func TestParseJWT(t *testing.T) {
	token, err := GenerateJWT(1)
	if err != nil {
	}
	claims, err := ParseJWT(token)
	if err != nil {
		panic(err)
	}
	t.Log(claims.UserID)
	t.Log(claims.Issuer)
	t.Log(claims.IssuedAt)
}
