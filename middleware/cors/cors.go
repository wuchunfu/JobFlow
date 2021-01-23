package cors

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func Cors() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 请求方法
		method := ctx.Request.Method
		// 请求头部
		origin := ctx.Request.Header.Get("Origin")
		// 声明请求头keys
		var headerKeys []string
		for k, _ := range ctx.Request.Header {
			headerKeys = append(headerKeys, k)
		}
		headerStr := strings.Join(headerKeys, ", ")
		if headerStr != "" {
			headerStr = fmt.Sprintf("Access-Control-Allow-Origin, Access-Control-Allow-Headers, Access-Control-Request-Headers, Access-Control-Request-Method, %s", headerStr)
		} else {
			headerStr = "Access-Control-Allow-Origin, Access-Control-Allow-Headers, Access-Control-Request-Headers, Access-Control-Request-Method"
		}
		if origin != "" {
			// 这是允许访问所有域
			ctx.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			// 服务器支持的所有跨域请求的方法,为了避免浏览次请求的多次'预检'请求
			ctx.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
			// header的类型
			ctx.Writer.Header().Set("Access-Control-Allow-Headers", "Access-Control-Request-Headers, Access-Control-Request-Method, User-Agent, Authorization, Content-Length, X-CSRF-Token, X-Token, Token, session, X_Requested_With, Accept, Accept-Encoding, Origin Referer, Host, Connection, Accept-Language, DNT, X-CustomHeader, Keep-Alive, User-Agent, X-Requested-With, If-Modified-Since, Cache-Control, Content-Type, Pragma, Sec-Fetch-Mode, Sec-Fetch-Site, Sec-Fetch-Dest")
			// 允许跨域设置, 可以返回其他子段
			// 跨域关键设置 让浏览器可以解析
			ctx.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type, Expires, Last-Modified, Pragma, FooBar")
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
