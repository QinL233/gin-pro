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

//对外提供初始化路由的方法
func InitHttp(prefix string) *gin.Engine {
	r := gin.New()
	//默认日志
	r.Use(gin.Logger())
	//500
	r.Use(gin.Recovery())

	//路由表
	kap := r.Group(prefix)
	for _, router := range routerTables {
		router(kap)
	}
	return r
}
