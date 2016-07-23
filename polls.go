package main

// PollsCol resource used to manage the pollsCol.
type PollsCol struct{}

// Polls ...
var Polls = PollsCol{}

// Get ...
func (p *PollsCol) Get(id int) Poll {
	return Poll{}
}

// Index ...
func (p *PollsCol) Index() []Poll {
	return []Poll{}
}
