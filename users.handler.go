package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func (d *dbStruct) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	type parameters struct {
		Name string `json:"name"`
	}
	// params := &parameters{}
	var params *parameters
	err := json.NewDecoder(r.Body).Decode(&params)

	if err != nil {
		respondWithErrors(w, http.StatusBadRequest, fmt.Sprintf("error parsing json: %s", err.Error()))
		return
	}

	query := `insert into users(id, created_at, updated_at, name) values ($1, $2, $3, $4)`
	// As we are using
	// alter table users add column api_key varchar(64) unique not null default (encode(sha256(random()::text::bytea),'hex'))
	// So no need to add $5 in query to insert random hex values.
	// Learn difference among sha256, md5, digest and hex, ASCII.
	uuidN := uuid.New()
	cAt := time.Now().UTC()
	res, err := d.db.Exec(query, uuidN, cAt, cAt, &params.Name)
	if err != nil {
		respondWithErrors(w, http.StatusInternalServerError, fmt.Sprintf("Unable to insert into users table %s", err.Error()))
		return
	}

	rf, err := res.RowsAffected()
	if err != nil {
		respondWithErrors(w, http.StatusInternalServerError, fmt.Sprintf("error from psql while insert into users table %s", err.Error()))
		return
	}

	respondWithJson(w, http.StatusCreated, fmt.Sprintf("no of docs inserted : %d", rf))
}

func (d *dbStruct) handlerGetUser(w http.ResponseWriter, r *http.Request) {
	apiKey, err := GetAPIKey(r.Header)

	if err != nil {
		respondWithErrors(w, 403, err.Error())
		return
	}

	query := `select * from users where api_key = $1`
	var user User

	err = d.db.QueryRow(query, apiKey).Scan(&user.ID, &user.CreatedAt, &user.UpdateAt, &user.Name, &user.ApiKey)
	if err != nil {
		log.Println("This is error from scan of records ", err.Error())
		respondWithErrors(w, 404, "there is no record found in postgresql")
		return
	}
	respondWithJson(w, http.StatusOK, &user)
}
