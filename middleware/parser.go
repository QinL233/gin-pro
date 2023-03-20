package middleware

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"github.com/unknwon/com"
	"reflect"
)

/**
参数解析器
1、分组各类传参
2、传参校验
3、根据传参转发到指定handler
4、封装返回
*/
type Parser interface {
	//解析方式
	Body(params ...string)
	Query(params ...string)
	Path(params ...string)
	Form(params ...string)
}

type AbstractParser[T Handler] struct {
	C         *gin.Context
	IsContext bool
	//用于反射的参数/server
	param T
}

//自定义参数以及提取器
var paramTables = make(map[string]func(g *gin.Context) (string, error))

func RegisterParam(paramName string, f func(g *gin.Context) (string, error)) {
	paramTables[paramName] = f
}

//通用并从body中获取参数传递
func (b *AbstractParser[T]) Body(params ...string) {
	//1、从body读取序列化数据
	body, err := b.C.GetRawData()
	if len(body) < 1 {
		ParamError(b.C)
		return
	}
	if err != nil {
		Error(b.C, err)
		return
	}
	//2、构建对象：将二进制反序列化为对象
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	if err = json.Unmarshal(body, &b.param); err != nil {
		ParamError(b.C)
		return
	}
	//3、尝试从body以外的其它途径获取 param 的值
	if len(params) > 0 {
		v := reflect.ValueOf(b.param).Elem() //通过类型创建对象
		for _, paramName := range params {
			fv := v.FieldByName(capitalize(paramName))
			if fv.CanSet() && fv.IsNil() {
				if f, ok := paramTables[paramName]; ok {
					if cv, err := f(b.C); err != nil {
						Error(b.C, err)
						return
					} else {
						fv.SetString(cv)
					}
				}
			}
		}
	}
	b.serverHandler()
	return
}

//通用从 Query 获取参数传递:url?id=1
func (b *AbstractParser[T]) Query(params ...string) {
	//1、通过参数名从param提取参数构建map
	t := reflect.TypeOf(b.param).Elem() //反射获取类型
	v := reflect.New(t).Elem()          //通过类型创建对象
	for _, paramName := range params {
		fv := v.FieldByName(capitalize(paramName))
		if fv.CanSet() {
			if value, exist := b.C.GetQuery(paramName); exist {
				//根据参数类型从param中获取
				switch fv.Kind() {
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					fv.SetInt(com.StrTo(value).MustInt64())
				case reflect.Float64, reflect.Float32:
					fv.SetFloat(com.StrTo(value).MustFloat64())
				case reflect.String:
					fv.SetString(value)
				default:
					ParamError(b.C)
					return
				}
			} else if f, ok := paramTables[paramName]; ok {
				if cv, err := f(b.C); err != nil {
					Error(b.C, err)
					return
				} else {
					fv.SetString(cv)
				}
			}
		}
	}
	b.serverHandlerWithReflect(v)
	return
}

//通用从 param 获取参数传递:url/{id}/{params...}
func (b *AbstractParser[T]) Path(params ...string) {
	//1、通过参数名从param提取参数构建map
	t := reflect.TypeOf(b.param).Elem() //反射获取类型
	v := reflect.New(t).Elem()          //通过类型创建对象
	for _, paramName := range params {
		fv := v.FieldByName(capitalize(paramName))
		if fv.CanSet() {
			value := b.C.Param(paramName)
			if value != "" {
				//根据参数类型从param中获取
				switch fv.Kind() {
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					fv.SetInt(com.StrTo(value).MustInt64())
				case reflect.Float64, reflect.Float32:
					fv.SetFloat(com.StrTo(value).MustFloat64())
				case reflect.String:
					fv.SetString(value)
				default:
					ParamError(b.C)
					return
				}
			} else if f, ok := paramTables[paramName]; ok {
				if cv, err := f(b.C); err != nil {
					Error(b.C, err)
					return
				} else {
					fv.SetString(cv)
				}
			}
		}
	}
	b.serverHandlerWithReflect(v)
	return
}

//通用从 form 中获取文件或参数并反射到serviceParam中并执行
func (b *AbstractParser[T]) Form(params ...string) {
	//1、通过参数名从param提取参数构建map
	t := reflect.TypeOf(b.param).Elem() //反射获取类型
	v := reflect.New(t).Elem()          //通过类型创建对象
	form, err := b.C.MultipartForm()
	if err != nil {
		ParamError(b.C)
		return
	}
	for _, paramName := range params {
		fv := v.FieldByName(capitalize(paramName))
		if fv.CanSet() {
			if file, ok := form.File[paramName]; ok {
				if len(file) > 0 {
					fv.Set(reflect.ValueOf(file))
				}
			} else if value, ok := form.Value[paramName]; ok {
				if len(value) > 0 {
					//根据参数类型从param中获取
					switch fv.Kind() {
					case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
						fv.SetInt(com.StrTo(value[0]).MustInt64())
					case reflect.Float64, reflect.Float32:
						fv.SetFloat(com.StrTo(value[0]).MustFloat64())
					case reflect.String:
						fv.SetString(value[0])
					default:
						ParamError(b.C)
						return
					}
				}
			} else if f, ok := paramTables[paramName]; ok {
				if cv, err := f(b.C); err != nil {
					Error(b.C, err)
					return
				} else {
					fv.SetString(cv)
				}
			}
		}
	}
	b.serverHandlerWithReflect(v)
	return
}

//将value映射成server并自行
func (b *AbstractParser[T]) serverHandlerWithReflect(v reflect.Value) {
	//1、将value转换为server
	m := v.Addr().Interface()
	if r, ok := m.(T); ok {
		b.param = r
		b.serverHandler()
	} else {
		Error(b.C, errors.New("传参与结构体不一致"))
		return
	}

}

//执行param校验以及绑定的server方法获取返回
func (b *AbstractParser[T]) serverHandler() {
	//1、全局存储param方便取证
	b.C.Set(REQUEST_PARAM, fmt.Sprintf("%v", b.param))
	//2、校验参数
	if err := Valid(b.param); err != nil {
		Error(b.C, err)
		return
	}
	//3、执行handler
	if b.IsContext {
		if _, err := b.param.ContextHandler(b.param, b.C); err != nil {
			Error(b.C, err)
		}
	} else {
		if r, err := b.param.Handler(b.param); err != nil {
			Error(b.C, err)
		} else {
			Success(b.C, r)
		}
	}
}

//首字母大写
func capitalize(str string) string {
	var upperStr string
	vv := []rune(str)
	for i := 0; i < len(vv); i++ {
		if i == 0 {
			if vv[i] >= 97 && vv[i] <= 122 { //判断是大写字母
				vv[i] -= 32 // string的码表相差32位
				upperStr += string(vv[i])
			} else {
				return str
			}
		} else {
			upperStr += string(vv[i])
		}
	}
	return upperStr
}
