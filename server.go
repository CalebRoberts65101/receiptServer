package main

import (
	"github.com/labstack/echo/v4"
)

const PORT_NUMBER = "8080"

var recordStore = make(map[string]processReceiptRecord)

var server *echo.Echo

func startServer() {
	server = echo.New()
	server.Validator = getValidator()
	server.POST("/receipts/process", processReceipt)
	server.GET("/receipts/:id/points", getPoints)
	err := server.Start(":" + PORT_NUMBER)
	if err != nil {
		server.Logger.Fatal(err)
	}
}

func main() {
	startServer()
}
