package server

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	"github.com/Nestik55/task0/storage/api"
)

const portNum string = ":8080"

func Handler(w http.ResponseWriter, r *http.Request) {
	uid := r.URL.Query().Get("query")
	if uid != "" {
		answer, ok := api.GetOrder(uid)
		if !ok {
			fmt.Fprintf(w, "Данного элемента нет в кеше и бд.")
			return
		}
		if data, err := json.MarshalIndent(answer, "", " "); err != nil {
			panic(err)
		} else {
			w.Write(data)
		}
		return
	}

	templ, err := template.ParseFiles("server/resource/resource.html")
	if err != nil {
		fmt.Println("cringe", err)
	}
	templ.Execute(w, nil)

}

func Run() {
	http.HandleFunc("/home", Handler)

	fmt.Println("To close connection CTRL+C :-)")

	err := http.ListenAndServe(portNum, nil)
	if err != nil {
		panic(err)
	}
}
