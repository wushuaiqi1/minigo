package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	issuer = "mini-go"
	expire = time.Hour * 24 * 7
	secret = []byte("qz001qz002qz003")
)

// JwtClaims Payload
type JwtClaims struct {
	UserID string `json:"user_id"` // 非数据库自增ID 防止用户通过遍历ID推测业务数据量或者恶意爬取
	jwt.RegisteredClaims
}

func GenerateJWT(uid int64) (string, error) {
	encodeUid, err := EncodeUId(uid)
	if err != nil {
		return "", err
	}
	expireTime := time.Now().Add(expire)

	payload := JwtClaims{
		encodeUid,
		jwt.RegisteredClaims{
			Issuer:    issuer,
			ExpiresAt: jwt.NewNumericDate(expireTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}
	// 使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	// token转换为字符串的形式
	return token.SignedString(secret)
}

func ParseJWT(tokenString string) (*JwtClaims, error) {
	// 解析token
	token, err := jwt.ParseWithClaims(tokenString, &JwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil {
		return nil, err
	}
	// 判断Payload是否过期
	claims, ok := token.Claims.(*JwtClaims)
	if !ok {
		return nil, jwt.ErrSignatureInvalid
	}
	if claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, jwt.ErrTokenExpired
	}
	if token.Valid {
		return claims, nil
	}
	return nil, jwt.ErrSignatureInvalid
}
