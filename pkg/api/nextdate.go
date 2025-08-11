package api

import (
	"fmt"
	"go-final-project/pkg/utils"
	"net/http"
	"time"
)

const format = "20060102"

func nextDayHandler(w http.ResponseWriter, r *http.Request) {
	now := r.FormValue("now")
	date := r.FormValue("date")
	repeat := r.FormValue("repeat")

	nowDate, err := time.Parse(format, now)
	if err != nil {
		fmt.Println(err)
		return
	}

	if date == "" || repeat == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if now == "" {
		now = time.Now().UTC().Format(format)
	}

	nextDate, err := utils.NextDate(nowDate, date, repeat)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(nextDate))
}
