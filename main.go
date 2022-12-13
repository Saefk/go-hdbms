package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

//type Patient

type Patient struct {
	ID        int
	Firstname string
	Lastname  string
	Ward      string
}

type Doctor struct {
	ID                 int
	DateOfBirth        string
	NationalID         string
	RegistrationNumber string
	Rank               string
	Department         string
	DirectReport       string
	PhoneNumber        string
	Email              string
	Address            string
}

type Case struct {
	ID                   int
	PatientID            int
	RegisterDate         string
	TriageID             int
	PresentingComplaint  string
	History              string
	PastMedicalHistory   string
	DrugHistory          string
	Allergies            string
	Examination          string
	Investigation        string
	ProvisionalDiagnosis string
	DischargeDate        string
	CurrentWard          string
	DateCreated          string
	CreatedBy            int
}

type CaseNote struct {
	ID                  int
	CaseID              int
	PatientID           int
	DoctorID            int
	DateCreated         string
	LastEdited          string
	Editors             string
	CaseType            string
	PresentingComplaint string
	History             string
	Examination         string
	Investigation       string
	PastMedicalHistory  string
	DrugHistory         string
	Diagnoses           string
	CurrentManagement   string
}

var tpl *template.Template

var db *sql.DB

const DBPASS = "#Okibgnina321"

func main() {

	var err error
	// this will only parse templates that match the templates/*. patterhtml
	tpl, err = template.ParseGlob("templates/*.html")
	if err != nil {
		fmt.Println("template err")
		panic(err.Error())
	}
	//sqlOpen takes in driver name + data source name
	//returns type DB struct
	db, err = sql.Open("mysql", "root:"+DBPASS+"@tcp(localhost:3306)/newdatabase")
	if err != nil {
		// fmt.Println("error validating sql.Open arguments")
		panic(err.Error())
	}
	defer db.Close()
	http.HandleFunc("/patient/name", patientSearchByNameHandler)
	http.HandleFunc("/patient/ward", patientSearchbyWardHandler)
	http.HandleFunc("/patient/register", patientRegistrationHandler)
	http.HandleFunc("/patient/remove/", removePatientByIdHandler)
	http.HandleFunc("/patient/update/redirect/", updateRedirectHandler)
	http.HandleFunc("/patient/update/details/", updatePatientDetailsHandler)
	http.HandleFunc("/", homePageHandler)

	http.HandleFunc("/case/new", caseAddHandler)
	http.ListenAndServe("localhost:8080", nil)

	// check if connected to DB
	// err = db.Ping()
	// if err != nil {
	// 	fmt.Println("error verifying with connection db.Ping")
	// 	panic(err.Error())

	// }

	// insert, err := db.Query("INSERT INTO `newdatabase`.`patients` (`firstname`,`lastname`,`id`) VALUES('jimmy','D','2')")
	// if err != nil {
	// 	panic(err.Error())
	// }
	// defer insert.Close()

	// fmt.Println("Connection Successful")
}

const Redir = 307

func updateRedirectHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	id := r.FormValue("patientId")
	row := db.QueryRow("Select * FROM newdatabase.patients WHERE id = ?;", id)
	var p Patient
	err := row.Scan(&p.ID, &p.Firstname, &p.Lastname, &p.Ward)
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/patient/name", Redir)
		return
	}
	tpl.ExecuteTemplate(w, "patient/patientupdate.html", p)
}

func updatePatientDetailsHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	id := r.FormValue("patientId")
	firstname := r.FormValue("firstName")
	lastname := r.FormValue("lastName")
	ward := r.FormValue("WardName")
	// upStmt := "UPDATE `newdatabase`.`patients` SET `firstname` = ?,`lastname`= ?,`ward`= ? WHERE (`id` = ?);"
	upStmt := "UPDATE `newdatabase`.`patients` SET `firstname` = ?, `lastname` = ?, `ward` = ? WHERE (`id` = ?);"

	stmt, err := db.Prepare(upStmt)
	if err != nil {
		fmt.Println("Error preparing stmt")
		panic(err)
	}

	fmt.Println("db.Prepare err: ", err)
	fmt.Println("db.Prepare stmt", stmt)

	defer stmt.Close()

	var res sql.Result
	res, err = stmt.Exec(firstname, lastname, ward, id)
	rowsAff, _ := res.RowsAffected()
	if err != nil || rowsAff != 1 {
		print(rowsAff)
		fmt.Println(err)
		tpl.ExecuteTemplate(w, "patient/result.html", "Error when attempting to update details")
		return
	}

	tpl.ExecuteTemplate(w, "patient/result.html", "Patient details successfully updated.")
}

