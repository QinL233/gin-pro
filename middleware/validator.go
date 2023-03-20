package middleware

import (
	"errors"
	"fmt"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
)

//通用校验器
func Valid(data interface{}) error {
	validate := validator.New()
	err := validate.Struct(data)
	if err != nil {
		//翻译
		trans := validateTransInit(validate)
		var result string
		for _, value := range err.(validator.ValidationErrors).Translate(trans) {
			result += value
		}
		return errors.New(result)
	}
	return nil
}

// 数据验证翻译器
func validateTransInit(validate *validator.Validate) ut.Translator {
	// 万能翻译器，保存所有的语言环境和翻译数据
	uni := ut.New(zh.New())
	// 翻译器
	trans, _ := uni.GetTranslator("zh")
	//验证器注册翻译器
	err := zhTranslations.RegisterDefaultTranslations(validate, trans)
	if err != nil {
		fmt.Println(err)
	}
	return trans
}
