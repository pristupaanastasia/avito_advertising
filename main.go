package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/lib/pq"

	"strings"
	"time"

	//"log"
	"net/http"
)

type Advertisement struct {
	Id          int       `json:"id"`
	Price       int       `json:"price"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Image       []string  `json:"image"`
	Update      time.Time `json:"update"`
}

var database *sql.DB

func SpaceTriming(adv *Advertisement) {
	(*adv).Name = strings.TrimSpace((*adv).Name)
	(*adv).Description = strings.TrimSpace((*adv).Description)
	for i, _ := range (*adv).Image {
		(*adv).Image[i] = strings.TrimSpace((*adv).Image[i])
	}
}

func ErrorHandler(w http.ResponseWriter, err interface{}, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(err)
}

func CreateHandler(w http.ResponseWriter, r *http.Request) {
	var id int
	r.ParseForm()
	price := r.FormValue("price")
	name := r.FormValue("name")
	desc := r.FormValue("description")
	image := r.Form["image"]
	update := time.Now().Format("2006-01-02")
	err := database.QueryRow("insert into advertisement (price, name, description, image, update ) values ($1, $2, $3, $4, $5) returning id",
		price, name, desc, pq.Array(image), update).Scan(&id)
	if err != nil {
		ErrorHandler(w, err, 400)
		//panic(400)
		return
	}
	json_return := map[string]int{"id": id}
	json_data, errno := json.Marshal(json_return)
	if errno != nil {
		ErrorHandler(w, errno, 400)
		return
	}
	w.Write(json_data)
}
func FindHandler(w http.ResponseWriter, r *http.Request) {
	//r.ParseForm()
	id := r.FormValue("id")
	//fields := r.FormValue("fields")
	ad := database.QueryRow("select * from advertisement where id = $1", id)
	adv := Advertisement{}
	err := ad.Scan(&adv.Id, &adv.Price, &adv.Name, &adv.Description, pq.Array(&adv.Image), &adv.Update)
	SpaceTriming(&adv)
	if err != nil {
		ErrorHandler(w, err, 400)
	}
	fmt.Println("|", r.Form)
	//fmt.Println("-",id)
	json_data, errno := json.Marshal(adv)
	if errno != nil {
		ErrorHandler(w, errno, 400)
		return
	}
	w.Write(json_data)
}
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	page := r.FormValue("page")
	sort := r.FormValue("sort")
	ad, err := database.Query("select * from advertisement order by $2 limit 10 offset ($1* 10 -10) ", page, sort) //думаю такая выборка самая оптимальная
	if err != nil {
		panic(400)
		return
	}
	list := []Advertisement{}
	for ad.Next() {
		l := Advertisement{}

		err := ad.Scan(&l.Id, &l.Price, &l.Name, &l.Description, pq.Array(&l.Image), &l.Update)

		if err != nil {
			panic(400)
			continue
		}
		list = append(list, l)
	}
	defer ad.Close()
	json_data, errno := json.Marshal(list)
	if errno != nil {
		ErrorHandler(w, errno, 400)
		return
	}
	w.Write(json_data)
}

func main() {
	connStr := "user=postgres dbname=avito_db sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	database = db
	defer db.Close()
	http.HandleFunc("/ad", IndexHandler)
	http.HandleFunc("/find", FindHandler)
	http.HandleFunc("/create", CreateHandler)
	http.ListenAndServe(":9000", nil)
}
