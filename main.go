package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"net/url"

	//"strings"
	"time"

	//"log"
	"net/http"
)

type Advertisement struct {
	Id          int
	Price       int
	Name        string
	Description string
	Image       []string
	Update      time.Time
}

var database *sql.DB

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
	body, _ := url.Parse(r.RequestURI)
	v, _ := url.ParseQuery(body.RawQuery)
	fmt.Println(price)
	update := time.Now().Format("2006-01-02")
	err := database.QueryRow("insert into advertisement (price, name, description, image, update ) values ($1, $2, $3, $4, $5) returning id",
		price, name, desc, pq.Array(v["image"]), update).Scan(&id)
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
	id := r.FormValue("id")
	ad, err := database.Query("select * from advertisement where id = $1", id)
	adv := Advertisement{}
	err = ad.Scan(&adv.Id, &adv.Price, &adv.Name, &adv.Description, &adv.Image, &adv.Update)
	if err != nil {
		panic(err)
	}
	json_data, errno := json.Marshal(adv)
	if errno != nil {
		panic(400)
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

		err := ad.Scan(&l.Id, &l.Price, &l.Name, &l.Description, &l.Image, &l.Update)

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
	http.HandleFunc("/list", FindHandler)
	http.HandleFunc("/create", CreateHandler)
	http.ListenAndServe(":9000", nil)
}
