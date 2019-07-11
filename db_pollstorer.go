package main

import (
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
)

type dbPollStore struct {
	db *gorm.DB
}

var _ PollStorer = &dbPollStore{}

func NewDbPollStore() (*dbPollStore, error) {
	// db, err := gorm.Open("sqlite3", "straws.db")
	db, err := gorm.Open("postgres", "host=10.14.12.11 port=5432 user=postgres password=docker sslmode=disable")
	if err != nil {
		return nil, errors.Wrap(err, "failed to create database")
	}
	err = db.AutoMigrate(&Poll{}, &Answer{}).Error
	if err != nil {
		return nil, errors.Wrap(err, "failed to create tables for polls and answers")
	}
	return &dbPollStore{
		db,
	}, nil
}

func (db *dbPollStore) New(p *Poll) (int, error) {
	err := db.db.Create(p).Error
	if err != nil {
		return 0, errors.Wrap(err, "failed to save new poll")
	}
	log.Println(p.ID, p)
	return p.ID, nil
}
func (db *dbPollStore) Get(id int) (Poll, bool) {
	var p Poll
	if db.db.Preload("Answers").First(&p, id).RecordNotFound() {
		log.Println("poll not found", id)
		return p, false
	}
	return p, true
}
func (db *dbPollStore) Vote(id int, answer int, ip string) bool {
	var p Poll
	err := db.db.Preload("Answers").First(&p, id).Error
	if err != nil {
		return false
	}

	if !p.Vote(answer, ip) {
		return false
	}
	err = db.db.Save(&p).Error
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}
