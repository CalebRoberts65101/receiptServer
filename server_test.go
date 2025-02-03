package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
)

const processURL = "http://localhost:8080/receipts/process"
const pointsURL = "http://localhost:8080/receipts/%s/points"

func checkScore(t *testing.T, query []byte, score int) {
	var input = bytes.NewBuffer(query)
	res, err := http.Post(processURL, "application/json", input)
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	defer res.Body.Close()

	var r processReceiptResult
	err = json.NewDecoder(res.Body).Decode(&r)
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	res, err = http.Get(fmt.Sprintf(pointsURL, r.ID))
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	var p getPointsResult
	err = json.NewDecoder(res.Body).Decode(&p)
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	if p.Points != score {
		t.Fail()
	}

}

func TestExample1(t *testing.T) {
	var input = []byte(
		`
{
  "retailer": "Target",
  "purchaseDate": "2022-01-01",
  "purchaseTime": "13:01",
  "items": [
    {
      "shortDescription": "Mountain Dew 12PK",
      "price": "6.49"
    },{
      "shortDescription": "Emils Cheese Pizza",
      "price": "12.25"
    },{
      "shortDescription": "Knorr Creamy Chicken",
      "price": "1.26"
    },{
      "shortDescription": "Doritos Nacho Cheese",
      "price": "3.35"
    },{
      "shortDescription": "   Klarbrunn 12-PK 12 FL OZ  ",
      "price": "12.00"
    }
  ],
  "total": "35.35"
}
		`,
	)
	var points = 28
	checkScore(t, input, points)

}

func TestExample2(t *testing.T) {
	var input = []byte(
		`
{
  "retailer": "M&M Corner Market",
  "purchaseDate": "2022-03-20",
  "purchaseTime": "14:33",
  "items": [
    {
      "shortDescription": "Gatorade",
      "price": "2.25"
    },{
      "shortDescription": "Gatorade",
      "price": "2.25"
    },{
      "shortDescription": "Gatorade",
      "price": "2.25"
    },{
      "shortDescription": "Gatorade",
      "price": "2.25"
    }
  ],
  "total": "9.00"
}
		`,
	)
	var points = 109
	checkScore(t, input, points)

}

func TestExample3(t *testing.T) {
	var input = []byte(
		`
{
    "retailer": "Walgreens",
    "purchaseDate": "2022-01-02",
    "purchaseTime": "08:13",
    "total": "2.65",
    "items": [
        {"shortDescription": "Pepsi - 12-oz", "price": "1.25"},
        {"shortDescription": "Dasani", "price": "1.40"}
    ]
}
		`,
	)
	var points = 15
	checkScore(t, input, points)

}

func TestExample4(t *testing.T) {
	var input = []byte(
		`
{
    "retailer": "Target",
    "purchaseDate": "2022-01-02",
    "purchaseTime": "13:13",
    "total": "1.25",
    "items": [
        {"shortDescription": "Pepsi - 12-oz", "price": "1.25"}
    ]
}
		`,
	)
	var points = 31
	checkScore(t, input, points)

}
