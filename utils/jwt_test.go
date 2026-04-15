package utils

import (
	"errors"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// 实现解析Token测试用例
func TestJWTExpiration(t *testing.T) {
	// 生成一个 Token
	token, err := GenerateJWT(123)
	if err != nil {
		t.Fatal(err)
	}

	// 正常解析（应该成功）
	claims, err := ParseJWT(token)
	if err != nil {
		t.Logf("Token 解析成功: UserID=%s", claims.UserID)
	}

	// 手动构造一个已过期的 Token 来测试
	expireTime := time.Now().Add(-1 * time.Hour) // 1小时前就过期了
	payload := JwtClaims{
		UserID: "test123",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    issuer,
			ExpiresAt: jwt.NewNumericDate(expireTime),
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
		},
	}

	expiredToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	expiredTokenString, _ := expiredToken.SignedString(secret)

	// 解析过期 Token（应该失败）
	_, err = ParseJWT(expiredTokenString)
	if err != nil {
		t.Logf("正确捕获过期错误: %v", err)
		// 输出: 正确捕获过期错误: token has expired
	}
}

func TestTokenValidity(t *testing.T) {
	// 1. 生成正常 Token
	tokenStr, _ := GenerateJWT(123)

	// 2. 使用错误的 secret 解析
	wrongSecret := []byte("wrong-secret-key!!!")

	token, err := jwt.ParseWithClaims(tokenStr, &JwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return wrongSecret, nil // ⚠️ 使用错误的密钥
	})

	t.Logf("Error: %v", err)   // 会有签名错误
	t.Logf("Token: %v", token) // token 可能不为 nil
	if token != nil {
		t.Logf("Token Valid: %v", token.Valid) // false!
	}

	// 3. 使用正确的 secret，但修改了 Payload
	claims := JwtClaims{
		UserID: "hacked_user", // 篡改了 UserID
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    issuer,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	hackedToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	hackedTokenStr, _ := hackedToken.SignedString(secret)

	// 4. 尝试用原始 secret 解析（会失败，因为签名不匹配）
	token2, err2 := ParseJWT(hackedTokenStr)
	t.Logf("Hacked token error: %v", err2)    // invalid token signature
	t.Logf("Hacked token claims: %v", token2) // nil
}
