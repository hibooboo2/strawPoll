package main

import (
	"fmt"
	"html/template"
	"net/http"
)

func main() {
	http.HandleFunc("/", makePoll)
	http.HandleFunc("/newpoll/", newPoll)
	http.HandleFunc("/poll/", viewPoll)
	http.HandleFunc("/poll/r/", pollResults)
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

func viewPoll(w http.ResponseWriter, req *http.Request) {

}

func pollResults(w http.ResponseWriter, req *http.Request) {

}

func newPoll(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		fmt.Fprintf(w, "Nope dip")
	}
	t := template.New("baseTemplate")     // Create a template.
	t, _ = t.ParseFiles("view/poll.html") // Parse template file.
	aPoll := Poll{
		Question: req.PostFormValue("question"),
		Answers:  req.Form["option"],
	}
	t.ExecuteTemplate(w, "poll", aPoll) // merge.

}
func makePoll(w http.ResponseWriter, req *http.Request) {
	t := template.New("main")              // Create a template.
	t, _ = t.ParseFiles("view/index.html") // Parse template file.
	t.ExecuteTemplate(w, "main", nil)      // merge.
}

type Poll struct {
	Question    string
	Answers     []string
	Multiselect bool
	PerIp       bool
	PerBrowser  bool
}
