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
	"time"

	"github.com/gorilla/mux"
	"github.com/skratchdot/open-golang/open"
)

var (
	flagListen = flag.String("l", "8080", "Used to define which port is listened on.")
	dbToUse    = flag.String("db", "scribble", "Used to define which poll storage is used")
	logs       = flag.Bool("showlogs", false, "set to show logs.Other wise now logs.")
	autoopen   = flag.Bool("dontopen", false, "Don't Autoopen polls")
	t          = template.Must(getTemplates())
	polls      PollStorer
)

func init() {
	flag.Parse()
	if !*logs {
		log.SetOutput(ioutil.Discard)
	}
	switch *dbToUse {
	case "inmemory":
		polls = &inMemoryPollStore{make(map[int]*Poll)}
	case "scribble":
		polls = NewScribbleStorer()
	default:
		log.Fatalln("Invalid db type: ", *dbToUse)
	}
	go func() {
		for {
			time.Sleep(time.Second)
			temp, err := getTemplates()
			if err != nil {
				log.Println(err)
				continue
			}
			t = temp
		}
	}()
}

func getTemplates() (*template.Template, error) {
	files, err := AssetDir("view")
	if err != nil {
		return nil, err
	}
	temp := template.New("base")
	for _, f := range files {
		data, err := Asset("view/" + f)
		if err != nil {
			return nil, err
		}
		temp, err = temp.Parse(string(data))
		if err != nil {
			return nil, err
		}
	}
	return temp, nil
}

func main() {
	log.Println("Starting with PID:", os.Getpid(), "Current version:", version, os.Args)
	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/", makePoll)
	r.HandleFunc("/newpoll/", newPoll)
	r.HandleFunc("/poll/{id:[0-9]+}/", viewPoll).Methods("GET")
	r.HandleFunc("/poll/{id:[0-9]+}/vote/", votePoll).Methods("POST")
	r.HandleFunc("/poll/{id:[0-9]+}/r/", pollResults).Methods("GET")
	if !strings.HasPrefix(*flagListen, ":") {
		*flagListen = fmt.Sprintf(":%s", *flagListen)
	}
	if !*autoopen {
		go open.Start("http://localhost" + *flagListen)
	}
	log.Println("Starting on port: ", *flagListen)
	err := http.ListenAndServe(*flagListen, r)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		fmt.Fprintln(os.Stderr, *flagListen)
	}
}

func viewPoll(w http.ResponseWriter, req *http.Request) {
	log.Println("Viewing poll.")

	pollID, err := strconv.Atoi(mux.Vars(req)["id"])
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	thePoll, ok := polls.Get(pollID)
	if !ok {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	err = t.ExecuteTemplate(w, "poll", thePoll) // merge.
	if err != nil {
		log.Println(err)
	}
}

func pollResults(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	pollID, err := strconv.Atoi(mux.Vars(req)["id"])
	if err != nil {
		return
	}
	thePoll, _ := polls.Get(pollID)
	log.Println(req.PostFormValue("alreadyVoted"))
	thePoll.AlreadyVoted = req.URL.Query().Get("alreadyVoted") == "true"
	if thePoll.AlreadyVoted {
		thePoll.IP = req.RemoteAddr
	} else {
		thePoll.IP = ""
	}
	err = t.ExecuteTemplate(w, "results", thePoll) // merge.
	if err != nil {
		log.Println(err)
		fmt.Fprintln(w, err.Error())
	}
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
	n, err := strconv.Atoi(req.PostFormValue("Answer"))
	if err != nil {
		http.Error(w, "Invalid answer", http.StatusBadRequest)
		return
	}
	voted := polls.Vote(pollID, n, strings.Split(req.RemoteAddr, ":")[0])
	if !voted {
		log.Println("Unable to vote. Probably already voted: " + req.RemoteAddr)
		http.Redirect(w, req, fmt.Sprintf("/poll/%d/r/?alreadyVoted=true", pollID), http.StatusSeeOther)
		return
	}
	http.Redirect(w, req, fmt.Sprintf("/poll/%d/r/", pollID), http.StatusSeeOther)
}

func newPoll(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		fmt.Fprintf(w, "Nope dip")
	}
	id, err := polls.New(req.PostFormValue("question"), req.Form["option"], req.PostFormValue("PerIP") == "on")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, req, fmt.Sprintf("/poll/%d/", id), http.StatusSeeOther)
	log.Println("New poll made")
}

func makePoll(w http.ResponseWriter, req *http.Request) {
	t.ExecuteTemplate(w, "main", nil) // merge.
}

func answerStringsToAnswers(theStrings []string) []Answer {
	nonblank := []Answer{}
	for _, aString := range theStrings {
		if len(aString) > 0 {
			nonblank = append(nonblank, Answer{Value: aString, Total: 0})
		}
	}
	return nonblank
}
