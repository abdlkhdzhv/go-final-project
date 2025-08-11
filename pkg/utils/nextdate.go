package utils

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	date, err := time.Parse("20060102", dstart)
	if err != nil {
		return "", fmt.Errorf("неверный формат даты: %w", err)
	}

	if repeat == "" {
		return "", errors.New("пустой интервал")
	}

	parts := strings.Fields(repeat)

	if len(parts) == 2 && parts[0] == "d" {
		interval, err := strconv.Atoi(parts[1])
		if err != nil {
			return "", fmt.Errorf("неверный интервал: %w", err)
		}
		if interval <= 0 {
			return "", errors.New("интервал должен быть больше нуля")
		}
		if interval >= 400 {
			return "", errors.New("превышен лимит интервала")
		}

		for !date.After(now) {
			date = date.AddDate(0, 0, interval)
		}
		return date.Format("20060102"), nil
	}

	if repeat == "y" {
		for !date.After(now) {
			date = date.AddDate(1, 0, 0)
		}
		return date.Format("20060102"), nil
	}

	return "", errors.New("неподдерживаемый формат правила")
}
