package main

//PollStorer Interface used to allow crud of polls.
type PollStorer interface {
	New(p *Poll) (int, error)
	Get(id int) (Poll, bool)
	Vote(id int, answer int, ip string) bool
}
