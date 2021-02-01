package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
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

func CreateHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	price := r.FormValue("price")
	name := r.FormValue("name")
	desc := r.FormValue("description")
	image := r.FormValue("image")
	fmt.Println("price ", price, "name ", name, "desc ", desc, "image  ", image)
	u, _ := url.Parse(image)
	u.RawQuery, u.Fragment = "", ""
	fmt.Printf("%s\n", u)
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
		panic(400)
		return
	}
	w.Write(json_data)
}

func main() {
	connStr := "user=postgres password= dbname=avito_db sslmode=disable"
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
