package main

import (
	"fmt"
	"log"

	scribble "github.com/nanobox-io/golang-scribble"
)

type scribblePollStorer struct {
	d *scribble.Driver
}

//NewScribbleStorer creates and returns a new PollStorer that is backed using scribble.
func NewScribbleStorer() PollStorer {
	d, err := scribble.New("./straws", nil)
	if err != nil {
		panic(err)
	}
	s := scribblePollStorer{d}
	return &s
}

// Store stores the poll in the underlying data store.
func (ps *scribblePollStorer) New(question string, answers []string, PerIP bool) (int, error) {
	p := &Poll{
		Question: question,
		Answers:  answerStringsToAnswers(answers),
		IPSUsed:  make(map[string]bool),
		PerIP:    PerIP,
	}
	var id int
	err := ps.d.Read("poll", "nextID", &id)
	if err != nil {
		log.Println("No next id found.", err)
	}
	p.ID = id
	err = ps.d.Write("poll", fmt.Sprintf("%d", id), p)
	if err != nil {
		return 0, err
	}
	id++
	err = ps.d.Write("poll", "nextID", &id)
	if err != nil {
		panic(err)
	}
	return p.ID, nil
}

//Get get poll by id
func (ps *scribblePollStorer) Get(id int) (Poll, bool) {
	p := Poll{}
	err := ps.d.Read("poll", fmt.Sprintf("%d", id), &p)
	if err != nil {
		return Poll{}, false
	}
	return p, true
}

//Vote vote in a poll. Return true if your vote was saved and used false otherwise.
func (ps *scribblePollStorer) Vote(id int, answer int, ip string) bool {
	poll := Poll{}
	err := ps.d.Read("poll", fmt.Sprintf("%d", id), &poll)
	if err != nil {
		return false
	}
	if poll.PerIP && poll.IPSUsed[ip] {
		return false
	}
	if len(poll.Answers) > answer {
		poll.Answers[answer].Total++
		if poll.PerIP {
			poll.IPSUsed[ip] = true
		}
		err := ps.d.Write("poll", fmt.Sprintf("%d", poll.ID), poll)
		if err != nil {
			log.Println("Failed to save poll after voting. Vote lost: ", id, answer, ip)
			return false
		}
		return true
	}
	log.Println("Tried to vote with an invalid answer")
	return false
}
