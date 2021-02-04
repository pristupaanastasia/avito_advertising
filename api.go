package main

import (
	"encoding/json"
	"github.com/lib/pq"
	"net/http"
	"time"
)

func CreateHandler(w http.ResponseWriter, r *http.Request) {
	var id int
	r.ParseForm()
	price := r.FormValue("price")
	name := r.FormValue("name")
	desc := r.FormValue("description")
	image := r.Form["image"]
	if Valid(name, desc, image, w) {
		return
	}
	update := time.Now().Format("2006-01-02")
	err := database.QueryRow("insert into advertisement (price, name, description, image, update ) values ($1, $2, $3, $4, $5) returning id",
		price, name, desc, pq.Array(image), update).Scan(&id)
	if err != nil {
		ErrorHandler(w, err, 400)
		return
	}
	json_return := map[string]int{"id": id, "status": 201} //если возникает ошибка, то не возвращает id
	json_data, errno := json.Marshal(json_return)
	if errno != nil {
		ErrorHandler(w, errno, 418)
		return
	}
	w.Write(json_data)
}

func FindHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	ad := database.QueryRow("select * from advertisement where id = $1", id)
	adv := Advertisement{}
	adv1 := Advertisement_js{}
	err := ad.Scan(&adv.Id, &adv.Price, &adv.Name, &adv.Description, pq.Array(&adv.Image), &adv.Update)
	if err != nil {
		ErrorHandler(w, err, 400)
		return
	}
	SpaceTriming(&adv)
	if _, ok := r.Form["fields"]; !ok {
		adv.Description = ""
		buf := adv.Image[0] //переменная buf введена, что бы не было коллизии данных в памяти
		adv.Image = []string{buf}
	}
	CopyAdv(&adv, &adv1)
	json_data, errno := json.Marshal(adv1)
	if errno != nil {
		ErrorHandler(w, errno, 400)
		return
	}
	w.Write(json_data)
}
