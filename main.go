package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"strconv"
)

var db *sql.DB

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// User struct to store user data
type User struct {
	UserID int    `json:"user_id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Role   string `json:"role"`
}

// Doctor struct
type Doctor struct {
	DoctorID   int     `json:"doctor_id"`
	UserID     int     `json:"user_id"`
	Specialty  string  `json:"specialty"`
	Experience int     `json:"experience"`
	Rating     float64 `json:"rating"`
	Bio        string  `json:"bio"`
	Schedule   string  `json:"schedule"`
}

// Add user
func addUserHandler(w http.ResponseWriter, r *http.Request) {
	var user struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil || user.Name == "" || user.Email == "" || user.Password == "" || user.Role == "" {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	query := `INSERT INTO users (name, email, password, role) VALUES ($1, $2, $3, $4) RETURNING user_id`
	var userID int
	err = db.QueryRow(query, user.Name, user.Email, user.Password, user.Role).Scan(&userID)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Response{
		Status:  "success",
		Message: fmt.Sprintf("User %s added with ID %d", user.Name, userID),
	})
}

// Get all users
func getUsersHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT user_id, name, email, role FROM users")
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.UserID, &user.Name, &user.Email, &user.Role)
		if err != nil {
			http.Error(w, "Error reading data", http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// Update user
func updateUserHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil || user.UserID == 0 || user.Name == "" || user.Email == "" || user.Role == "" {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	query := `UPDATE users SET name=$1, email=$2, role=$3 WHERE user_id=$4`
	_, err = db.Exec(query, user.Name, user.Email, user.Role, user.UserID)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Response{
		Status:  "success",
		Message: fmt.Sprintf("User %d updated", user.UserID),
	})
}

// Delete user
func deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	var user struct {
		UserID int `json:"user_id"`
	}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil || user.UserID == 0 {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	query := `DELETE FROM users WHERE user_id=$1`
	_, err = db.Exec(query, user.UserID)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Response{
		Status:  "success",
		Message: fmt.Sprintf("User %d deleted", user.UserID),
	})
}

// Add doctor
func addDoctorHandler(w http.ResponseWriter, r *http.Request) {
	var doctor struct {
		UserID     int     `json:"user_id"`
		Specialty  string  `json:"specialty"`
		Experience int     `json:"experience"`
		Rating     float64 `json:"rating"`
		Bio        string  `json:"bio"`
		Schedule   string  `json:"schedule"`
	}
	err := json.NewDecoder(r.Body).Decode(&doctor)
	if err != nil || doctor.UserID == 0 || doctor.Specialty == "" {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	query := `INSERT INTO doctors (user_id, specialty, experience, rating, bio, schedule)
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING doctor_id`
	var doctorID int
	err = db.QueryRow(query, doctor.UserID, doctor.Specialty, doctor.Experience, doctor.Rating, doctor.Bio, doctor.Schedule).Scan(&doctorID)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Response{
		Status:  "success",
		Message: fmt.Sprintf("Doctor added with ID %d", doctorID),
	})
}

// Update doctor
func updateDoctorHandler(w http.ResponseWriter, r *http.Request) {
	var doctor Doctor
	err := json.NewDecoder(r.Body).Decode(&doctor)
	if err != nil || doctor.DoctorID == 0 || doctor.Specialty == "" {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	query := `UPDATE doctors SET specialty=$1, experience=$2, rating=$3, bio=$4, schedule=$5 WHERE doctor_id=$6`
	_, err = db.Exec(query, doctor.Specialty, doctor.Experience, doctor.Rating, doctor.Bio, doctor.Schedule, doctor.DoctorID)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Response{
		Status:  "success",
		Message: fmt.Sprintf("Doctor %d updated", doctor.DoctorID),
	})
}

// Delete doctor
func deleteDoctorHandler(w http.ResponseWriter, r *http.Request) {
	var doctor struct {
		DoctorID int `json:"doctor_id"`
	}
	err := json.NewDecoder(r.Body).Decode(&doctor)
	if err != nil || doctor.DoctorID == 0 {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	query := `DELETE FROM doctors WHERE doctor_id=$1`
	_, err = db.Exec(query, doctor.DoctorID)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Response{
		Status:  "success",
		Message: fmt.Sprintf("Doctor %d deleted", doctor.DoctorID),
	})
}

// Get user by ID
func getUserByIdHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем user_id из параметров запроса
	userId := r.URL.Query().Get("user_id")
	if userId == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	// Преобразуем user_id в integer
	id, err := strconv.Atoi(userId)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	query := `SELECT user_id, name, email, role FROM users WHERE user_id=$1`
	var user User
	err = db.QueryRow(query, id).Scan(&user.UserID, &user.Name, &user.Email, &user.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			http.Error(w, "Database error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func main() {
	connStr := "user=postgres password=0000 dbname=medical_portal sslmode=disable"
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Serve static files
	fs := http.FileServer(http.Dir("./medical-portal"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Serve the HTML file
	http.HandleFunc("/", servehtml)

	// User-related routes
	http.HandleFunc("/add-user", addUserHandler)
	http.HandleFunc("/get-users", getUsersHandler)
	http.HandleFunc("/update-user", updateUserHandler)
	http.HandleFunc("/delete-user", deleteUserHandler)

	// Doctor-related routes
	http.HandleFunc("/add-doctor", addDoctorHandler)
	http.HandleFunc("/update-doctor", updateDoctorHandler)
	http.HandleFunc("/delete-doctor", deleteDoctorHandler)

	// Get user by ID
	http.HandleFunc("/get-user", getUserByIdHandler)

	log.Println("Server started on 127.0.0.1:8080")
	err = http.ListenAndServe("127.0.0.1:8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func servehtml(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "medical-portal/index.html")
}
