package main

import (
	"fmt"
	"net/http"

	"task-manager/internal/handler"
	"task-manager/internal/repository"
)

func main() {
	repo := repository.NewTaskRepository()
	taskHandler := handler.NewTaskHandler(repo)

	http.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./web/index.html")
	})

	http.HandleFunc("GET /tasks", taskHandler.GetTasks)
	http.HandleFunc("POST /tasks", taskHandler.CreateTask)
	http.HandleFunc("DELETE /tasks/{id}", taskHandler.DeleteTask)
	http.HandleFunc("PUT /tasks/{id}", taskHandler.UpdateTask)

	fmt.Println("server started on http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
