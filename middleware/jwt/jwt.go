package jwt

import (
	"github.com/golang-jwt/jwt"
	logger "github.com/sirupsen/logrus"
	"strconv"
	"time"
)

var jwtKey = "a_secret_create"

type Claims struct {
	UserId int64
	jwt.StandardClaims
}

// 生成jwt
func ReleaseToken(userId int64, username string) map[string]string {
	// 过期时间为七天
	expireTime := time.Now().Add(7 * 24 * time.Hour).Unix()
	claims := make(jwt.MapClaims)
	// 过期时间
	claims["exp"] = expireTime
	claims["uid"] = userId
	// 签发时间
	claims["iat"] = time.Now().Unix()
	claims["issuer"] = "go-task"
	claims["username"] = username
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtKey))
	if err != nil {
		logger.Error("生成 Token 失败 " + err.Error())
		return nil
	}
	tokenMap := make(map[string]string, 10)
	tokenMap["userId"] = strconv.FormatInt(userId, 10)
	tokenMap["token"] = tokenString
	tokenMap["expireTime"] = time.Unix(expireTime, 0).Format("2006-01-02 15:04:05")
	return tokenMap
}

// 还原jwt
func ParseToken(tokenString string) (*jwt.Token, *Claims, error) {
	Claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, Claims, func(token *jwt.Token) (i interface{}, err error) {
		return []byte(jwtKey), nil
	})
	return token, Claims, err
}
