package cors

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Cors() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 请求方法
		method := ctx.Request.Method
		// 请求头部
		origin := ctx.Request.Header.Get("Origin")
		if origin != "" {
			// 这是允许访问所有域
			ctx.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			// 服务器支持的所有跨域请求的方法,为了避免浏览次请求的多次'预检'请求
			ctx.Writer.Header().Set("Access-Control-Allow-Methods", "*")
			// header的类型
			ctx.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Access-Token, Access-Control-Request-Method, Access-Control-Request-Headers, Authorization, X-CSRF-Token, X-Token, Token")
			// 允许跨域设置, 可以返回其他子段
			// 跨域关键设置 让浏览器可以解析
			ctx.Writer.Header().Set("Access-Control-Expose-Headers", "*")
			// 缓存请求信息 单位为秒
			ctx.Writer.Header().Set("Access-Control-Max-Age", "172800")
			//  跨域请求是否需要带cookie信息 默认设置为true
			ctx.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			// 设置返回格式是json
			ctx.Writer.Header().Set("Content-Type", "application/json;charset=UTF-8")
		}

		// 放行所有OPTIONS方法
		if method == "OPTIONS" {
			ctx.AbortWithStatus(http.StatusNoContent)
			return
		}
		// 处理请求
		ctx.Next()
	}
}
