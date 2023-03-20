package api

import "github.com/go-playground/validator/v10"

var validateBookingType validator.Func = func(fl validator.FieldLevel) bool {
	if Type, ok := fl.Field().Interface().(string); ok {
		switch Type {
		case MaternityNursing, GiftPackage:
			return true
		default:
			return false

		}
	}

	return false
}
