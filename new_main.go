package main
import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)
//noinspection ALL
const (
	DB_USER     = "postgres"
	DB_PASSWORD = "12345"
	DB_NAME     = "dvdrental"
)
//noinspection ALL
type details struct {
	new_id			string `json:"new_id"`
	name           string `json:"name"`
	source         string `json:"source"`
	phone_number   string `json:"phone_number"`
	experience      string `json:"experience"`
	ctc            string `json:"ctc"`
	ectc           string `json:"ectc"`
	np             string `json:"np"`
	status         string `json:"status"`
	interview_date string `json:"interview_date"`
	email          string `json:"email"`       //required
	applied_for    string `json:"applied_for"` //required`
}
type JsonResponse struct {
	Type    string    `json:"type"`
	Data    []details `json:"data"`
	Message string    `json:"message"`
}
func main() {
	router := mux.NewRouter()
	// Get all details
	router.HandleFunc("/details/", Getdetails).Methods("GET")
	// Create the details
	//Get Request to fetch a single data by the key - email
	router.HandleFunc("/details/", Createdetails).Methods("POST")
	// Delete a specific details by the detailsID
	router.HandleFunc("/details/{detail_email}", Getdetailsbyemail).Methods("POST")
	router.HandleFunc("/details/{detailsid}", Deletedetails).Methods("DELETE")
	// Delete all detailss
	router.HandleFunc("/details/{detailsid}", Updatedetails).Methods("PUT")
	router.HandleFunc("/details/", Deletedetailss).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8000", router))
}
// Get all details
//noinspection ALL
func Getdetails(w http.ResponseWriter, r *http.Request) {
	db := setupDB()
	printMessage("Getting all details...")
	// Get all details from details table that don't have detailsID = "1"
	rows, err := db.Query("SELECT * FROM details where new_id <> $1", "1")
	printMessage("fetching all details...")
	checkErr(err)
	var det []details
	for rows.Next() {
		var _name string
		var _source string
		var _phone_number string
		var _experience string
		var _ctc string
		var _ectc string
		var _np string
		var _status string
		var _interview_date string
		var _email string       //required
		var _applied_for string //required`
		err = rows.Scan(&_name, &_source, &_phone_number, &_experience, &_ctc, &_ectc, &_np, &_status, &_interview_date, &_email, &_applied_for)
		checkErr(err)
		det = append(det, details{name: _name, source: _source, phone_number: _phone_number, experience: _experience, ctc: _ctc, ectc: _ectc, np: _np, status: _status, interview_date: _interview_date, email: _email, applied_for: _applied_for})
	}
	var response = JsonResponse{Type: "success", Data: det}
	json.NewEncoder(w).Encode(response)
}

// Create a details
//noinspection ALL
func Createdetails(w http.ResponseWriter, r *http.Request) {
	details_name := r.FormValue("name")
	details_source := r.FormValue("source")
	details_phone_number := r.FormValue("phone_number")
	details_experience := r.FormValue("experience")
	details_ctc := r.FormValue("ctc")
	details_ectc := r.FormValue("ectc")
	details_np := r.FormValue("np")
	details_status := r.FormValue("status")
	details_interview_date := r.FormValue("interview_date")
	details_email := r.FormValue("email")
	details_applied_for := r.FormValue("applied_for")
	var response = JsonResponse{}
	if details_name == "" {
		response = JsonResponse{Type: "error", Message: "You are missing detailsName parameter."}
	} else {
		db := setupDB()
		printMessage("Inserting details into DB")
		var lastInsertID int
		err := db.QueryRow("INSERT INTO details(details_name,details_source,details_phone_number,details_experience,details_ctc,details_ectc,details_np,details_status,details_interview_date,details_email,details_applied_for) VALUES($1, $2,$3,$4,$5,$6,$7,$8,$9,$10) returning id;",details_name, details_source, details_phone_number, details_experience, details_ctc, details_ectc, details_np, details_status, details_interview_date, details_email, details_applied_for).Scan(&lastInsertID)
		checkErr(err)
		response = JsonResponse{Type: "success", Message: "The details has been inserted successfully!"}
	}
	json.NewEncoder(w).Encode(response)
}

// Delete a details//checked
func Deletedetails(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	detail_name := params["name"]
	var response = JsonResponse{}
	if detail_name == "" {
		response = JsonResponse{Type: "error", Message: "You are missing detailsID parameter."}
	} else {
		db := setupDB()
		printMessage("Deleting details from DB")
		_, err := db.Exec("DELETE FROM detailss where name = &1", detail_name)
		checkErr(err)
		response = JsonResponse{Type: "success", Message: "The details has been deleted successfully!"}
	}
	json.NewEncoder(w).Encode(response)
}

// Delete all detailss//chekecd
func Deletedetailss(w http.ResponseWriter, r *http.Request) {
	db := setupDB()
	printMessage("Deleting all detailss...")
	_, err := db.Exec("DELETE FROM details")
	checkErr(err)
	printMessage("All detailss have been deleted successfully!")
	var response = JsonResponse{Type: "success", Message: "All detailss have been deleted successfully!"}
	json.NewEncoder(w).Encode(response)
}

//Update a details
//noinspection ALL//checked all
func Updatedetails(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	details_name := params["name"]
	details_source := params["source"]
	details_phone_number := params["phone_number"]
	details_experience := params["experience"]
	details_ctc := params["ctc"]
	details_ectc := params["ectc"]
	details_np := params["np"]
	details_status := params["status"]
	details_interview_date := params["interview_date"]
	details_email := params["email"]
	details_applied_for := params["applied_for"]
	//db:=setupDB()
	var response = JsonResponse{}
	if details_name == "" {
	} else {
		db := setupDB()
		printMessage("Updating details in DB")
		_, err := db.Exec(`UPDATE details set name=$1,source=$2 ,phone_number=$3, experience=$4, ctc=$5 ,ectc=$6,np=$7,status=$8,interview_date=$9,email=$10,applied_for=$11,where name=$12 RETURNING name`, details_source, details_phone_number, details_experience, details_ctc, details_ectc, details_np, details_status, details_interview_date, details_email, details_applied_for)
		checkErr(err)
		response = JsonResponse{Type: "success", Message: "The details has been updated successfully!"}
	}
	json.NewEncoder(w).Encode(response)
}
//Get details for a single email
//noinspection ALL
func Getdetailsbyemail(w http.ResponseWriter, r *http.Request) {
	var response = JsonResponse{}
	printMessage("Getting details from an email.........")
	params := mux.Vars(r)
	var main_data []details
	detail_email := params["email"] //getting user query email
	if detail_email == "" {
		response = JsonResponse{Type: "error", Message: "You are missing email parameter."}
	} else {
		db := setupDB()
		printMessage("Getting details from DB by email")
		rows, err := db.Query("SELECT name FROM details where email = $1", detail_email)
		checkErr(err)
		for rows.Next() {
			var name string
			err = rows.Scan(&name)
			checkErr(err)
			main_data = append(main_data, details{name: name})
		}
	}
	response = JsonResponse{Type: "success", Data: main_data}
	json.NewEncoder(w).Encode(response)
}

// DB set up
func setupDB() *sql.DB {
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", DB_USER, DB_PASSWORD, DB_NAME)
	db, err := sql.Open("postgres", dbinfo)
	checkErr(err)
	return db
}

// Function for handling messages
func printMessage(message string) {
	fmt.Println("")
	fmt.Println(message)
	fmt.Println("")
}

// Function for handling errors
func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
