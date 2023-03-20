package middleware

import "github.com/gin-gonic/gin"

//对外提供初始化路由的方法
func InitHttp(prefix string) *gin.Engine {
	r := gin.New()
	//默认日志
	r.Use(gin.Logger())
	//500
	r.Use(gin.Recovery())

	//自定义拦截器
	for _, interceptor := range interceptorTables {
		r.Use(interceptor)
	}

	//自定义路由表
	kap := r.Group(prefix)
	for _, router := range routerTables {
		router(kap)
	}
	return r
}
