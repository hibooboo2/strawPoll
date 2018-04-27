package main

import "log"

type inMemoryPollStore struct {
	data map[int]*Poll
}

var _ PollStorer = &inMemoryPollStore{}

// Store stores the poll in the underlying data store.
func (ps *inMemoryPollStore) New(question string, answers []string, PerIP bool) (int, error) {
	p := &Poll{
		Question: question,
		Answers:  answerStringsToAnswers(answers),
		IPSUsed:  make(map[string]bool),
		PerIP:    PerIP,
	}
	p.ID = len(ps.data)
	ps.data[p.ID] = p
	return p.ID, nil
}

//Get get poll by id
func (ps *inMemoryPollStore) Get(id int) (Poll, bool) {
	p, ok := ps.data[id]
	if !ok {
		return Poll{}, false
	}
	return *p, true
}

//Vote vote in a poll. Return true if your vote was saved and used false otherwise.
func (ps *inMemoryPollStore) Vote(id int, answer int, ip string) bool {
	poll, ok := ps.data[id]
	if !ok {
		return false
	}
	if poll.PerIP && poll.IPSUsed[ip] {
		return false
	}
	theAnswers := poll.Answers
	if len(theAnswers) > answer {
		theAnswers[answer].Total++
		if poll.PerIP {
			poll.IPSUsed[ip] = true
		}
		return true
	}
	log.Println("Tried to vote with an invalid answer")
	return false
}
