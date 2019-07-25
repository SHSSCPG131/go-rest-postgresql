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
type new_details struct {
	Id			int `json:"id"`
	Name           string `json:"name"`
	Source         string `json:"source"`
	Phone_number   string `json:"phone_number"`
	Experience      string `json:"experience"`
	Ctc            string `json:"ctc"`
	Ectc           string `json:"ectc"`
	Np             string `json:"np"`
	Status         string `json:"status"`
	Interview_date string `json:"interview_date"`
	Email          string `json:"email"`       //required
	Applied_for    string `json:"applied_for"` //required`
}
type JsonResponse struct {
	Type    string    `json:"type"`
	Data    []new_details `json:"data"`
	Message string    `json:"message"`
}
func main() {
	router := mux.NewRouter()
	// Get all new_details
	router.HandleFunc("/new_details/", Getnew_details).Methods("GET")
	// Create the new_details
	//Get Request to fetch a single data by the key - email
	router.HandleFunc("/new_details/", Createnew_details).Methods("POST")
	// Delete a specific new_details by the new_detailsID
	router.HandleFunc("/new_details/{detail_email}", Getnew_detailsbyemail).Methods("POST")
	//router.HandleFunc("/new_details/{new_detailsid}", Deletenew_details).Methods("DELETE")
	// Delete all new_detailss
	router.HandleFunc("/new_details/{new_detailsid}", Updatenew_details).Methods("PUT")
	router.HandleFunc("/new_details/", Deletenew_detailss).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8000", router))
}
// Get all new_details
//noinspection ALL
func Getnew_details(w http.ResponseWriter, r *http.Request) {
	db := setupDB()
	printMessage("Getting all new_details...")
	// Get all new_details from new_details table that don't have new_detailsID = "1"
	rows, err := db.Query("SELECT * FROM new_details where id <> $1", "1")
	printMessage("fetching all new_details...")
	checkErr(err)
	var det []new_details
	for rows.Next() {
		var _id int
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
		err = rows.Scan(&_id,&_name, &_source, &_phone_number, &_experience, &_ctc, &_ectc, &_np, &_status, &_interview_date, &_email, &_applied_for)
		checkErr(err)
		det = append(det, new_details{Id:_id,Name: _name, Source: _source, Phone_number: _phone_number, Experience: _experience, Ctc: _ctc, Ectc: _ectc, Np: _np, Status: _status, Interview_date: _interview_date, Email: _email, Applied_for: _applied_for})
	}
	var response = JsonResponse{Type: "success", Data: det}
	json.NewEncoder(w).Encode(response)
}
// Create a new_details
func Createnew_details(w http.ResponseWriter, r *http.Request) {
	new_details_id:=r.FormValue("id")
	new_details_name := r.FormValue("name")
	new_details_source := r.FormValue("source")
	new_details_phone_number := r.FormValue("phone_number")
	new_details_experience := r.FormValue("experience")
	new_details_ctc := r.FormValue("ctc")
	new_details_ectc := r.FormValue("ectc")
	new_details_np := r.FormValue("np")
	new_details_status := r.FormValue("status")
	new_details_interview_date := r.FormValue("interview_date")
	new_details_email := r.FormValue("email")
	new_details_applied_for := r.FormValue("applied_for")
	var response = JsonResponse{}
	if new_details_name == "" || new_details_id==""{
		response = JsonResponse{Type: "error", Message: "You are missing Name and ID parameter."}
	} else {
		db := setupDB()
		printMessage("Inserting new_details into DB" + new_details_id)
		var lastInsertID int
		err := db.QueryRow("INSERT INTO new_details(new_details_id,new_details_email,new_details_applied_for,new_details_name,new_details_source,new_details_phone_number,new_details_experience,new_details_ctc,new_details_ectc,new_details_np,new_details_status,new_details_interview_date) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12) returning id;",new_details_id,new_details_email,new_details_applied_for,new_details_name, new_details_source, new_details_phone_number, new_details_experience, new_details_ctc, new_details_ectc, new_details_np, new_details_status, new_details_interview_date).Scan(&lastInsertID)
		checkErr(err)
		response = JsonResponse{Type: "success", Message: "The new_details has been inserted successfully!"}
	}
	json.NewEncoder(w).Encode(response)
}
// Delete a new_details//checked
func Deletenew_details(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	detail_name := params["name"]
	var response = JsonResponse{}
	if detail_name == "" {
		response = JsonResponse{Type: "error", Message: "You are missing new_detailsID parameter."}
	} else {
		db := setupDB()
		printMessage("Deleting new_details from DB")
		_, err := db.Exec("DELETE FROM new_detailss where name = &1", detail_name)
		checkErr(err)
		response = JsonResponse{Type: "success", Message: "The new_details has been deleted successfully!"}
	}
	json.NewEncoder(w).Encode(response)
}
// Delete all new_detailss//chekecd
func Deletenew_detailss(w http.ResponseWriter, r *http.Request) {
	db := setupDB()
	printMessage("Deleting all new_detailss...")
	_, err := db.Exec("DELETE FROM new_details")
	checkErr(err)
	printMessage("All new_detailss have been deleted successfully!")
	var response = JsonResponse{Type: "success", Message: "All new_detailss have been deleted successfully!"}
	json.NewEncoder(w).Encode(response)
}
//Update a new_details
//noinspection ALL//checked all
func Updatenew_details(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	new_details_name := params["name"]
	new_details_source := params["source"]
	new_details_phone_number := params["phone_number"]
	new_details_experience := params["experience"]
	new_details_ctc := params["ctc"]
	new_details_ectc := params["ectc"]
	new_details_np := params["np"]
	new_details_status := params["status"]
	new_details_interview_date := params["interview_date"]
	new_details_email := params["email"]
	new_details_applied_for := params["applied_for"]
	//db:=setupDB()
	var response = JsonResponse{}
	if new_details_name == "" {
	} else {
		db := setupDB()
		printMessage("Updating new_details in DB")
		_, err := db.Exec(`UPDATE new_details set name=$1,source=$2 ,phone_number=$3, experience=$4, ctc=$5 ,ectc=$6,np=$7,status=$8,interview_date=$9,email=$10,applied_for=$11,where name=$12 RETURNING name`, new_details_source, new_details_phone_number, new_details_experience, new_details_ctc, new_details_ectc, new_details_np, new_details_status, new_details_interview_date, new_details_email, new_details_applied_for)
		checkErr(err)
		response = JsonResponse{Type: "success", Message: "The new_details has been updated successfully!"}
	}
	json.NewEncoder(w).Encode(response)
}
//Get new_details for a single email
//noinspection ALL
func Getnew_detailsbyemail(w http.ResponseWriter, r *http.Request) {
	var response = JsonResponse{}
	printMessage("Getting new_details from an email.........")
	params := mux.Vars(r)
	var main_data []new_details
	detail_email := params["email"] //getting user query email
	if detail_email == "" {
		response = JsonResponse{Type: "error", Message: "You are missing email parameter."}
	} else {
		db := setupDB()
		printMessage("Getting new_details from DB by email")
		rows, err := db.Query("SELECT name FROM new_details where email = $1", detail_email)
		checkErr(err)
		for rows.Next() {
			var _name string
			err = rows.Scan(&_name)
			checkErr(err)
			main_data = append(main_data, new_details{Name: _name})
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

