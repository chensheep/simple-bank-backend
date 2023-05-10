package api

import (
	"github.com/chensheep/simple-bank-backend/util"
	"github.com/go-playground/validator/v10"
)

var currencyValidator validator.Func = func(fl validator.FieldLevel) bool {

	v, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}

	return util.IsSupportedCurrency(v)
}
