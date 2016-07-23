package main

// Poll ...
type Poll struct {
	DBPoll
	Answers      []*Answer
	IPSUsed      map[string]bool
	AlreadyVoted bool
	IP           string
}

// Answer ...
type Answer struct {
	Value string
	Total int
}

// DBPoll ... The info in the database table.
type DBPoll struct {
	ID          int    `json:"ID,string" db:"id"`
	Multiselect bool   `db:"multi_select"`
	PerBrowser  bool   `db:"per_browser"`
	PerIP       bool   `db:"per_ip"`
	Question    string `db:"question"`
}

// Save ...
func (p *DBPoll) Save() {
}

const savePollQuery = `INSERT INTO polls (question, )`
