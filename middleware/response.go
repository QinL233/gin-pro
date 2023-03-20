package middleware

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"

	"net/http"
	"reflect"
	"strings"
)

/**
gin context全局变量用于存储传参，返回，错误日志
*/
var REQUEST_PARAM = "request_param"
var RESPONSE_BODY = "response_body"
var RESPONSE_ERR = "response_err"

/**
标准类型的返回封装
*/
type Response[T any] struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data T      `json:"data"`
}

type PageResponse[T any] struct {
	Response[T]
	//当前页
	CurrentPage int `json:"currentPage"`
	//每页数量
	PageSize int `json:"pageSize"`
	//总页数
	PageCount int `json:"pageCount"`
	//总数量
	TotalCount int `json:"totalCount"`
}

type NodeResult[T any] struct {
	Node      T               `json:"node"`
	ChildList []NodeResult[T] `json:"childList"`
}

/**
封装返回
*/
//通用成功返回
func Success(c *gin.Context, data ...any) {
	rs := fmt.Sprintf("%v", data)
	if len(rs) > 255 {
		c.Set(RESPONSE_BODY, rs[:255])
	} else {
		c.Set(RESPONSE_BODY, rs)
	}

	if len(data) > 0 && data[0] != nil {
		//如果data类型是Response则不重复封装
		if strings.HasPrefix(reflect.TypeOf(data[0]).String(), "middleware.PageResponse") {
			c.JSON(http.StatusOK, data[0])
		} else {
			c.JSON(http.StatusOK, Response[any]{
				Code: http.StatusOK,
				Msg:  "success",
				Data: data[0],
			})
		}
	} else {
		c.JSON(http.StatusOK, Response[any]{
			Code: http.StatusOK,
			Msg:  "success",
			Data: nil,
		})
	}
	return
}

//通用错误
func Error(c *gin.Context, errMsg error, errCode ...int) {
	fmt.Println(fmt.Errorf("[err] %v", errMsg))
	if len(errMsg.Error()) > 255 {
		c.Set(RESPONSE_ERR, errMsg.Error()[:255])
	} else {
		c.Set(RESPONSE_ERR, errMsg.Error())
	}
	if len(errCode) > 0 {
		c.JSON(http.StatusOK, Response[any]{
			Code: errCode[0],
			Msg:  fmt.Sprintf("%v", errMsg),
		})
	} else {
		c.JSON(http.StatusOK, Response[any]{
			Code: 500,
			Msg:  fmt.Sprintf("%v", errMsg),
		})
	}
	c.Abort()
	return
}

//参数错误
func ParamError(c *gin.Context) {
	Error(c, errors.New("参数错误"))
	return
}
