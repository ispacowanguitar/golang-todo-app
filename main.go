package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type todo struct {
	gorm.Model
	Item      string `json:"title"`
	Completed bool   `json:"completed"`
	UserID    int64  `json:"user_id"`
}

func addTodo(context *gin.Context) {
	var newTodo todo

	if err := context.BindJSON(&newTodo); err != nil {
		return
	}
	newTodo.Completed = false
	fmt.Print(newTodo)
	todoRepo().Create(&newTodo)
	context.JSON(http.StatusCreated, newTodo)
}

func getTodos(context *gin.Context) {
	var records []todo
	todoRepo().Find(&records)
	context.JSON(http.StatusOK, records)
}

func getTodoById(id string) (*todo, error) {
	var todo todo
	result := todoRepo().First(&todo, id)
	if result.Error != nil {
		return nil, errors.New("not found")
	}
	return &todo, nil
}

func getTodo(context *gin.Context) {
	todo, err := getTodoById(context.Param("id"))

	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
	} else {
		context.JSON(http.StatusOK, &todo)
	}
}

func updateTodo(c *gin.Context) {
	id := c.Param("id")
	todo, err := getTodoById(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
	} else {
		c.BindJSON(&todo)
		todo.Completed = true
		db.Save(&todo)
		c.JSON(http.StatusOK, todo)
	}
}

var db *gorm.DB
var err error

func todoRepo() *gorm.DB {
	godotenv.Load()
	host := os.Getenv("HOST")
	dbPort := os.Getenv("DBPORT")
	userName := os.Getenv("USER")
	dbName := os.Getenv("NAME")
	password := os.Getenv("PASSWORD")

	dbURI := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s port=%s", host, userName, dbName, password, dbPort)
	fmt.Print("PW: ", dbURI)

	db, err = gorm.Open(postgres.Open(dbURI), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database")
	} else {
		fmt.Print("Successfully connected to database")
	}

	db.AutoMigrate(&todo{}, &user{})
	return db
}

func main() {
	// get env vars
	todoRepo()

	// Close connection when main func finishes
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	router := gin.Default()
	LoadUserRoutes(router)
	router.GET("/todos", getTodos)
	router.GET("/todos/:id", getTodo)
	router.PATCH("/todos/:id", updateTodo)
	router.POST("/todos", addTodo)
	router.Run("localhost:9090")
}
