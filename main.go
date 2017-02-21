package main

import (
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

var flagListen = flag.String("l", os.Getenv("PORT"), "Used to define which port is listened on.")
var logs = flag.Bool("showlogs", false, "set to show logs.Other wise now logs.")

func init() {
	flag.Parse()

	if !*logs {
		log.SetOutput(ioutil.Discard)
	}
}

func main() {
	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/", makePoll)
	r.HandleFunc("/newpoll/", newPoll)
	r.HandleFunc("/poll/{id:[0-9]+}/", viewPoll).Methods("GET")
	r.HandleFunc("/poll/{id:[0-9]+}/vote/", votePoll).Methods("POST")
	r.HandleFunc("/poll/{id:[0-9]+}/r/", pollResults).Methods("GET")
	http.Handle("/", r)
	log.Println("server started.")
	if !strings.HasPrefix(*flagListen, ":") {
		*flagListen = fmt.Sprintf(":%s", *flagListen)
	}
	err := http.ListenAndServe(*flagListen, nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		fmt.Fprintln(os.Stderr, *flagListen)
	}
}

func viewPoll(w http.ResponseWriter, req *http.Request) {
	log.Println("Viewing poll.")

	t := template.New("baseTemplate")     // Create a template.
	t, _ = t.ParseFiles("view/poll.html") // Parse template file.
	pollId, err := strconv.Atoi(mux.Vars(req)["id"])
	if err != nil {
		t := template.New("main")              // Create a template.
		t, _ = t.ParseFiles("view/index.html") // Parse template file.
		t.ExecuteTemplate(w, "main", nil)      // merge.
		return
	}
	thePoll := polls[pollId]
	t.ExecuteTemplate(w, "poll", thePoll) // merge.
	log.Println("Viewed poll.")

}

func pollResults(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		fmt.Fprintf(w, "Nope dip")
	}

	t := template.New("baseTemplate")        // Create a template.
	t, _ = t.ParseFiles("view/results.html") // Parse template file.
	pollId, err := strconv.Atoi(mux.Vars(req)["id"])
	if err != nil {
		return
	}
	thePoll := polls[pollId]
	log.Println(req.PostFormValue("alreadyVoted"))
	thePoll.AlreadyVoted = req.URL.Query().Get("alreadyVoted") == "true"
	if thePoll.AlreadyVoted {
		thePoll.IP = req.RemoteAddr
	} else {
		thePoll.IP = ""
	}
	t.ExecuteTemplate(w, "results", thePoll) // merge.
}

func votePoll(w http.ResponseWriter, req *http.Request) {
	log.Println("Voting")
	err := req.ParseForm()
	if err != nil {
		fmt.Fprintf(w, "Nope dip")
	}

	vars := mux.Vars(req)
	pollId, err := strconv.Atoi(vars["id"])
	if err != nil {
		panic("shit")
	}
	if polls[pollId].PerIp && polls[pollId].IPSUsed[req.RemoteAddr] {
		log.Println("Duped ip: " + req.RemoteAddr)
		http.Redirect(w, req, fmt.Sprintf("/poll/%d/r/?alreadyVoted=true", pollId), http.StatusSeeOther)
		return
	}
	theAnswers := polls[pollId].Answers
	chosenAnswer := req.PostFormValue("Answer")
	for _, ans := range theAnswers {
		if ans.Value == chosenAnswer {
			ans.Total = ans.Total + 1
		}
	}
	if polls[pollId].PerIp {
		polls[pollId].IPSUsed[req.RemoteAddr] = true
	}
	http.Redirect(w, req, fmt.Sprintf("/poll/%d/r/", pollId), http.StatusSeeOther)
}

func newPoll(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		fmt.Fprintf(w, "Nope dip")
	}
	t := template.New("baseTemplate")     // Create a template.
	t, _ = t.ParseFiles("view/poll.html") // Parse template file.
	aPoll := &Poll{
		Question: req.PostFormValue("question"),
		Answers:  removeBlanks(req.Form["option"]),
		IPSUsed:  make(map[string]bool),
		PerIp:    req.PostFormValue("perIP") == "on",
	}
	storePoll(aPoll)
	http.Redirect(w, req, fmt.Sprintf("/poll/%d/", aPoll.Id), http.StatusSeeOther)
	log.Println("New poll made")

}

func makePoll(w http.ResponseWriter, req *http.Request) {
	t := template.New("main")              // Create a template.
	t, _ = t.ParseFiles("view/index.html") // Parse template file.
	t.ExecuteTemplate(w, "main", nil)      // merge.
}

type Poll struct {
	Question     string
	Answers      []*Answer
	Multiselect  bool
	PerIp        bool
	PerBrowser   bool
	Id           int
	IPSUsed      map[string]bool
	AlreadyVoted bool
	IP           string
}
type Answer struct {
	Value string
	Total int
}

func removeBlanks(theStrings []string) []*Answer {
	nonblank := []*Answer{}
	for _, aString := range theStrings {
		if len(aString) > 0 {
			nonblank = append(nonblank, &Answer{Value: aString, Total: 0})
		}
	}
	return nonblank
}

var polls map[int]*Poll = make(map[int]*Poll)

func storePoll(thePoll *Poll) int {
	thePoll.Id = len(polls)
	polls[thePoll.Id] = thePoll
	return thePoll.Id
}
