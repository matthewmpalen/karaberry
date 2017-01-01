package main

import (
	"encoding/json"
	"html/template"
	"net/http"
)

var (
	templateFolder = ""
	templates      = template.Must(template.ParseFiles(templateFolder + "/home.html"))
)

func RenderTemplate(w http.ResponseWriter, name string, context map[string]interface{}) {
	err := templates.ExecuteTemplate(w, name+".html", context)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func JSON(w http.ResponseWriter, i interface{}, statusCode int) {
	js, err := json.Marshal(i)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(js)
}
