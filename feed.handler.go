package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func (d *dbStruct) handlerCreateFeed(w http.ResponseWriter, r *http.Request, user User) {
	type parameters struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}
	defer r.Body.Close()
	params := &parameters{}
	err := json.NewDecoder(r.Body).Decode(params)

	if err != nil {
		respondWithErrors(w, 400, fmt.Sprintf("error parsing the json: %s", err.Error()))
		return
	}

	query := `insert into feeds(id, created_at, updated_at, name, url, user_id) values($1, $2, $3, $4, $5, $6)`
	uuid := uuid.New()
	cAt := time.Now().UTC()

	sqlRes, err := d.db.Exec(query, &uuid, &cAt, &cAt, &params.Name, &params.URL, &user.ID)

	if err != nil {
		respondWithErrors(w, 400, fmt.Sprintf("error while executing and inserting into feeds %s", err.Error()))
		return
	}

	rowsAffected, err := sqlRes.RowsAffected()

	if err != nil {
		respondWithErrors(w, 400, fmt.Sprintf("error in checking rows affected while inserting %s", err.Error()))
		return
	}

	respondWithJson(w, 201, rowsAffected)
}

// Feeds query with foreign key :: create table feeds(id uuid primary key, created_at timestamp not null, updated_at timestamp not null, name text not null, url text unique not null, user_id UUID references users(id) on delete cascade);
// Feeds query
// query := insert into feeds (id, created_at, updated_at, name, url, user_id) values($1, $2, $3, $4, $5, $6)
