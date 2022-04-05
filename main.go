package main

import (
	"Project-Rest-Api/config"
	"Project-Rest-Api/routes"

	"github.com/labstack/echo/v4"
)

func main() {
    e := echo.New()

    config.ConnectDB()
	routes.UserRoute(e)
    e.Logger.Fatal(e.Start(":6000"))
}