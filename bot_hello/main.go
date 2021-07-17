package main

import (
	tb "gopkg.in/tucnak/telebot.v2"
	"log"
	"strconv"
	"time"
)

func main() {
	b, err := tb.NewBot(tb.Settings{
		// You can also set custom API URL.
		// If field is empty it equals to "https://api.telegram.org".
		URL: "https://api.telegram.org",
		// Token here
		Token:  "",
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		log.Fatal(err)
		return
	}

	b.Handle("/hello", func(m *tb.Message) {
		b.Send(m.Sender, m.Payload)
	})

	b.Handle(tb.OnQuery, func (q *tb.Query) {
		log.Println(q.Text)
		urls := []string{
			"https://t0922496-nus-13072021.s3.amazonaws.com/f11.jpg",
			"https://t0922496-nus-13072021.s3.amazonaws.com/f12.jpg",
		}
	
		results := make(tb.Results, len(urls)) // []tb.Result
		for i, url := range urls {
			result := &tb.PhotoResult{
				ResultBase:  tb.ResultBase{},
				URL:         url,
				Width:       0,
				Height:      0,
				Title:       "",
				Description: "",
				Caption:     q.Text,
				ParseMode:   "",
				ThumbURL:    url,
				Cache:       "",
			}

			results[i] = result
			// needed to set a unique string ID for each result
			results[i].SetResultID(strconv.Itoa(i))
		}
	
		err := b.Answer(q, &tb.QueryResponse{
			Results:   results,
			CacheTime: 10, // a minute
		})
	
		if err != nil {
			log.Println(err)
		}
	})

	b.Start()
}

