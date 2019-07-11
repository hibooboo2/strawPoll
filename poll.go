package main

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"log"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

//Poll used to represent a poll.
type Poll struct {
	ID           int `json:"ID" gorm:"primary_key" `
	Question     string
	Answers      []Answer `gorm:"auto_preload"`
	Multiselect  bool
	PerIP        bool
	PerBrowser   bool
	IPSUsed      *IPList `gorm:"type:varchar"`
	AlreadyVoted bool
	IP           string
}

type IPList map[string]bool

var _ driver.Valuer = &IPList{}
var _ sql.Scanner = &IPList{}

func (i *IPList) Scan(src interface{}) error {
	data, ok := src.(string)
	if !ok {
		return errors.New(fmt.Sprintf("invalid data type: %T", src))
	}
	ips := strings.Split(data, "\n")
	*i = IPList{}
	for _, ip := range ips {
		if ip == "" {
			continue
		}
		(*i)[ip] = true
	}
	return nil
}

func (i *IPList) Value() (driver.Value, error) {
	var data bytes.Buffer
	for k, ok := range *i {
		if ok {
			data.WriteString(k)
			data.WriteByte('\n')
		}
	}
	return data.String(), nil
}

func NewPoll(question string, answers []string, PerIP bool) *Poll {
	p := &Poll{
		Question: question,
		Answers:  answerStringsToAnswers(answers),
		PerIP:    PerIP,
		IPSUsed:  &IPList{},
	}
	(*p.IPSUsed)["dick"] = true
	return p
}

func (p *Poll) Used(ip string) bool {
	return (*p.IPSUsed)[ip]
}

func (p *Poll) Vote(answer int, ip string) bool {
	if p.PerIP && p.Used(ip) {
		return false
	}
	if !(len(p.Answers) > answer) {
		log.Println("Tried to vote with an invalid answer")
		return false
	}

	p.Answers[answer].Total++
	if p.PerIP {
		(*p.IPSUsed)[ip] = true
	}
	return true
}

//Answer used to represent an answer and the number of votes it has.
type Answer struct {
	gorm.Model `json:"-"`
	PollID     int `gorm:"index" json:"-"`
	Value      string
	Total      int
}
