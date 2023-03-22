package middleware

import (
	"errors"
	"github.com/gin-gonic/gin"
)

//对外提供业务服务的handler接口(多态特性)
type Handler interface {
	//常规业务函数，返回data
	Handler(service Handler) (any, error)
	//将context传递到全局
	ContextHandler(c *gin.Context, service Handler) (any, error)
	//需要重载的接口
	Exec() (any, error)
	ContextExec(c *gin.Context) (any, error)
}

type AbstractHandler struct {
}

//controller调用业务入口，使用多态的特性执行-若需要增强可直接重载该方法
func (s *AbstractHandler) Handler(service Handler) (any, error) {
	return service.Exec()
}
func (s *AbstractHandler) ContextHandler(c *gin.Context, service Handler) (any, error) {
	return service.ContextExec(c)
}

//真正业务类需要实现的方法
func (s *AbstractHandler) Exec() (any, error) {
	return nil, errors.New("not server")
}
func (s *AbstractHandler) ContextExec(c *gin.Context) (any, error) {
	return nil, errors.New("not server")
}
