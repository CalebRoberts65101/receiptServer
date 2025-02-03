package main

import (
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

func getValidator() *CustomValidator {
	v := validator.New()
	v.RegisterValidation("retailerValidator", retailerValidator)
	v.RegisterValidation("purchaseDateValidator", purchaseDateValidator)
	v.RegisterValidation("purchaseTimeValidator", purchaseTimeValidator)
	v.RegisterValidation("moneyValidator", moneyValidator)
	v.RegisterValidation("shortDescriptionValidator", shortDescriptionValidator)
	v.RegisterValidation("getPointsValidator", getPointsValidator)
	return &CustomValidator{validator: v}
}

func retailerValidator(fl validator.FieldLevel) bool {
	pattern := "^[\\w\\s\\-&]+$"
	matched, err := regexp.MatchString(pattern, fl.Field().String())
	if err != nil {
		fmt.Println(err)
		return false
	} else {
		return matched
	}
}

func purchaseDateValidator(fl validator.FieldLevel) bool {
	layout := "2006-01-02"
	_, err := time.Parse(layout, fl.Field().String())
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

func purchaseTimeValidator(fl validator.FieldLevel) bool {
	layout := "15:04"
	_, err := time.Parse(layout, fl.Field().String())
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

func moneyValidator(fl validator.FieldLevel) bool {
	pattern := "^\\d+\\.\\d{2}$"
	matched, err := regexp.MatchString(pattern, fl.Field().String())
	if err != nil {
		fmt.Println(err)
		return false
	} else {
		return matched
	}
}

func shortDescriptionValidator(fl validator.FieldLevel) bool {
	pattern := "^[\\w\\s\\-]+$"
	matched, err := regexp.MatchString(pattern, fl.Field().String())
	if err != nil {
		fmt.Println(err)
		return false
	} else {
		return matched
	}
}

func getPointsValidator(fl validator.FieldLevel) bool {
	pattern := "^\\S+$"
	matched, err := regexp.MatchString(pattern, fl.Field().String())
	if err != nil {
		fmt.Println(err)
		return false
	} else {
		return matched
	}
}
