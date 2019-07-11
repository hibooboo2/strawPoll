package main

type inMemoryPollStore struct {
	data map[int]*Poll
}

var _ PollStorer = &inMemoryPollStore{}

// Store stores the poll in the underlying data store.
func (ps *inMemoryPollStore) New(p *Poll) (int, error) {
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
	return poll.Vote(answer, ip)
}
