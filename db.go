package main

import (
	"strconv"
	"encoding/json"
	"github.com/HouzuoGuo/tiedot/db"
	"github.com/fatih/structs"
)

// appDB represents a tiedot DB.
type appDB struct {
	*db.DB
}

// Open opens a tiedot database reads into the method receiver.
func (appdb *appDB) Open(path string) error {
	var err error
	appdb.DB, err = db.OpenDB(path)
	if err != nil {
		return err
	}

	return nil
}	

// Initialize creates the neccessary collections and searchable
// indexes for the database.
func (appdb *appDB) Initialize(path string) error {
	err := appdb.Create("Tracks")
	if err != nil {
		return err
	}

	tracks := appdb.Use("Tracks")
	if err := tracks.Index([]string{"Title"}); err != nil {
		return err
	}
	if err := tracks.Index([]string{"Album"}); err != nil {
		return err
	}
	if err := tracks.Index([]string{"Artist"}); err != nil {
		return err
	}
	
	return nil
}

func (appdb *appDB) AddTrack(track *Track) error {
	tracks := appdb.Use("Tracks")
	record := structs.Map(track)
	_, err := tracks.Insert(record)
	if err != nil {
		return err
	}

	return nil
}
		
func (appdb *appDB) Query(index, query string) (string, error) {
	tracks := appdb.Use("Tracks")

	var q interface{}
	json.Unmarshal([]byte(`[{"in": ["`+index+`"], "eq": "`+query+`"}]`), &q)
	
	queryResult := make(map[int]struct{})
	if err :=  db.EvalQuery(q, tracks, &queryResult); err != nil {
		return "", err
	}

	resultDocs := make(map[string]interface{})
	for docID := range queryResult {
		doc, _ := tracks.Read(docID)
		if doc != nil {
			resultDocs[strconv.Itoa(docID)] = doc
		}
	}
	
	results, err := json.Marshal(resultDocs)
	if err != nil {
		return "", err
	}
		
	return string(results), nil
}

