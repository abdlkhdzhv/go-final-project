package api

import (
	"encoding/json"
	"go-final-project/pkg/db"
	"go-final-project/pkg/utils"
	"net/http"
	"strconv"
	"time"
)

func AddTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}

	if task.Title == "" {
		writeJson(w, map[string]string{"error": "Не указан заголовок задачи"})
		return
	}

	now := time.Now()

	if task.Date == "" {
		task.Date = now.Format("20060102")
	}

	t, err := time.Parse("20060102", task.Date)
	if err != nil {
		writeJson(w, map[string]string{"error": "Неверный формат даты"})
		return
	}

	if t.Before(now) {
		if task.Repeat == "" {
			task.Date = now.Format("20060102")
		} else {
			next, err := utils.NextDate(now, task.Date, task.Repeat)
			if err != nil {
				writeJson(w, map[string]string{"error": err.Error()})
				return
			}
			task.Date = next
		}
	} else {
		if task.Repeat != "" {
			if _, err := utils.NextDate(now, task.Date, task.Repeat); err != nil {
				writeJson(w, map[string]string{"error": err.Error()})
				return
			}
		}
	}

	id, err := db.AddTask(&task)
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}

	writeJson(w, map[string]string{"id": strconv.FormatInt(id, 10)})
}

func writeJson(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(data)
}
