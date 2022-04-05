package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type user struct {
	gorm.Model
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Todos     []todo `json:"todos"`
}

func addUser(c *gin.Context) {
	var newUser user
	if err := c.BindJSON(&newUser); err != nil {
		c.JSON(http.StatusCreated, err.Error())
	}
	db.Create(&newUser)
	c.JSON(http.StatusCreated, newUser)
}

func getUsers(c *gin.Context) {
	var users []user
	db.Preload("Todos").Find(&users)
	c.JSON(http.StatusOK, users)
}

func getUserById(c *gin.Context) {
	var user user
	db.Preload("Todos").First(&user, c.Param("id"))
	c.JSON(http.StatusFound, user)
}

func LoadUserRoutes(router *gin.Engine) {
	router.GET("/users", getUsers)
	router.GET("/users/:id", getUserById)
	router.POST("/users", addUser)
}
