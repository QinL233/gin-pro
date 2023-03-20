package test

import (
	"fmt"
	"github.com/QinL233/gin-pro/middleware"
	"github.com/gin-gonic/gin"
	"testing"
)

/**
重载handler
*/
type TestService struct {
	middleware.AbstractHandler
	str string
}

func (s *TestService) Handler(service middleware.Handler) (any, error) {
	s.str = "test server"
	return service.Exec()
}

//实现方法
type TestParam struct {
	TestService
	Token string `validate:"required"`
	Param string `validate:"required"`
}

func (p *TestParam) Exec() (any, error) {
	return fmt.Sprintf("[%s]handler [%s]token [%s]param", p.str, p.Token, p.Param), nil
}

type Test2Param struct {
	TestService
	Token string `validate:"required"`
	Param string `validate:"required"`
}

func (p *Test2Param) Exec() (any, error) {
	return fmt.Sprintf("[%s]handler [%s]token [%s]param", p.str, p.Token, p.Param), nil
}

type Test3Param struct {
	TestService
	Token string `validate:"required"`
	Param string `validate:"required"`
}

func (p *Test3Param) Exec() (any, error) {
	return (&Test2Param{Token: p.Token, Param: p.Param, TestService: TestService{str: "test3"}}).Exec()
}

//启动http
func TestWeb(t *testing.T) {
	//加载到路由
	middleware.RegisterRouter(func(g *gin.RouterGroup) {
		home := g.Group("/test").Use()
		{
			//get请求
			home.GET("/1", func(c *gin.Context) {
				//query参数解析并执行
				(&middleware.AbstractParser[*TestParam]{C: c}).Query("param", "token")
			})
			home.GET("/2/:param/:token", func(c *gin.Context) {
				//query参数解析并执行
				(&middleware.AbstractParser[*Test2Param]{C: c}).Path("param", "token")
			})
			home.GET("/3", func(c *gin.Context) {
				//query参数解析并执行
				(&middleware.AbstractParser[*Test3Param]{C: c}).Query("param", "token")
			})
		}
	})
	middleware.RegisterParam("token", func(g *gin.Context) (string, error) {
		return g.GetHeader("token"), nil
	})
	gin.SetMode("debug")
	middleware.InitHttp("/qz").Run(":8085")
}
