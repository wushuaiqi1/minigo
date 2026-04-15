package utils

import "time"

type MyCustomClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func GenerateJWT(username string) (string, error) {
	// 创建一个我们自己的声明
	myClaims := MyCustomClaims{
		username,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
			Issuer:    "test",
		},
	}
	// 使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, myClaims)
}

func ParseJWT(tokenString string) (*MyCustomClaims, error) {
	// 解析token
	token, err := jwt.ParseWithClaims(tokenString, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*MyCustomClaims); ok && token.Valid {
		return claims, nil
	}
}
