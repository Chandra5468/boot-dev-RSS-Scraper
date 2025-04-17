package main

import (
	"log"
	"sync"
	"time"
)

/*
This scrapper function is going to be running in the background as our server runs.
*/

func startScrapping(
	d *dbStruct,
	concurrency int,
	timeBetweenRequest time.Duration) {
	log.Printf("Scrapping on %v goroutines every %s duration", concurrency, timeBetweenRequest)

	ticker := time.NewTicker(timeBetweenRequest)

	// NewTicker returns a new [Ticker] containing a channel that will send the current time on the channel after each tick.
	// The period of the ticks is specified by the duration argument.

	for ; ; <-ticker.C { // If timeBetweenRequest is 1 minute this for loop will execute for every 1 minute
		// Get next(Upcoming or new) feeds to fetch
		query := `select * from feeds order by last_fetched_at asc nulls first limit $1`

		sRows, err := d.db.Query(query, int32(concurrency))

		if err != nil {
			log.Println("error while getting rows ", err)
			continue
		}
		wg := &sync.WaitGroup{}
		for sRows.Next() {
			wg.Add(1)
			feed := &Feed{}
			err = sRows.Scan(&feed.ID, &feed.CreatedAt, &feed.UpdatedAt, &feed.Name, &feed.Url, &feed.UserId, &feed.LastFetchedAt)
			if err != nil {
				log.Println("Error at scanning rows-------", err)
				wg.Done()
				continue
			} else {
				go scrapeFeed(d, wg, feed)
			}
		}
		wg.Wait()
	}
	// Below for loop will not work as do while request. Above for acts like do while
	/*
		for range ticker.C{}
	*/

}

func scrapeFeed(d *dbStruct, wg *sync.WaitGroup, feed *Feed) {
	defer wg.Done()

	// Mark feed as fetched
	query := `update feeds set last_fetched_at = $1 where id = $2`
	// query := `update feeds set last_fetched_at = $1 where id = $2 returning *`
	updatedAt := time.Now().UTC()
	updatedFeed := &Feed{}
	err := d.db.QueryRow(query, &updatedAt, &feed.ID).Scan(&updatedFeed.ID, &updatedFeed.CreatedAt, &updatedFeed.UpdatedAt, &updatedFeed.Name, &updatedFeed.Url, &updatedFeed.UserId, &updatedFeed.LastFetchedAt)

	// err = rws.Scan(&updatedFeed.ID, &updatedFeed.CreatedAt, &updatedFeed.UpdatedAt, &updatedFeed.Name, &updatedFeed.Url, &updatedFeed.UserId, &updatedFeed.LastFetchedAt)

	if err != nil {
		log.Printf("Unable to get updated record of feed post reading %s", err.Error())
		return
	} else {
		log.Println("This is updated document ", updatedFeed)
	}
	// res, err := d.db.Exec(query, &updatedAt, &feed.ID)

	// if err != nil {
	// 	log.Println("Error while updating scrape feed", err.Error())
	// 	return
	// }
	// aff, err := res.RowsAffected()
	// if err != nil {
	// 	log.Println("Number of rows affected error...", err.Error())
	// 	return
	// }
	// log.Printf("Number of rows affected is %d", aff)
}
