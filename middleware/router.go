package middleware

import "github.com/gin-gonic/gin"

/**
controller路由表
*/
var routerTables = make([]func(g *gin.RouterGroup), 0)

//用于controller注册路由
func RegisterRouter(f func(g *gin.RouterGroup)) {
	routerTables = append(routerTables, f)
}
