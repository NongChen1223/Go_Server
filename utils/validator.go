package utils

import (
	"github.com/go-playground/validator/v10"
)

type Validator interface {
	GetMessages() ValidatorMessages
}

type ValidatorMessages map[string]string

/**
 *  GetErrorMsg
 *  @Description: 程序会判断请求结构体是否实现了 Validator 接口 如果 request 实现了 Validator 接口，则可以自定义错误信息
 *  @param request 参数是被验证的请求结构体
 *  @param err 是用 Go 自带的验证器库 validator 验证参数时返回的错误
 *  @return string
 */
func GetErrorMsg(request interface{}, err error) string {
	// 处理特殊错误
	if err != nil && err.Error() == "EOF" {
		return "请求数据为空，请提供必要的信息"
	}
	if _, isValidatorErrors := err.(validator.ValidationErrors); isValidatorErrors {
		_, isValidator := request.(Validator)
		for _, v := range err.(validator.ValidationErrors) {
			// 若 request 结构体实现 Validator 接口即可实现自定义错误信息
			if isValidator {
				if message, exist := request.(Validator).GetMessages()[v.Field()+"."+v.Tag()]; exist {
					return message
				}
			}
			return v.Error()
		}
	}
	return "参数格式错误，请检查输入"
}
