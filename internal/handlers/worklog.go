package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time-tracker/internal/db"
	"time-tracker/internal/models"
)

// GetWorkLogs godoc
// @Summary Get work logs for a user
// @Description Get work logs for a user within a specific period
// @Tags worklogs
// @Produce  json
// @Param user_id query int true "User ID"
// @Param start_date query string true "Start Date"
// @Param end_date query string true "End Date"
// @Success 200 {array} models.WorkLog
// @Router /worklogs [get]
func GetWorkLogs(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	startDate := r.URL.Query().Get("start_date")
	endDate := r.URL.Query().Get("end_date")

	query := `
        SELECT task_id, SUM(EXTRACT(EPOCH FROM (end_time - start_time))/3600) AS hours
        FROM work_logs
        WHERE user_id = $1 AND start_time >= $2 AND end_time <= $3
        GROUP BY task_id
        ORDER BY hours DESC`

	rows, err := db.Database.Query(query, userID, startDate, endDate)
	if err != nil {
		log.Printf("Failed to execute query: %v", err)
		http.Error(w, "Failed to execute query", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var workLogs []models.WorkLog
	for rows.Next() {
		var workLog models.WorkLog
		if err := rows.Scan(&workLog.TaskID, &workLog.Hours); err != nil {
			log.Printf("Failed to scan row: %v", err)
			http.Error(w, "Failed to scan row", http.StatusInternalServerError)
			return
		}
		workLogs = append(workLogs, workLog)
	}

	json.NewEncoder(w).Encode(workLogs)
}
