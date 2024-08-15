package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time-tracker/internal/db"
	"time-tracker/internal/models"
)

// StartTask godoc
// @Summary Start a task for a user
// @Description Start a task for a user
// @Tags tasks
// @Accept  json
// @Produce  json
// @Param task body models.Task true "Task"
// @Success 200 {object} models.Task
// @Router /tasks/start [post]
func StartTask(w http.ResponseWriter, r *http.Request) {
	var task models.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	query := `INSERT INTO work_logs (user_id, task_id, start_time, description) VALUES ($1, $2, NOW(), $3) RETURNING id`
	err := db.Database.QueryRow(query, task.UserID, task.TaskID, task.Description).Scan(&task.ID)
	if err != nil {
		log.Printf("Failed to start task: %v", err)
		http.Error(w, "Failed to start task", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(task)
}
