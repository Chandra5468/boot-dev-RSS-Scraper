package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (d *dbStruct) handlerCreateFeedFollow(w http.ResponseWriter, r *http.Request, user User) {

	type parameters struct {
		FeedId uuid.UUID `json:"feed_id"`
	}
	defer r.Body.Close()

	params := &parameters{}

	err := json.NewDecoder(r.Body).Decode(params)
	if err != nil {
		respondWithErrors(w, 400, fmt.Sprintf("error parsing json %s", err.Error()))
		return
	}

	query := `insert into feed_follows (id, created_at, updated_at, user_id, feed_id) values($1, $2, $3, $4, $5)`

	uuidN := uuid.New()
	cAt := time.Now().UTC()
	sR, err := d.db.Exec(query, &uuidN, &cAt, &cAt, &user.ID, &params.FeedId)

	if err != nil {
		respondWithErrors(w, 400, fmt.Sprintf("error while inserting into db %s", err.Error()))
		return
	}

	srf, err := sR.RowsAffected()
	if err != nil {
		respondWithErrors(w, 400, fmt.Sprintf("error while checking no of rows affected %s", err.Error()))
		return
	}

	respondWithJson(w, 201, srf)
}

func (d *dbStruct) handlerGetFeedFollows(w http.ResponseWriter, r *http.Request, user User) {
	query := `select * from feed_follows where user_id = $1`

	rws, err := d.db.Query(query, &user.ID)

	if err != nil {
		respondWithErrors(w, 400, fmt.Sprintf("error while searching for feed dollfows %s", err.Error()))
		return
	}
	allFeeds := []FeedFollows{}
	for rws.Next() {
		feed, err := FeedFlsAllScan(rws)
		if err != nil {
			respondWithErrors(w, 400, fmt.Sprintf("error while scanning all rows. There might be interupption in one row or many rows %s", err.Error()))
			return
		}
		allFeeds = append(allFeeds, *feed)
	}

	respondWithJson(w, 200, allFeeds)
}

func FeedFlsAllScan(rows *sql.Rows) (*FeedFollows, error) {
	feed := &FeedFollows{}
	err := rows.Scan(&feed.ID, &feed.CreatedAt, &feed.UpdateAt, &feed.UserId, &feed.FeedId)
	return feed, err
}

func (d *dbStruct) handlerDeleteFeedFollow(w http.ResponseWriter, r *http.Request, user User) {
	// defer r.Body.Close()
	feedFollowId := chi.URLParam(r, "feedFollowId")
	uid, err := uuid.Parse(feedFollowId)
	if err != nil {
		respondWithErrors(w, 400, fmt.Sprintf("could not parse uuid from string to uuid %s", err.Error()))
		return
	}
	query := `delete from feed_follows where id = $1 and user_id = $2`

	res, err := d.db.Exec(query, &uid, &user.ID)

	if err != nil {
		respondWithErrors(w, 400, fmt.Sprintf("error while deleting from feed_follows table %s", err.Error()))
		return
	}
	ra, err := res.RowsAffected()
	if err != nil {
		respondWithErrors(w, 400, fmt.Sprintf("error finding no of rows affected in feed_follows table %s", err.Error()))
		return
	}
	respondWithJson(w, 204, ra)
}
