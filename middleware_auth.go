package main

import (
	"log"
	"net/http"
)

type authHandler func(http.ResponseWriter, *http.Request, User)

func (d *dbStruct) middlewareAuth(handler authHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		handler(w, r, user)
	}
}
