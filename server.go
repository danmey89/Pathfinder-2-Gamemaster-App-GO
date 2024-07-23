package main

import (
	"html/template"
	"net/http"
)

var templates = template.Must(template.ParseFiles("templates/index.html"))

var fs = http.FileServer(http.Dir("./static"))

func indexHandler(w http.ResponseWriter, r *http.Request) {

	p := Page{
		Proficiencies: proficiencies,
		Abilities:     abilities,
		Data:          charMap,
	}

	err := templates.ExecuteTemplate(w, "index.html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
