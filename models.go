package main

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdateAt  time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	ApiKey    string    `json:"api_key"`
}

type Feed struct {
	ID            uuid.UUID    `json:"id"`
	CreatedAt     time.Time    `json:"created_at"`
	UpdatedAt     time.Time    `json:"updated_at"`
	Name          string       `json:"name"`
	Url           string       `json:"url"`
	UserId        uuid.UUID    `json:"userId"`
	LastFetchedAt sql.NullTime `json:"lastFetchedAt"`
}

type FeedFollows struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdateAt  time.Time `json:"updated_at"`
	UserId    uuid.UUID `json:"userId"`
	FeedId    uuid.UUID `json:"feedId"`
}
