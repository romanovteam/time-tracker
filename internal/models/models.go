package models

type User struct {
	ID             int    `json:"id"`
	PassportSerie  string `json:"passportSerie"`
	PassportNumber string `json:"passportNumber"`
	Surname        string `json:"surname"`
	Name           string `json:"name"`
	Patronymic     string `json:"patronymic,omitempty"`
	Address        string `json:"address"`
}

type WorkLog struct {
	TaskID int     `json:"taskID"`
	Hours  float64 `json:"hours"`
}

type Task struct {
	ID          int     `json:"id"`
	UserID      int     `json:"userID"`
	TaskID      int     `json:"taskID"`
	Description *string `json:"description,omitempty"`
}

type People struct {
	Surname    string `json:"surname"`
	Name       string `json:"name"`
	Patronymic string `json:"patronymic,omitempty"`
	Address    string `json:"address"`
}
