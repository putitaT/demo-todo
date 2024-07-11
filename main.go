package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/putitaT/todo-backend/database"
)

var db = database.ConnectDB()

type TodoDB struct {
	ID     int            `json:"id"`
	Title  sql.NullString `json:"title"`
	Status sql.NullString `json:"status"`
}

type Todo struct {
	ID     int
	Title  string
	Status string
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	r := gin.Default()

	// Configure CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// database.CreateTable()

	r.GET("/api/v1/todos", getTodosHandler)
	r.GET("/api/v1/todos/:id", getTodoByIdHandler)
	r.POST("/api/v1/todos", addTodoHandler)
	r.PUT("/api/v1/todos/:id", editTodoHandler)
	r.DELETE("/api/v1/todos/:id", deleteTodoHandler)
	r.PATCH("/api/v1/todos/:id/actions/status", updateStatusHandler)
	r.PATCH("/api/v1/todos/:id/actions/title", updateTitleHandler)

	srv := http.Server{
		Addr:    ":" + os.Getenv("PORT"),
		Handler: r,
	}

	closedChan := make(chan struct{})

	go func() {
		<-ctx.Done()
		fmt.Println("shutting down....")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				log.Println(err)
			}
		}

		close(closedChan)
	}()

	if err := srv.ListenAndServe(); err != nil {
		log.Println(err)
	}

	<-closedChan
	fmt.Println("bye")
}

func getTodosHandler(ctx *gin.Context) {
	rows, err := db.Query("SELECT id, title, status FROM todos ORDER BY id")
	if err != nil {
		log.Fatal("can't query all todos", err)
	}
	todos := []Todo{}
	for rows.Next() {
		var t TodoDB
		var todo Todo
		err := rows.Scan(&t.ID, &t.Title, &t.Status)
		if err != nil {
			log.Fatal("can't Scan row into variable", err)
		}

		todo = mapResponse(t)

		todos = append(todos, todo)
	}
	fmt.Println(todos)
	ctx.JSON(http.StatusOK, todos)
}

func getTodoByIdHandler(ctx *gin.Context) {
	rowId := ctx.Param("id")
	q := "SELECT id, title, status FROM todos where id=$1"
	row := db.QueryRow(q, rowId)

	var t TodoDB
	var todo Todo
	err := row.Scan(&t.ID, &t.Title, &t.Status)
	if err != nil {
		log.Fatal("can't Scan row into variable", err)
	}

	todo = mapResponse(t)

	ctx.JSON(http.StatusOK, todo)
}

func addTodoHandler(ctx *gin.Context) {
	var newTodo Todo

	if err := ctx.ShouldBindJSON(&newTodo); err != nil {
		ctx.Error(err)
	}

	q := "INSERT INTO todos (title, status) values ($1, $2)  RETURNING id"
	row := db.QueryRow(q, newTodo.Title, newTodo.Status)

	var id int
	err := row.Scan(&id)

	if err != nil {
		fmt.Println("can't scan id", err)
		return
	}

	fmt.Println("insert todo success id : ", id)
	ctx.JSON(http.StatusOK, id)

}

func editTodoHandler(ctx *gin.Context) {
	rowId, _ := strconv.Atoi(ctx.Param("id"))
	var editTodo Todo

	if err := ctx.ShouldBindJSON(&editTodo); err != nil {
		ctx.Error(err)
	}

	if _, err := db.Exec("UPDATE todos SET status=$2, title=$3 WHERE id=$1;", rowId, editTodo.Status, editTodo.Title); err != nil {
		log.Fatal("error execute update ", err)
	}

	fmt.Println("update success")

	ctx.JSON(http.StatusOK, Todo{ID: rowId, Status: editTodo.Status, Title: editTodo.Title})
}

func deleteTodoHandler(ctx *gin.Context) {
	rowId := ctx.Param("id")

	if _, err := db.Exec("DELETE FROM todos WHERE id=$1;", rowId); err != nil {
		log.Fatal("error delete todo", err)
	}

	fmt.Println("delete success")
	ctx.JSON(http.StatusOK, "Success")
}

func updateStatusHandler(ctx *gin.Context) {
	rowId := ctx.Param("id")

	var editTodo Todo

	if err := ctx.ShouldBindJSON(&editTodo); err != nil {
		ctx.Error(err)
	}

	if _, err := db.Exec("UPDATE todos SET status=$2 WHERE id=$1;", rowId, editTodo.Status); err != nil {
		log.Fatal("error patch status todo", err)
	}

	fmt.Println("patch status success")
	ctx.JSON(http.StatusOK, "Update Status Successful")
}

func updateTitleHandler(ctx *gin.Context) {
	rowId := ctx.Param("id")

	var editTodo Todo

	if err := ctx.ShouldBindJSON(&editTodo); err != nil {
		ctx.Error(err)
	}

	if _, err := db.Exec("UPDATE todos SET title=$2 WHERE id=$1;", rowId, editTodo.Title); err != nil {
		log.Fatal("error patch title todo", err)
	}

	fmt.Println("patch title success")
	ctx.JSON(http.StatusOK, "Update Title Successful")
}

func mapResponse(t TodoDB) Todo {
	var todo Todo
	todo.ID = t.ID
	todo.Title = t.Title.String
	todo.Status = t.Status.String
	return todo
}
