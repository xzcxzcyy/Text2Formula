package main

import (
    "fmt"
    tb "gopkg.in/tucnak/telebot.v2"
    "log"
    "os"
    "strconv"
    "time"
)

func main() {
    f, err := os.OpenFile("go-hello.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
    if err != nil {
        log.Fatalf("Error opening file: %v", err)
    }
    defer f.Close()
    log.SetOutput(f)
    log.Println("Start printing log.")

    myBot, err := tb.NewBot(tb.Settings{
        // You can also set custom API URL.
        // If field is empty it equals to "https://api.telegram.org".
        URL: "https://api.telegram.org",
        // token here
        Token:  token,
        Poller: &tb.LongPoller{Timeout: 10 * time.Second},
    })

    if err != nil {
        log.Fatal(err)
        return
    }

    myBot.Handle("/render", func(m *tb.Message) {
        pngFilePath, _, err := renderTex(fmt.Sprintf("%v-%v", m.Chat.ID, m.ID), m.Payload)
        if err != nil {
            log.Printf("Error: On command render: %v", err)
            return
        }
        a := &tb.Photo{
            File:    tb.FromDisk(pngFilePath),
            Caption: m.Payload,
        }
        resultMsg, err := myBot.Send(m.Sender, a)
        if err != nil {
            log.Println(err)
        }
        if resultMsg != nil {
            log.Println(a.FileID)
        }
    })

    myBot.Handle(tb.OnQuery, func(q *tb.Query) {
        /**
          tmp
        */
        // pictureURL := "https://en.wikipedia.org/wiki/Scalable_Vector_Graphics#/media/File:SVG_Logo.svg"

        //privateChannelRecipient := &tb.User{ID: chanID}
        queryID := q.ID

        if q.Text == "" {
            return
        }
        log.Printf("ID: %v, Query: %v", queryID, q.Text)

        curAnswerFilePath, sizeInfo, err := renderTex(queryID, q.Text)
        if err != nil {
            log.Printf("Error: On Handle tb.OnQuery: When rendering: %v", err)
            return
        }
        //log.Printf("tb.Onquery handler gets image size: %v", sizeInfo)
        //log.Printf("On Handle tb.OnQuery: %v", curAnswerFilePath)

        answerURL := uploadFileToS3(bucketID, curAnswerFilePath)

        urls := []string{
            answerURL,
        }
        //picture := &tb.Photo{
        //    File: tb.FromURL(answerURL),
        //}
        //
        //if !picture.InCloud() {
        //    _, err := myBot.Send(privateChannelRecipient, picture)
        //    if err != nil {
        //        log.Println(err)
        //    }
        //}

        results := make(tb.Results, len(urls)) // []tb.Result
        for i, url := range urls {
            result := &tb.PhotoResult{
                ResultBase:  tb.ResultBase{},
                URL:         url,
                Width:       sizeInfo.X,
                Height:      sizeInfo.Y,
                Title:       "",
                Description: "",
                Caption:     q.Text,
                ParseMode:   "",
                ThumbURL:    url,
            }

            results[i] = result
            // needed to set a unique string ID for each result
            results[i].SetResultID(strconv.Itoa(i))
        }

        err = myBot.Answer(q, &tb.QueryResponse{
            Results:   results,
            CacheTime: 60, // a minute
        })

        if err != nil {
            log.Println(err)
        }
    })

    myBot.Start()
}
