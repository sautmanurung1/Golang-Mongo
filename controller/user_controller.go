package controller

import (
	"Project-Rest-Api/config"
	"Project-Rest-Api/models"
	"Project-Rest-Api/response"
	"context"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = config.GetCollection(config.DB, "users")
var validate = validator.New()

func CreateUser(c echo.Context) error{
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var user models.User
	defer cancel()

	if err := c.Bind(&user); err != nil{
		return c.JSON(http.StatusBadRequest, response.UserResponse{
			Status : http.StatusBadRequest,
			Message : "error",
			Data : &echo.Map{"data" : err.Error()}})
	}

	if validationErr := validate.Struct(&user); validationErr != nil{
		return c.JSON(http.StatusBadRequest, response.UserResponse{
			Status: http.StatusBadRequest, 
			Message: "error", 
			Data: &echo.Map{
				"data": validationErr.Error()}})
	}

	newUser := models.User{
		ID: primitive.NewObjectID(),
		Name: user.Name,
		Location: user.Location,
		Tittle: user.Tittle,
	}

	result, err := userCollection.InsertOne(ctx, newUser)
	if err != nil{
		return c.JSON(http.StatusInternalServerError, response.UserResponse{
			Status: http.StatusInternalServerError,
			Message: "Error",
			Data: &echo.Map{
				"Data" : err.Error(),
			},
		})
	}
	return c.JSON(http.StatusCreated,response.UserResponse{
		Status: http.StatusCreated,
		Message: "Success Create User",
		Data: &echo.Map{"Data" : result},
	})
}

func GetAUser(c echo.Context) error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    userId := c.Param("userId")
    var user models.User
    defer cancel()
    objId, _ := primitive.ObjectIDFromHex(userId)
    err := userCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&user)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, response.UserResponse{
			Status: http.StatusInternalServerError, 
			Message: "error", 
			Data: &echo.Map{"data": err.Error()}})
    }
    return c.JSON(http.StatusOK, response.UserResponse{
		Status: http.StatusOK, 
		Message: "success", 
		Data: &echo.Map{"data": user}})
}

func EditAUser(c echo.Context) error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    userId := c.Param("userId")
    var user models.User
    defer cancel()
    objId, _ := primitive.ObjectIDFromHex(userId)
    //validate the request body
    if err := c.Bind(&user); err != nil {
        return c.JSON(http.StatusBadRequest, response.UserResponse{
			Status: http.StatusBadRequest, 
			Message: "error", 
			Data: &echo.Map{"data": err.Error()}})
    }
    //use the validator library to validate required fields
    if validationErr := validate.Struct(&user); validationErr != nil {
        return c.JSON(http.StatusBadRequest, response.UserResponse{
			Status: http.StatusBadRequest, 
			Message: "error", 
			Data: &echo.Map{"data": validationErr.Error()}})
    }
    update := bson.M{"name": user.Name, "location": user.Location, "title": user.Tittle}
    result, err := userCollection.UpdateOne(ctx, bson.M{"id": objId}, bson.M{"$set": update})
    if err != nil {
        return c.JSON(http.StatusInternalServerError, response.UserResponse{
			Status: http.StatusInternalServerError, 
			Message: "error", 
			Data: &echo.Map{"data": err.Error()}})
    }
    //get updated user details
    var updatedUser models.User
    if result.MatchedCount == 1 {
        err := userCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&updatedUser)
        if err != nil {
            return c.JSON(http.StatusInternalServerError, response.UserResponse{
				Status: http.StatusInternalServerError, 
				Message: "error", 
				Data: &echo.Map{"data": err.Error()}})
        }
    }
    return c.JSON(http.StatusOK, response.UserResponse{
		Status: http.StatusOK, 
		Message: "success", 
		Data: &echo.Map{"data": updatedUser}})
}

func DeleteAUser(c echo.Context) error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    userId := c.Param("userId")
    defer cancel()
    objId, _ := primitive.ObjectIDFromHex(userId)
    result, err := userCollection.DeleteOne(ctx, bson.M{"id": objId})
    if err != nil {
        return c.JSON(http.StatusInternalServerError, response.UserResponse{
			Status: http.StatusInternalServerError, 
			Message: "error", 
			Data: &echo.Map{"data": err.Error()}})
    }
    if result.DeletedCount < 1 {
        return c.JSON(http.StatusNotFound, response.UserResponse{
			Status: http.StatusNotFound, 
			Message: "error", 
			Data: &echo.Map{"data": "User with specified ID not found!"}})
    }
    return c.JSON(http.StatusOK, response.UserResponse{
		Status: http.StatusOK, 
		Message: "success", 
		Data: &echo.Map{"data": "User successfully deleted!"}})
}

func GetAllUsers(c echo.Context) error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    var users []models.User
    defer cancel()
    results, err := userCollection.Find(ctx, bson.M{})
    if err != nil {
        return c.JSON(http.StatusInternalServerError, response.UserResponse{
			Status: http.StatusInternalServerError, 
			Message: "error", 
			Data: &echo.Map{"data": err.Error()}})
    }

    //reading from the db in an optimal way
    defer results.Close(ctx)
    for results.Next(ctx) {
        var singleUser models.User
        if err = results.Decode(&singleUser); err != nil {
            return c.JSON(http.StatusInternalServerError, response.UserResponse{
				Status: http.StatusInternalServerError, 
				Message: "error", 
				Data: &echo.Map{"data": err.Error()}})
        }
        users = append(users, singleUser)
    }
    return c.JSON(http.StatusOK, response.UserResponse{
		Status: http.StatusOK, 
		Message: "success", 
		Data: &echo.Map{
			"data": users},
		})
}