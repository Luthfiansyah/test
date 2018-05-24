package main

// import library
import (
	"io"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/julienschmidt/httprouter"

	"fmt"
	"net/http"
	"strconv"
)

var db *gorm.DB

func init() {
	//open a db connection
	var err error
	db, err = gorm.Open("mysql", "root:root@/demo?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&todo{})
}

func main() {

	// Logging to a file.
	f, _ := os.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(f)

	router := gin.Default()

	v1 := router.Group("/api/v1/")
	{
		v1.GET("/", func(c *gin.Context) {
			c.JSON(200, "Welcome!")
		})
		v1.POST("todos/", createTodo)
		v1.GET("todos/", fetchAllTodo)
		v1.GET("todos/:id", fetchSingleTodo)
		v1.PUT("todos/:id", updateTodo)
		v1.DELETE("todos/:id", deleteTodo)
	}

	// router := httprouter.New()
	// router.GET("/", index)

	// log.Fatal(http.ListenAndServe(":8080", router))
	router.Run("localhost:8080")

}

type (
	// todo describes a todo type
	todo struct {
		gorm.Model
		Title     string `json:"title"`
		Completed int    `json:"completed"`
	}

	// transformedTodo represents a formatted todo
	transformedTodo struct {
		ID        uint   `json:"id"`
		Title     string `json:"title"`
		Completed bool   `json:"completed"`
	}
)

func index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome!\n")
}

// createTodo add a new todo
func createTodo(c *gin.Context) {
	completed, _ := strconv.Atoi(c.PostForm("completed"))
	todo := todo{Title: c.PostForm("title"), Completed: completed}
	db.Save(&todo)
	c.JSON(http.StatusCreated, gin.H{"status": http.StatusCreated, "message": "Todo item created successfully!", "resourceId": todo.ID})
}

// fetchAllTodo fetch all todos
func fetchAllTodo(c *gin.Context) {
	var todos []todo
	var _todos []transformedTodo

	db.Find(&todos)

	if len(todos) <= 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No todo found!"})
		return
	}

	//transforms the todos for building a good response
	for _, item := range todos {
		completed := false
		if item.Completed == 1 {
			completed = true
		} else {
			completed = false
		}
		_todos = append(_todos, transformedTodo{ID: item.ID, Title: item.Title, Completed: completed})
	}
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": _todos})
}

// fetchSingleTodo fetch a single todo
func fetchSingleTodo(c *gin.Context) {
	var todo todo
	todoID := c.Param("id")

	db.First(&todo, todoID)

	if todo.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No todo found!"})
		return
	}

	completed := false
	if todo.Completed == 1 {
		completed = true
	} else {
		completed = false
	}

	_todo := transformedTodo{ID: todo.ID, Title: todo.Title, Completed: completed}
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": _todo})
}

// updateTodo update a todo
func updateTodo(c *gin.Context) {
	var todo todo
	todoID := c.Param("id")

	db.First(&todo, todoID)

	if todo.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No todo found!"})
		return
	}

	db.Model(&todo).Update("title", c.PostForm("title"))
	completed, _ := strconv.Atoi(c.PostForm("completed"))
	db.Model(&todo).Update("completed", completed)
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Todo updated successfully!"})
}

// deleteTodo remove a todo
func deleteTodo(c *gin.Context) {
	var todo todo
	todoID := c.Param("id")

	db.First(&todo, todoID)

	if todo.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No todo found!"})
		return
	}

	db.Delete(&todo)
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Todo deleted successfully!"})
}
