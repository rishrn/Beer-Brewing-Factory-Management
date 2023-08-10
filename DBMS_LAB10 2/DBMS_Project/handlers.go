package main

import (
	"encoding/json"
	"fmt"
	_ "fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

func (app *application) homeHandler(w http.ResponseWriter, r *http.Request) {
	tmp, err := template.ParseFiles("./ui/beerHome.html")
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	err = tmp.Execute(w, nil)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (app *application) registerForm(w http.ResponseWriter, r *http.Request) {
	tmp, err := template.ParseFiles("./ui/buyer_registration.html")
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	err = tmp.Execute(w, nil)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

}

func (app *application) buyerRegistration(w http.ResponseWriter, r *http.Request) {

	// checking for method
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/register-form", http.StatusSeeOther)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	buyer_id, _ := strconv.Atoi(r.PostForm.Get("buyerid"))
	fname := r.PostForm.Get("firstname")
	lname := r.PostForm.Get("lastname")
	phone := r.PostForm.Get("phone")
	city := r.PostForm.Get("city")
	street := r.PostForm.Get("street")
	quantity, _ := strconv.Atoi(r.PostForm.Get("quantity"))
	advance_amount, _ := strconv.Atoi(r.PostForm.Get("adv_amount"))
	id_proof := r.PostForm.Get("id_proof")
	beer_id, _ := strconv.Atoi(r.PostForm.Get("beer_id"))
	beer_name := r.PostForm.Get("beer_name")

	fmt.Println(buyer_id, fname, lname, phone, city, street, quantity, advance_amount, id_proof, beer_id, beer_name)

	insertStatement := `
	  insert into "202001213_db"."Buyer"("buyer_id","first_name","last_name","city","street","quantity","adv_amount","id_proof","beer_id","beer_name") values ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
	`
	_, err = app.db.Exec(insertStatement, buyer_id, fname, lname, city, street, quantity, advance_amount, id_proof, beer_id, beer_name)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Fatal(err)
		return
	}

	insertStatement2 := `
	  insert into "202001213_db"."Buyer_Contact"("buyer_id","contact_no") values($1,$2) 
	`
	_, err = app.db.Exec(insertStatement2, buyer_id, phone)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Fatal(err)
		return
	}

	w.Write([]byte(fmt.Sprintf("%s %s Inserted Successfully", fname, lname)))

}

func (app *application) beerHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	read := `
	  select * from "202001213_db"."Beer"
	`
	rows, err := app.db.Query(read)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	var beers []Beer

	for rows.Next() {
		var beer Beer
		err := rows.Scan(&beer.Beer_id, &beer.Beer_name, &beer.Alcohol, &beer.Cost, &beer.Availability)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		beers = append(beers, beer)
	}

	json.NewEncoder(w).Encode(beers)

	//
	// data, err := json.Marshal(beers)
	// if err != nil {
	// 	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	// 	return
	// }

	// file, err := os.Create("./ui/data.json")
	// file.Write(data)

	// if err != nil {
	// 	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	// 	return
	// }

}

func (app *application) updateFormHandler(w http.ResponseWriter, r *http.Request) {
	tmp, err := template.ParseFiles("./ui/beer_update.html")
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	err = tmp.Execute(w, nil)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

}

func (app *application) beerUpdateHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	beer_id, _ := strconv.Atoi(r.PostForm.Get("beerid"))
	beer_name := r.PostForm.Get("beername")
	alcohol, _ := strconv.Atoi(r.PostForm.Get("alcohol"))
	cost, _ := strconv.Atoi(r.PostForm.Get("cost"))
	availability := r.PostForm.Get("availability")
	var isavailable = true
	if availability == "No" {
		isavailable = false
	}

	fmt.Println(beer_id, beer_name, alcohol, cost, availability)

	update := `
	 update "202001213_db"."Beer"
	 set "alcohol"=$3,
	 "cost"=$4,
     "availability"=$5 
	 where "beer_name"=$2 and "beer_id"=$1
	`
	_, err = app.db.Exec(update, beer_id, beer_name, alcohol, cost, isavailable)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Fatal(err) // 5 8160 false
	}

	w.Write([]byte(fmt.Sprintf("%s Updated Successfully", beer_name)))

}

func (app *application) deleteFormHandler(w http.ResponseWriter, r *http.Request) {
	tmp, err := template.ParseFiles("./ui/beer_delete_form.html")
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Fatal(err)
		return
	}
	err = tmp.Execute(w, nil)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Fatal(err)
		return
	}

}

func (app *application) deleteBeerHandler(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Fatal(err)
		return
	}

	beer_name := r.PostForm.Get("beername")

	delete := `
	 delete from "202001213_db"."Beer"
	 where "beer_name"=$1
     `
	_, err = app.db.Query(delete, beer_name)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Fatal(err)
		return
	}

	w.Write([]byte(fmt.Sprintf("%s Deleted Successfully", beer_name)))

}
