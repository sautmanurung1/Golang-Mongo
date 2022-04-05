package routes

import (
	"Project-Rest-Api/controller"

	"github.com/labstack/echo/v4"
)

func UserRoute(e *echo.Echo) {
	// All Route User In Here
	e.POST("/user", controller.CreateUser)
    e.GET("/user/:userId", controller.GetAUser)
    e.PUT("/user/:userId", controller.EditAUser)
    e.DELETE("/user/:userId", controller.DeleteAUser)
    e.GET("/users", controller.GetAllUsers)
}