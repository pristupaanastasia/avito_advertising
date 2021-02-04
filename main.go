package main

import (
	"database/sql"
	"encoding/json"
	"errors"
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
	Description string    `json:"description,omitempty"`
	Image       []string  `json:"image"`
	Update      time.Time `json:"update"`
}

type Advertisement_js struct {
	Id          int      `json:"id"`
	Price       int      `json:"price"`
	Name        string   `json:"name"`
	Description string   `json:"description,omitempty"`
	Image       []string `json:"image"`
	Update      string   `json:"update"`
}

type Erstruct struct {
	Err  error
	Code int
	Name string
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

func Valid(name string, desc string, image []string, w http.ResponseWriter) bool {
	if len(image) > 3 {
		erro := Erstruct{errors.New("error upload"), 400, "Количество ссылок не должно превышать 3"}
		ErrorHandler(w, erro, 400)
		return true
	}
	if len(desc) > 1000 {
		erro := Erstruct{errors.New("error upload"), 400, "Описание должно быть не больше 1000 символов"}
		ErrorHandler(w, erro, 400)
		return true
	}
	if len(name) > 200 && len(name) == 0 {
		erro := Erstruct{errors.New("error upload"), 400, "Название не должно быть больше 200 символов"}
		ErrorHandler(w, erro, 400)
		return true
	}
	return false
}

func CopyAdv(l *Advertisement, adv *Advertisement_js) {
	(*adv).Name = (*l).Name
	(*adv).Id = (*l).Id
	(*adv).Price = (*l).Price
	(*adv).Description = (*l).Description
	(*adv).Image = (*l).Image
	(*adv).Update = (*l).Update.Format("2006-01-02")
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	page := r.FormValue("page")
	sort := "update"
	switch r.FormValue("sort") {
	case "price_desc":
		sort = "price desc"
	case "price":
		sort = "price"
	case "update_desc":
		sort = "update desc"
	}
	ad, err := database.Query("select * from advertisement order by "+sort+
		" limit 10 offset ($1* 10 -10) ", page) //думаю выборка сразу в бд самая оптимальная
	if err != nil {
		ErrorHandler(w, err, 400)
		return
	}
	list := []Advertisement_js{}
	for ad.Next() {
		l := Advertisement{}
		adv := Advertisement_js{}
		err := ad.Scan(&l.Id, &l.Price, &l.Name, &l.Description, pq.Array(&l.Image), &l.Update)
		if err != nil {
			ErrorHandler(w, err, 400)
			continue
		}
		SpaceTriming(&l)
		l.Description = ""
		buf := l.Image[0] 
		l.Image = []string{buf}
		CopyAdv(&l, &adv)
		list = append(list, adv)
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
	connStr := "host=db port=5432 user=postgres password=1805 dbname=avito_db sslmode=disable"
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
