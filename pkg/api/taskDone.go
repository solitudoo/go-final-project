package api

import (
	"go-final-project/pkg/db"
	"net/http"
	"time"
)

func taskDoneHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	task, err := db.GetTask(id)
	if err != nil {
		writeJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if task.Repeat == "" || task.Repeat == "0" {
		err = db.DeleteTask(id)
		if err != nil {
			writeJSONError(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}
	next, err := NextDate(time.Now().UTC(), task.Date, task.Repeat)
	if err != nil {
		writeJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = db.UpdateDate(next, id)
	if err != nil {
		writeJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, map[string]any{}, http.StatusOK)

}
