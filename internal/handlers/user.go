package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time-tracker/internal/db"
	"time-tracker/internal/models"

	"github.com/gorilla/mux"
)

// GetUsers godoc
// @Summary Get all users with filtering and pagination
// @Description Get all users with filtering and pagination
// @Tags users
// @Produce  json
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Param filter query string false "Filter"
// @Success 200 {array} models.User
// @Router /users [get]
func GetUsers(w http.ResponseWriter, r *http.Request) {
	limit := r.URL.Query().Get("limit")
	offset := r.URL.Query().Get("offset")
	filter := r.URL.Query().Get("filter")

	query := "SELECT * FROM users WHERE 1=1"
	if filter != "" {
		query += " AND (name LIKE '%' || $1 || '%' OR surname LIKE '%' || $1 || '%')"
	}
	if limit != "" {
		query += " LIMIT " + limit
	}
	if offset != "" {
		query += " OFFSET " + offset
	}

	rows, err := db.Database.Query(query, filter)
	if err != nil {
		http.Error(w, "Failed to execute query", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.PassportSerie, &user.PassportNumber, &user.Surname, &user.Name, &user.Patronymic, &user.Address); err != nil {
			http.Error(w, "Failed to scan row", http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}

	json.NewEncoder(w).Encode(users)
}

// DeleteUser godoc
// @Summary Delete a user by ID
// @Description Delete a user by ID
// @Tags users
// @Produce  json
// @Param id path int true "User ID"
// @Success 200 {string} string "Successfully deleted"
// @Router /users/{id} [delete]
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]

	query := `DELETE FROM users WHERE id = $1`
	res, err := db.Database.Exec(query, userID)
	if err != nil {
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		http.Error(w, "Failed to retrieve rows affected", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "No user found to delete", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Successfully deleted"))
}

// UpdateUser godoc
// @Summary Update a user's information
// @Description Update a user's information
// @Tags users
// @Accept  json
// @Produce  json
// @Param id path int true "User ID"
// @Param user body models.User true "User"
// @Success 200 {object} models.User
// @Router /users/{id} [put]
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]

	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if !ValidatePassportSerieAndNumber(user.PassportSerie, user.PassportNumber) {
		http.Error(w, "Invalid passport series or number format", http.StatusBadRequest)
		return
	}

	query := `UPDATE users SET passport_serie = $1, passport_number = $2, surname = $3, name = $4, patronymic = $5, address = $6 WHERE id = $7`
	_, err := db.Database.Exec(query, user.PassportSerie, user.PassportNumber, user.Surname, user.Name, user.Patronymic, user.Address, userID)
	if err != nil {
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

// AddUser godoc
// @Summary Add a new user
// @Description Add a new user and enrich data from external API
// @Tags users
// @Accept  json
// @Produce  json
// @Param user body models.User true "User"
// @Success 201 {object} models.User
// @Router /users [post]
func AddUser(w http.ResponseWriter, r *http.Request) {
	var user struct {
		PassportNumber string `json:"passportNumber"`
	}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if !ValidatePassportNumber(user.PassportNumber) {
		http.Error(w, "Invalid passport number format", http.StatusBadRequest)
		return
	}

	passportParts := strings.Split(user.PassportNumber, " ")
	passportSerie := passportParts[0]
	passportNumber := passportParts[1]

	apiURL := fmt.Sprintf("%s?passportSerie=%s&passportNumber=%s", os.Getenv("API_URL"), passportSerie, passportNumber)
	resp, err := http.Get(apiURL)
	var externalUser models.People

	if err != nil || resp.StatusCode != http.StatusOK {
		log.Println("Failed to fetch user data from external API, creating user with default values")
		externalUser = models.People{
			Surname:    "Unknown",
			Name:       "Unknown",
			Patronymic: "Unknown",
			Address:    "Unknown",
		}
	} else {
		defer resp.Body.Close()
		err = json.NewDecoder(resp.Body).Decode(&externalUser)
		if err != nil {
			log.Println("Failed to parse external API response, creating user with default values")
			externalUser = models.People{
				Surname:    "Unknown",
				Name:       "Unknown",
				Patronymic: "Unknown",
				Address:    "Unknown",
			}
		}
	}

	query := `INSERT INTO users (passport_serie, passport_number, surname, name, patronymic, address) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err = db.Database.Exec(query, passportSerie, passportNumber, externalUser.Surname, externalUser.Name, externalUser.Patronymic, externalUser.Address)
	if err != nil {
		http.Error(w, "Failed to insert user data into the database", http.StatusInternalServerError)
		return
	}

	newUser := models.User{
		PassportSerie:  passportSerie,
		PassportNumber: passportNumber,
		Surname:        externalUser.Surname,
		Name:           externalUser.Name,
		Patronymic:     externalUser.Patronymic,
		Address:        externalUser.Address,
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newUser)
}
