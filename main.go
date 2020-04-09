package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"net/http"
	"strconv"
)

type (
	todoModel struct {
		gorm.Model
		Title     string `json:"title"`
		Completed int    `json:"completed"`
	}
	fmtTodo struct {
		ID        uint   `json:"id"`
		Title     string `json:"title"`
		Completed bool   `json:"completed"`
	}
)

func (todoModel) TableName() string {
	return "todos"
}

var db *gorm.DB

func init() {
	var err error
	var constr = "root:123456@(localhost)/test?charset=utf8&parseTime=True&loc=Local"
	db, err = gorm.Open("mysql", constr)
	if err != nil {
		panic("数据库连接失败")
	}
	db.AutoMigrate(&todoModel{})
}

const (
	JSON_SUCCESS int = 1
	JSON_ERROR   int = 0
)

func main() {
	r := gin.Default()
	v1 := r.Group("/api/v1/todo")
	{
		v1.GET("/query/list", queryAll)
		v1.POST("/add", addOne)
		v1.GET("/query/one/:id", queryOne)
		v1.POST("/update/:id", updateOne)
		v1.POST("/delete/:id", deleteOne)
	}
	r.Run()
}

//删除
func deleteOne(c *gin.Context) {
	id := c.Param("id")
	var todo todoModel
	db.First(&todo, id)
	if todo.ID == 0 {
		c.JSON(http.StatusOK, gin.H{"status": JSON_ERROR, "message": "id isn't exist"})
		return
	}
	db.Delete(&todo)
	c.JSON(http.StatusOK, gin.H{
		"status":  JSON_SUCCESS,
		"message": "ok",
	})
}

//更新一条条目
func updateOne(c *gin.Context) {
	id := c.Param("id")
	var todo todoModel
	db.First(&todo, id)
	if todo.ID == 0 {
		c.JSON(http.StatusOK, gin.H{"status": JSON_ERROR, "message": "id isn't exist"})
		return
	}
	db.Model(&todo).Update("title", c.PostForm("title"))
	completed, _ := strconv.Atoi(c.PostForm("completed"))
	db.Model(&todo).Update("completed", completed)
	c.JSON(http.StatusOK, gin.H{
		"status":  JSON_SUCCESS,
		"message": "ok",
	})
}

//获取单条条目
func queryOne(c *gin.Context) {
	id := c.Param("id")
	var todo todoModel
	db.First(&todo, id)
	if todo.ID == 0 {
		c.JSON(http.StatusOK, gin.H{"status": JSON_ERROR, "message": "id isn't exist"})
		return
	}
	completed := false
	if todo.Completed == 1 {
		completed = true
	}
	_todo := fmtTodo{
		ID:        todo.ID,
		Title:     todo.Title,
		Completed: completed,
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  JSON_SUCCESS,
		"message": "ok",
		"data":    _todo,
	})
}

/**
查询所有条目
*/
func queryAll(context *gin.Context) {
	var todos []todoModel
	var _todos []fmtTodo
	db.Find(&todos)
	if len(todos) <= 0 {
		context.JSON(http.StatusOK, gin.H{
			"status":  JSON_SUCCESS,
			"message": "todo list is empty",
		})
		return
	}
	//格式化
	for _, item := range todos {
		completed := false
		if item.Completed == 1 {
			completed = true
		}
		_todos = append(_todos, fmtTodo{
			ID:        item.ID,
			Title:     item.Title,
			Completed: completed,
		})
	}
	context.JSON(http.StatusOK, gin.H{
		"status":  JSON_SUCCESS,
		"message": "ok",
		"data":    _todos,
	})
}

/**
保存todo条目
*/
func addOne(context *gin.Context) {
	completed, err := strconv.Atoi(context.PostForm("completed"))
	if err != nil {
		context.JSON(http.StatusOK, gin.H{"status": JSON_ERROR, "message": "invalid param"})
		return
	}
	title := context.PostForm("title")
	todo := todoModel{Title: title, Completed: completed}
	db.Save(&todo)
	context.JSON(http.StatusOK,gin.H{
		"status":JSON_SUCCESS,
		"message":"ok",
		"data":todo,
	})
}
