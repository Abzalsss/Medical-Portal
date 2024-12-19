package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

var db *sql.DB

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
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

	query := `
		INSERT INTO doctors (user_id, specialty, experience, rating, bio, schedule)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING doctor_id`
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

// Add appointment
func addAppointmentHandler(w http.ResponseWriter, r *http.Request) {
	var appointment struct {
		PatientID       int    `json:"patient_id"`
		DoctorID        int    `json:"doctor_id"`
		AppointmentDate string `json:"appointment_date"`
		Status          string `json:"status"`
	}
	err := json.NewDecoder(r.Body).Decode(&appointment)
	if err != nil || appointment.PatientID == 0 || appointment.DoctorID == 0 || appointment.AppointmentDate == "" {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	query := `
		INSERT INTO appointments (patient_id, doctor_id, appointment_date, status)
		VALUES ($1, $2, $3, $4)
		RETURNING appointment_id`
	var appointmentID int
	err = db.QueryRow(query, appointment.PatientID, appointment.DoctorID, appointment.AppointmentDate, appointment.Status).Scan(&appointmentID)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Response{
		Status:  "success",
		Message: fmt.Sprintf("Appointment added with ID %d", appointmentID),
	})
}

// Add review
func addReviewHandler(w http.ResponseWriter, r *http.Request) {
	var review struct {
		PatientID int    `json:"patient_id"`
		DoctorID  int    `json:"doctor_id"`
		Rating    int    `json:"rating"`
		Comment   string `json:"comment"`
	}
	err := json.NewDecoder(r.Body).Decode(&review)
	if err != nil || review.PatientID == 0 || review.DoctorID == 0 || review.Rating < 1 || review.Rating > 5 {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	query := `
		INSERT INTO reviews (patient_id, doctor_id, rating, comment)
		VALUES ($1, $2, $3, $4)
		RETURNING review_id`
	var reviewID int
	err = db.QueryRow(query, review.PatientID, review.DoctorID, review.Rating, review.Comment).Scan(&reviewID)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Response{
		Status:  "success",
		Message: fmt.Sprintf("Review added with ID %d", reviewID),
	})
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

	http.HandleFunc("/add-user", addUserHandler)
	http.HandleFunc("/add-doctor", addDoctorHandler)
	http.HandleFunc("/add-appointment", addAppointmentHandler)
	http.HandleFunc("/add-review", addReviewHandler)

	log.Println("Server started on 127.0.0.1:8080")
	err = http.ListenAndServe("127.0.0.1:8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
