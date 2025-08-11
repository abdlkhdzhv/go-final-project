package api

import (
	"encoding/json"
	"fmt"
	"go-final-project/pkg/db"
	"go-final-project/pkg/utils"
	"io"
	"net/http"
	"strconv"
	"time"
)

type TasksResp struct {
	Tasks []map[string]string `json:"tasks"`
}

func TasksHandler(w http.ResponseWriter, r *http.Request) {
	tasks, err := db.Tasks(50)
	if err != nil {
		data := map[string]interface{}{
			"error": err.Error(),
		}
		jsonData, err := json.Marshal(data)
		if err != nil {
			fmt.Println(err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(jsonData)
		return
	}

	list := make([]map[string]string, 0, len(tasks))
	for _, t := range tasks {
		list = append(list, map[string]string{
			"id":      strconv.FormatInt(t.ID, 10),
			"date":    t.Date,
			"title":   t.Title,
			"comment": t.Comment,
			"repeat":  t.Repeat,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	writeJson(w, &TasksResp{list})
}

func GetTaskHundler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJson(w, map[string]string{"error": "Ожидается параметр id"})
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}
	writeJson(w, map[string]string{
		"id":      strconv.FormatInt(task.ID, 10),
		"date":    task.Date,
		"title":   task.Title,
		"comment": task.Comment,
		"repeat":  task.Repeat,
	})
}

func DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJson(w, map[string]string{"error": "not found id"})
		return
	}
	if err := db.DeleteTask(id); err != nil {
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}
	writeJson(w, map[string]any{})
}

func DoneTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJson(w, map[string]string{"error": "not found id"})
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}

	if task.Repeat == "" {
		if err := db.DeleteTask(id); err != nil {
			writeJson(w, map[string]string{"error": err.Error()})
			return
		}
		writeJson(w, map[string]any{})
		return
	}

	now := time.Now()
	next, err := utils.NextDate(now, task.Date, task.Repeat)
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}

	if err := db.UpdateDate(next, id); err != nil {
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}
	writeJson(w, map[string]any{})
}

func UpdateTaskHundler(w http.ResponseWriter, r *http.Request) {
	var payload map[string]any
	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}
	if err := json.Unmarshal(body, &payload); err != nil {
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}

	idStr := fmt.Sprint(payload["id"])
	if idStr == "" || idStr == "<nil>" {
		writeJson(w, map[string]string{"error": "Missing task ID"})
		return
	}
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeJson(w, map[string]string{"error": "Invalid task ID"})
		return
	}

	title := fmt.Sprint(payload["title"])
	if title == "" {
		writeJson(w, map[string]string{"error": "Empty title"})
		return
	}

	comment := fmt.Sprint(payload["comment"])
	repeat := fmt.Sprint(payload["repeat"])
	date := fmt.Sprint(payload["date"])

	now := time.Now()
	if date == "" {
		date = now.Format("20060102")
	}

	if _, err := time.Parse("20060102", date); err != nil {
		writeJson(w, map[string]string{"error": "Неверный формат даты"})
		return
	}

	if repeat != "" {
		nDate, err := utils.NextDate(now, date, repeat)
		if err != nil {
			writeJson(w, map[string]string{"error": err.Error()})
			return
		}
		date = nDate
	} else {
		dt, _ := time.Parse("20060102", date)
		if dt.Before(now) {
			date = now.Format("20060102")
		}
	}

	task := db.Task{
		ID:      id,
		Date:    date,
		Title:   title,
		Comment: comment,
		Repeat:  repeat,
	}

	if err := db.UpdateTask(&task); err != nil {
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}

	writeJson(w, map[string]string{
		"id":      idStr,
		"date":    task.Date,
		"title":   task.Title,
		"comment": task.Comment,
		"repeat":  task.Repeat,
	})
}
