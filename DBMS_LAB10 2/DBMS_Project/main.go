package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type application struct {
	db *sql.DB
}

type Beer struct {
	Beer_id      int     `json:"beer_id"`
	Beer_name    string  `json:"beer_name"`
	Alcohol      float32 `json:"alcohol"`
	Cost         float32 `json:"cost"`
	Availability bool    `json:"availability"`
}

func setup() (*sql.DB, error) {

	const (
		host     = "localhost"
		port     = 8080
		user     = "postgres"
		password = "maher"
		dbname   = "postgres"
	)

	pgStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err := sql.Open("postgres", pgStr)

	if err != nil {
		return nil, err
	}

	err = db.Ping()

	if err != nil {
		return nil, err
	}

	fmt.Println("Connection Successfull")

	return db, nil
}

func main() {
	
	db, err := setup()

	if err != nil {
		log.Fatal(err)
	}

	app := &application{
		db: db,
	}

	r := mux.NewRouter()

	r.HandleFunc("/home", app.homeHandler).Methods("GET")   //get
	r.HandleFunc("/show-beer", app.beerHandler).Methods("GET")     //get
	r.HandleFunc("/register-form", app.registerForm).Methods("GET")       //get
	r.HandleFunc("/buyer-register", app.buyerRegistration).Methods("POST")        // create
	r.HandleFunc("/beer-update-form", app.updateFormHandler).Methods("GET")     //get    
	r.HandleFunc("/beer-update", app.beerUpdateHandler).Methods("POST")    // update
	r.HandleFunc("/ber-delete-form",app.deleteFormHandler).Methods("GET"); // get
	r.HandleFunc("/beer-delete",app.deleteBeerHandler).Methods("POST");   // delete

	log.Fatal(http.ListenAndServe(":4000", r))

}
