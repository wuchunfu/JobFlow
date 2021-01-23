package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/wuchunfu/JobFlow/middleware/database"
	"github.com/wuchunfu/JobFlow/middleware/jwt"
	"github.com/wuchunfu/JobFlow/model/userModel"
	"net/http"
	"strings"
)

func Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 获取 Auth-Token header
		tokenString := ctx.GetHeader("Auth-Token")
		// vcalidate token formate
		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer") {
			ctx.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "权限不足"})
			ctx.Abort()
			return
		}

		tokenString = tokenString[7:]
		token, claims, err := jwt.ParseToken(tokenString)
		if err != nil || !token.Valid {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code": http.StatusUnauthorized,
				"data": nil,
				"msg":  "权限不足",
			})
			ctx.Abort()
			return
		}

		// 验证通过后获取Claiim中的userId
		userId := claims.UserId
		db := database.GetDB()
		//var user *userModel.User
		user := new(userModel.User)
		db.First(&user, userId)

		// 用户不存在
		if user.UserId < 0 {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code": http.StatusUnauthorized,
				"data": nil,
				"msg":  "权限不足",
			})
			ctx.Abort()
			return
		}
		// 用户存在，将user信息写入上下文
		ctx.Set("user", user)
		ctx.Next()
	}
}
