package main

import (
    tb "gopkg.in/tucnak/telebot.v2"
    "log"
    "strconv"
    "time"
)

func main() {
    pictureURL := "https://t0922496-nus-13072021.s3.amazonaws.com/mosfet.png"
    b, err := tb.NewBot(tb.Settings{
        // You can also set custom API URL.
        // If field is empty it equals to "https://api.telegram.org".
        URL: "https://api.telegram.org",
        // Token here
        Token:  Token,
        Poller: &tb.LongPoller{Timeout: 10 * time.Second},
    })

    if err != nil {
        log.Fatal(err)
        return
    }

    b.Handle("/render", func(m *tb.Message) {
        a := &tb.Photo{
            File:    tb.FromURL(pictureURL),
            Caption: m.Payload,
        }
        resultMsg, err := b.Send(m.Sender, a)
        if err != nil {
            log.Println(err)
        }
        if resultMsg != nil {
            log.Println(a.FileID)
        }
    })

    b.Handle("/ch", func(msg *tb.Message) {
        if msg.Payload != "" {
            pictureURL = msg.Payload
        }
    })

    b.Handle(tb.OnQuery, func(q *tb.Query) {
        privateChannelRecipient := &tb.User{ID: chanID}
        log.Println(q.Text)
        urls := []string{
            pictureURL,
        }
        picture := &tb.Photo{
            File: tb.FromURL(pictureURL),
        }

        if !picture.InCloud() {
            _, err := b.Send(privateChannelRecipient, picture)
            if err != nil {
                log.Println(err)
            }
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

        err = b.Answer(q, &tb.QueryResponse{
            Results:   results,
            CacheTime: 10, // a minute
        })

        if err != nil {
            log.Println(err)
        }
    })

    b.Start()
}
