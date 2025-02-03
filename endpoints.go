package main

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type processReceiptRequest struct {
	Retailer     string `json:"retailer" validate:"required,retailerValidator"`
	PurchaseDate string `json:"purchaseDate" validate:"required,purchaseDateValidator"`
	PurchaseTime string `json:"purchaseTime" validate:"required,purchaseTimeValidator"`
	Items        []struct {
		ShortDescription string `json:"shortDescription" validate:"required,shortDescriptionValidator"`
		Price            string `json:"price" validate:"required,moneyValidator"`
	} `json:"items" validate:"required,dive"`
	Total string `json:"total" validate:"required,moneyValidator"`
}

type processReceiptResult struct {
	ID string `json:"id"`
}

type processReceiptRecord struct {
	request processReceiptRequest
	result  processReceiptResult
	rules   []Rule
	points  int
}

func processReceipt(c echo.Context) error {
	var request processReceiptRequest
	err := c.Bind(&request)
	if err != nil {
		c.Echo().StdLogger.Println(err)
		return c.String(http.StatusBadRequest, "The receipt is invalid")
	}
	err = c.Validate(request)
	if err != nil {
		c.Echo().StdLogger.Println(err)
		return c.String(http.StatusBadRequest, "The receipt is invalid")
	}

	// Create record
	result := processReceiptResult{
		ID: uuid.New().String(),
	}
	rules := getAllRules()
	record := processReceiptRecord{
		request: request,
		result:  result,
		rules:   rules,
		points:  0,
	}

	// Process rules
	for _, rule := range record.rules {
		p, err := rule.score(&record.request)
		if err != nil {
			c.Echo().StdLogger.Println(err)
			return c.String(http.StatusBadRequest, "The receipt is invalid")
		}
		record.points += p
	}

	// Save result - In prod we would want to save to a DB but a in memory map is good enough.
	recordStore[record.result.ID] = record

	// Return ID
	return c.JSON(http.StatusOK, &record.result)
}

type getPointsRequest struct {
	ID string `param:"id" validation:"required,getPointsValidator"`
}

type getPointsResult struct {
	Points int `json:"points"`
}

func getPoints(c echo.Context) error {
	var request getPointsRequest
	err := c.Bind(&request)
	if err != nil {
		c.Echo().StdLogger.Println(err)
		return c.String(http.StatusNotFound, "No receipt found for that ID")
	}
	err = c.Validate(request)
	if err != nil {
		c.Echo().StdLogger.Println(err)
		return c.String(http.StatusNotFound, "No receipt found for that ID")
	}

	record, ok := recordStore[request.ID]
	if !ok {
		return c.String(http.StatusNotFound, "No receipt found for that ID")
	}

	result := getPointsResult{
		Points: record.points,
	}

	return c.JSON(http.StatusOK, result)
}
