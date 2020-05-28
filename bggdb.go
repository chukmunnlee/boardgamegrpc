package main

import (
	"fmt"
	"log"

	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

const (
	FIND_BY_NAME = "select gid, name, ranking, url from game where name like ? order by ranking limit ? offset ?"
)

type BoardgameDB struct {
	Username string
	Password string

	db *sql.DB
}

type Boardgame struct {
	Gid     int32
	Name    string
	Ranking uint32
	Url     string
}

type QueryResult struct {
	Error  error
	Result interface{}
}

// Connects and opens the database
func (bg *BoardgameDB) Open() error {
	dsn := fmt.Sprintf("%s:%s@/bgg", bg.Username, bg.Password)

	db, err := sql.Open("mysql", dsn)
	if nil != err {
		return err
	}
	bg.db = db

	return nil
}

// Closes the database
func (bg *BoardgameDB) Close() error {
	return bg.db.Close()
}

// Perform the query
func (bg *BoardgameDB) FindBoardgamesByName(query string, limit uint32, offset uint32) <-chan *QueryResult {

	c := make(chan *QueryResult)

	go func() {
		pattern := fmt.Sprintf("%%%s%%", query)

		rows, err := bg.db.Query(FIND_BY_NAME, pattern, limit, offset)
		if nil != err {
			c <- &QueryResult{Error: err}
			return
		}
		defer rows.Close()
		defer close(c)

		for rows.Next() {
			rec := &Boardgame{}
			err := rows.Scan(&rec.Gid, &rec.Name, &rec.Ranking, &rec.Url)
			if nil != err {
				c <- &QueryResult{Error: err}
				return
			}
			c <- &QueryResult{Result: *rec, Error: nil}
		}
	}()

	return c
}

// Utility functions
func scanBoardgameRecord(r *sql.Rows) (*Boardgame, error) {
	rec := &Boardgame{}
	err := r.Scan(&rec.Gid, &rec.Name, &rec.Ranking, &rec.Url)
	return rec, err
}

func checkError(err error) {
	if nil != err {
		log.Fatalf("Error: %v\n", err)
	}
}
