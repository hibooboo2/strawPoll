package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/codehack/go-relax"
	"github.com/codehack/go-relax/filter/logs"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var flagListen = flag.String("l", ":8080", "Used to define which port is listened on.")

func init() {
	flag.Parse()
}

func main() {
	env.config.database.Driver = "postgres"
	env.config.database.DSN = "postgres://wizardofmath:@localhost:5432/strawpoll?sslmode=disable"
	env.db = sqlx.MustConnect(env.config.database.Driver, env.config.database.DSN)

	env.svc = relax.NewService("/",
			CheckArgs: true, DisableAudienceCheck: true})
	env.svc.Use(&logs.Filter{})

	r := mux.NewRouter().StrictSlash(true)
	registerRoutes(r)
}

var env struct {
	db     *sqlx.DB
	config struct {
		database struct {
			DSN    string
			Driver string
		}
	}
	svc *relax.Service
}

func registerRoutes(r *mux.Router) {
	r.HandleFunc("/", makePoll)
	r.HandleFunc("/newpoll/", newPoll)
	r.HandleFunc("/poll/{id:[0-9]+}/", viewPoll).Methods("GET")
	r.HandleFunc("/poll/{id:[0-9]+}/vote/", votePoll).Methods("POST")
	r.HandleFunc("/poll/{id:[0-9]+}/r/", pollResults).Methods("GET")
	http.Handle("/", r)
	log.Println("server started.")
	http.ListenAndServe(*flagListen, nil)
}

func viewPoll(w http.ResponseWriter, req *http.Request) {
	log.Println("Viewing poll.")
	t := template.New("baseTemplate")
	t, _ = t.ParseFiles("view/poll.html")
	pollID, err := strconv.Atoi(mux.Vars(req)["id"])
	if err != nil {
		t := template.New("main")
		t, _ = t.ParseFiles("view/index.html")
		t.ExecuteTemplate(w, "main", nil)
		return
	}
	thePoll := &Poll{
		DBPoll: DBPoll{
			ID: pollID,
		},
	}
	t.ExecuteTemplate(w, "poll", thePoll)
	log.Println("Viewed poll.")

}

func pollResults(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		fmt.Fprintf(w, "Nope dip")
	}

	t := template.New("baseTemplate")
	t, _ = t.ParseFiles("view/results.html")
	pollID, err := strconv.Atoi(mux.Vars(req)["id"])
	if err != nil {
		return
	}
	thePoll := Polls.Get(pollID)
	log.Println(req.PostFormValue("alreadyVoted"))
	thePoll.AlreadyVoted = req.URL.Query().Get("alreadyVoted") == "true"
	if thePoll.AlreadyVoted {
		thePoll.IP = req.RemoteAddr
	} else {
		thePoll.IP = ""
	}
	t.ExecuteTemplate(w, "results", thePoll)
}

func votePoll(w http.ResponseWriter, req *http.Request) {
	log.Println("Voting")
	err := req.ParseForm()
	if err != nil {
		fmt.Fprintf(w, "Nope dip")
	}

	vars := mux.Vars(req)
	pollID, err := strconv.Atoi(vars["id"])
	if err != nil {
		panic("shit")
	}
	if Polls.Get(pollID).PerIP && Polls.Get(pollID).IPSUsed[req.RemoteAddr] {
		log.Println("Duped ip: " + req.RemoteAddr)
		http.Redirect(w, req, fmt.Sprintf("/poll/%d/r/?alreadyVoted=true", pollID), http.StatusSeeOther)
		return
	}
	theAnswers := Polls.Get(pollID).Answers
	chosenAnswer := req.PostFormValue("Answer")
	for _, ans := range theAnswers {
		if ans.Value == chosenAnswer {
			ans.Total = ans.Total + 1
		}
	}
	if Polls.Get(pollID).PerIP {
		Polls.Get(pollID).IPSUsed[req.RemoteAddr] = true
	}
	http.Redirect(w, req, fmt.Sprintf("/poll/%d/r/", pollID), http.StatusSeeOther)
}

func newPoll(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		fmt.Fprintf(w, "Nope dip")
	}
	t := template.New("baseTemplate")
	t, _ = t.ParseFiles("view/poll.html")
	aPoll := &Poll{
		DBPoll: DBPoll{
			Question: req.PostFormValue("question"),
			PerIP:    req.PostFormValue("perIP") == "on",
		},
		Answers: removeBlanks(req.Form["option"]),
		IPSUsed: make(map[string]bool),
	}
	storePoll(aPoll)
	http.Redirect(w, req, fmt.Sprintf("/poll/%d/", aPoll.ID), http.StatusSeeOther)
	log.Println("New poll made")

}

func makePoll(w http.ResponseWriter, req *http.Request) {
	t := template.New("main")
	t, _ = t.ParseFiles("view/index.html")
	t.ExecuteTemplate(w, "main", nil)
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

func storePoll(thePoll *Poll) int {
	thePoll.Save()
	return thePoll.ID
}
