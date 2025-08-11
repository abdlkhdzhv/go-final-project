package api

import (
	"net/http"
)

func Init() {
	http.HandleFunc("/api/nextdate", nextDayHandler)
	http.HandleFunc("/api/task", TaskHandler)
	http.HandleFunc("/api/tasks", TasksHandler)
	http.HandleFunc("/api/task/done", DoneTaskHandler)
}

func TaskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		AddTaskHandler(w, r)
	case http.MethodGet:
		GetTaskHundler(w, r)
	case http.MethodPut:
		UpdateTaskHundler(w, r)
	case http.MethodDelete:
		DeleteTaskHandler(w, r)
	default:
		w.WriteHeader(http.StatusBadRequest)
	}
}
