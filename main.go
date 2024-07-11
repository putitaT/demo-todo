package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/putitaT/todo-backend/database"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	r := gin.Default()
	// r.LoadHTMLGlob("./*.html")
	// r.GET("/", func(c *gin.Context) {
	// 	c.HTML(http.StatusOK, "index.html", nil)
	// })

	r.GET("/todos", getTodoHandler)

	// r.POST("/todos", postTodoHandler)

	// r.DELETE("/todos/:id", deleteToDoHandler)

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

func getTodoHandler(ctx *gin.Context) {
	db := database.ConnectDB()
	_ = db
}
