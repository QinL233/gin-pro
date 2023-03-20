package middleware

import "github.com/gin-gonic/gin"

/**
controller 拦截器表
*/
var interceptorTables = make([]func(g *gin.Context), 0)

//用于controller注册拦截器
func RegisterInterceptor(f func(g *gin.Context)) {
	interceptorTables = append(interceptorTables, f)
}
