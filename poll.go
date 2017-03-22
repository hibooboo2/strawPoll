package main

//Poll used to represent a poll.
type Poll struct {
	Question     string
	Answers      []Answer
	Multiselect  bool
	PerIP        bool
	PerBrowser   bool
	Id           int
	IPSUsed      map[string]bool
	AlreadyVoted bool
	IP           string
}

//Answer used to represent an answer and the number of votes it has.
type Answer struct {
	Value string
	Total int
}
