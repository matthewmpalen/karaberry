package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
)

var (
	templates *template.Template
)

func init() {
	names := []string{"home", "idle"}
	for i := range names {
		names[i] = fmt.Sprintf("%s/%s.html", Config.TemplateFolder, names[i])
	}
	templates = template.Must(template.ParseFiles(names...))
}

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