func patientSearchByNameHandler(w http.ResponseWriter, r *http.Request) {
	//if method is get => write template
	if r.Method == "GET" {
		tpl.ExecuteTemplate(w, "patinet/patientbyname.html", nil)
		return
	}
	r.ParseForm()
	firstname := r.FormValue("firstName")
	lastname := r.FormValue("lastName")
	// fmt.Println("Firstname", firstname)
	stmt := "SELECT * FROM patients where `firstname` = ? AND `lastname` = ?;"
	//Query row returns one row
	rows, err := db.Query(stmt, firstname, lastname)
	if err != nil {
		panic(err)
	}

	defer rows.Close()
	var patients []Patient
	for rows.Next() {
		var p Patient
		err := rows.Scan(&p.ID, &p.Firstname, &p.Lastname, &p.Ward)
		if err != nil {
			panic(err)
		}
		patients = append(patients, p)
	}
	// err := row.Scan(&P.ID, &P.Firstname, &P.Lastname, &P.Ward)

	tpl.ExecuteTemplate(w, "patient/patientbyname.html", patients)
}

func removePatientByIdHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	id := r.FormValue("patientId")

	del, err := db.Prepare("DELETE FROM `newdatabase`.`patients` WHERE (`id` =?);")
	if err != nil {
		panic(err)
	}
	defer del.Close()

	var res sql.Result
	res, err = del.Exec(id)
	rowsAff, _ := res.RowsAffected()
	fmt.Println("rowsAff:", rowsAff)
	if err != nil || rowsAff != 1 {
		fmt.Fprint(w, "Error deleting patient from register.")
		return
	}
	fmt.Println("err:", err)
	tpl.ExecuteTemplate(w, "patient/result.html", "Patient successfully deletetd from register.")

}

func patientSearchbyWardHandler(w http.ResponseWriter, r *http.Request) {
	//if method is get => write template
	if r.Method == "GET" {
		tpl.ExecuteTemplate(w, "patient/patientsearchbyward.html", nil)
		return
	}
	r.ParseForm()

	name := r.FormValue("WardName")
	fmt.Println("Ward: ", name)
	stmt := "SELECT * FROM patients where ward = ?;"
	//Query row returns one row
	rows, err := db.Query(stmt, name)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	var patients []Patient
	for rows.Next() {
		var p Patient

		err = rows.Scan(&p.ID, &p.Firstname, &p.Lastname, &p.Ward)
		if err != nil {
			panic(err)
		}
		patients = append(patients, p)
	}
	tpl.ExecuteTemplate(w, "patient/patientsearchbyward.html", patients)
}

func patientRegistrationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tpl.ExecuteTemplate(w, "patient/patientregister.html", nil)
		return
	}
	firstname := r.FormValue("firstName")
	lastname := r.FormValue("lastName")
	ward := r.FormValue("wardName")
	var err error

	if firstname == "" || lastname == "" {
		fmt.Println("Error registering", err)
		tpl.ExecuteTemplate(w, "patient/patientregister.html", "Error registering.Please ensure first name and last name filled in.")
		return
	}

	var ins *sql.Stmt
	ins, err = db.Prepare("INSERT INTO `newdatabase`.`patients`(`firstname`,`lastname`,`ward`) VALUES(?,?,?)")
	if err != nil {
		panic(err)
	}
	defer ins.Close()
	res, err := ins.Exec(firstname, lastname, ward)
	rowsAffec, _ := res.RowsAffected()
	//check for error + if theres more than 1 row affected
	if err != nil || rowsAffec != 1 {
		fmt.Println("Error inserting row: ", err)
		tpl.ExecuteTemplate(w, "patient/patientregister.html", "Error registering patient.Please check first and last name")
		return
	}

	lastInserted, _ := res.LastInsertId()
	rowsAffected, _ := res.RowsAffected()
	fmt.Println("ID of last row inserted: ", lastInserted)
	fmt.Println("No. of rows affected: ", rowsAffected)
	tpl.ExecuteTemplate(w, "patient/patientregister.html", "Patient successfully registered.")
}

func caseAddHandler(w http.ResponseWriter, r *http.Request) {

}

func homePageHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "this is the home page")
}
