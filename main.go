package main

import (
	"html/template"
	"net/http"
)

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	t := template.New("baseTemplate")    // Create a template.
	t, _ = t.ParseFiles("question.html") // Parse template file.
	aPoll := Poll{
		Question: "What are you doing for lunch?",
		Answers:  []string{"x", "y"},
	}
	t.ExecuteTemplate(w, "main", aPoll) // merge.
}

type Poll struct {
	Question    string
	Answers     []string
	Multiselect bool
	PerIp       bool
	PerBrowser  bool
}
