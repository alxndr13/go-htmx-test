package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/exp/slices"
)

type Films []Film

type Film struct {
	Id       string
	Title    string
	Director string
}

func main() {
	fmt.Println("Starting Server on http://localhost:8000")

	// films := Films{{
	// 	Id:       uuid.NewString(),
	// 	Title:    "Feuer und Flamme",
	// 	Director: "WDR",
	// }}
	films := Films{}

	defaultHandler := func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("index.html"))
		tmpl.Execute(w, films)
	}

	addHandler := func(w http.ResponseWriter, r *http.Request) {
		title := r.PostFormValue("title")
		director := r.PostFormValue("director")
		tmpl := template.Must(template.ParseFiles("index.html"))
		newFilm := Film{Id: uuid.NewString(), Title: title, Director: director}
		films = append(films, newFilm)
		err := tmpl.ExecuteTemplate(w, "film-list-element", films)
		if err != nil {
			fmt.Println("error", err)
		}

	}

	deleteHandler := func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("index.html"))
		uri, err := url.Parse(r.RequestURI)
		if err != nil {
			tmpl.ExecuteTemplate(w, "film-list-element", films)
		}
		id := strings.Split(uri.Path, "/")[2]
		fmt.Println(id)

		for k, v := range films {
			if v.Id == id {
				films = slices.Delete(films, k, k+1)
				break
			}
		}
		fmt.Println(films)
		tmpl.ExecuteTemplate(w, "film-list-element", films)
	}

	// define handlers
	http.HandleFunc("/", defaultHandler)
	http.HandleFunc("/add-film/", addHandler)
	http.HandleFunc("/delete-film/", deleteHandler)

	log.Fatal(http.ListenAndServe(":8000", nil))

}
