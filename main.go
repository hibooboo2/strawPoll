package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", makePoll)
	r.HandleFunc("/newpoll/", newPoll)
	r.HandleFunc("/poll/{id:[0-9]+}", viewPoll).Methods("GET")
	r.HandleFunc("/poll/{id:[0-9]+}/r/", pollResults).Methods("GET")
	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}

func viewPoll(w http.ResponseWriter, req *http.Request) {
	t := template.New("baseTemplate")     // Create a template.
	t, _ = t.ParseFiles("view/poll.html") // Parse template file.
	pollId, err := strconv.Atoi(mux.Vars(req)["id"])
	if err != nil {
		return
	}
	thePoll := polls[pollId]
	t.ExecuteTemplate(w, "poll", thePoll) // merge.
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
		Answers:  removeBlanks(req.Form["option"]),
	}
	storePoll(aPoll)
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

func removeBlanks(theStrings []string) []string {
	nonblank := []string{}
	for _, aString := range theStrings {
		if len(aString) > 0 {
			nonblank = append(nonblank, aString)
		}
	}
	return nonblank
}

var polls []Poll = []Poll{}

func storePoll(thePoll Poll) int {
	polls = append(polls, thePoll)
	return len(polls) - 1
}
