package migrations

import (
	"database/sql"
	"fmt"
	"log"
)

func Migrate(db *sql.DB) {
	userTable := `
    CREATE TABLE IF NOT EXISTS users (
        id SERIAL PRIMARY KEY,
        passport_serie INT NOT NULL,
        passport_number INT NOT NULL,
        surname VARCHAR(100) NOT NULL,
        name VARCHAR(100) NOT NULL,
        patronymic VARCHAR(100),
        address TEXT NOT NULL
    );`

	workLogTable := `
    CREATE TABLE IF NOT EXISTS work_logs (
        id SERIAL PRIMARY KEY,
        user_id INT NOT NULL,
        task_id INT NOT NULL,
        start_time TIMESTAMP NOT NULL,
        end_time TIMESTAMP,
        FOREIGN KEY (user_id) REFERENCES users(id)
    );`

	taskDescription := `
    ALTER TABLE tasks
    ADD COLUMN IF NOT EXISTS description TEXT;`

	_, err := db.Exec(userTable)
	if err != nil {
		log.Fatalf("Failed to execute migration for users table: %v", err)
	}

	_, err = db.Exec(workLogTable)
	if err != nil {
		log.Fatalf("Failed to execute migration for work_logs table: %v", err)
	}

	_, err = db.Exec(taskDescription)
	if err != nil {
		log.Fatalf("Failed to add description column to tasks table: %v", err)
	}

	fmt.Println("Migrations executed successfully")
}
