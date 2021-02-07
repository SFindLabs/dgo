package utils

import (
	cn "github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	translations "github.com/go-playground/validator/v10/translations/zh"
	"github.com/pkg/errors"
	"reflect"
)

var ValidateTrans *validator.Validate
var trans ut.Translator

func init() {
	zh := cn.New()
	uni := ut.New(zh, zh)
	trans, _ = uni.GetTranslator("zh")

	ValidateTrans = validator.New()
	ValidateTrans.RegisterTagNameFunc(func(field reflect.StructField) string {
		label := field.Tag.Get("label")
		if label == "" {
			return field.Name
		}
		return label
	})
	_ = translations.RegisterDefaultTranslations(ValidateTrans, trans)
}

/*type loginTransBind struct {
	Name    string `form:"name"  validate:"required" label:"账号"`
	Passwd  string `form:"passwd"  validate:"required" label:"密码"`
	Captcha string `form:"captcha"  validate:"required" label:"验证码"`
}*/

/*var params loginTransBind
callbackName := kbase.GetParam(c, "callback")
if err := c.Bind(&params); err != nil {
	kbase.SendErrorJsonStr(c, kcode.OPERATION_WRONG, callbackName)
	return
}
if err := kutils.ValidateTranslate(params); err != nil {
	kbase.SendErrorParamsJsonStr(c, kcode.OPERATION_WRONG, err, callbackName)
	return
}*/

func ValidateTranslate(params interface{}) error {
	if err := ValidateTrans.Struct(params); err != nil {
		for _, errVal := range err.(validator.ValidationErrors) {
			return errors.New(errVal.Translate(trans))
		}
	}
	return nil
}
