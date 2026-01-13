package api

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	date, err := time.Parse("20060102", dstart)
	if err != nil {
		return "", err
	}
	date = date.UTC()
	now = now.UTC()
	repeat = strings.TrimSpace(repeat)
	if repeat == "" {
		return "", errors.New("Неправильные данные")
	}
	array := strings.Split(repeat, " ")

	symbol := array[0]

	switch symbol {
	case "d":
		if len(array) != 2 {
			return "", errors.New("Неправильные данные")
		}
		interval, err := strconv.Atoi(array[1])
		if err != nil {
			return "", err
		}
		if interval <= 0 || interval > 400 {
			return "", errors.New("Недопустимое значение")
		}
		if date.After(now) {
			return date.Format("20060102"), nil
		}
		for {
			date = date.AddDate(0, 0, interval)
			if date.After(now) {
				break
			}
		}

	case "y":
		if len(array) != 1 {
			return "", errors.New("Неправильные данные")
		}
		if date.After(now) {
			return date.Format("20060102"), nil
		}
		for {
			date = date.AddDate(1, 0, 0)
			if date.After(now) {
				break
			}
		}
	default:
		return "", errors.New("Неправильные данные")
	}

	return date.Format("20060102"), nil

}

func nextDayHandler(res http.ResponseWriter, req *http.Request) {
	nowStr := req.FormValue("now")
	date := req.FormValue("date")
	repeat := req.FormValue("repeat")

	var now time.Time
	if nowStr == "" {
		now = time.Now().UTC()
		nowStr = now.Format("20060102")
	}

	now, err := time.Parse("20060102", nowStr)

	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	nextDate, err := NextDate(now, date, repeat)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	res.WriteHeader(http.StatusOK)
	res.Write([]byte(nextDate))
}
