package api

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"time"

	"go-final-project/pkg/db"
)

func checkDate(task *db.Task) error {
	now := time.Now()
	today := now.Format("20060102")

	if task.Date == ""  {
		task.Date = today
		return nil
	}

	if _, err := time.Parse("20060102", task.Date); err != nil {
		return err
	}

	var next string
	if task.Repeat != "" {
		nowDate, _ := time.Parse("20060102", today)
		n, err := NextDate(nowDate, task.Date, task.Repeat)
		if err != nil {
			return err
		}
		next = n
	}

	if task.Date < today {
		if task.Repeat == "" {
			task.Date = today
		} else {
			task.Date = next
		}
	}

	return nil
}


func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var task db.Task
	if err := json.Unmarshal(body, &task); err != nil {
		writeJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	if task.Title == "" {
		writeJSONError(w, "task title is required", http.StatusBadRequest)
		return
	}

	if err := checkDate(&task); err != nil {
		writeJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := db.AddTask(&task)
	if err != nil {
		writeJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	task.ID = strconv.Itoa(int(id))
	writeJSON(w, task, http.StatusCreated)
}

func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		addTaskHandler(w, r)
	case http.MethodGet:
		getTaskHandler(w, r)
	case http.MethodPut:
		putTaskHandler(w, r)
	case http.MethodDelete:
		deleteTaskHandler(w, r)

	default:
		writeJSONError(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func writeJSON(w http.ResponseWriter, data any, status int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func writeJSONError(w http.ResponseWriter, msg string, status int) {
	writeJSON(w, map[string]string{
		"error": msg,
	}, status)
}

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSONError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSONError(w, "task id is required", http.StatusBadRequest)
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeJSONError(w, "server error", http.StatusInternalServerError)
		return
	}

	writeJSON(w, task, http.StatusOK)
	return

}
func putTaskHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	
	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var task db.Task
	if err := json.Unmarshal(body, &task); err != nil {
		writeJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	if task.Title == "" {
		writeJSONError(w, "task title is required", http.StatusBadRequest)
		return
	}

	if err := checkDate(&task); err != nil {
		writeJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = db.UpdateTask(&task)
	if err != nil {
		writeJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, task, http.StatusCreated)
}
func deleteTaskHandler(w http.ResponseWriter, r *http.Request){
	id := r.URL.Query().Get("id")
	if id == ""{
		writeJSONError(w, "task id is required", http.StatusBadRequest)
		return
	}
	err := db.DeleteTask(id)
	if err != nil {
    	writeJSONError(w, err.Error(), http.StatusNotFound) 
    	return
	}
	writeJSON(w, map[string]any{}, http.StatusOK)
	return
	
}  