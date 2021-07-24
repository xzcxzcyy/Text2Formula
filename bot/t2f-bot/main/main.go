package main

import (
    "banson.moe/t2f-bot/config"
    "banson.moe/t2f-bot/database"
    "banson.moe/t2f-bot/network"
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

    dbClient := database.InitDatabase()

    myBot, err := tb.NewBot(tb.Settings{
        // You can also set custom API URL.
        // If field is empty it equals to "https://api.telegram.org".
        URL: "https://api.telegram.org",
        // token here
        Token:  config.Token,
        Poller: &tb.LongPoller{Timeout: 10 * time.Second},
    })

    if err != nil {
        log.Fatal(err)
        return
    }

    myBot.Handle("/render", func(m *tb.Message) {
        //pngFilePath, _, err := renderer.RenderTex(fmt.Sprintf("%v-%v", m.Chat.ID, m.ID), m.Payload)
        if m.Payload == "" {
            return
        }

        pictureInfo, err := dbClient.Get(m.Payload)
        if err != nil {
            log.Printf("Debug: On command \"/render\": %v", err)
        }

        var pictureUrl string
        if pictureInfo != nil {
            pictureUrl = pictureInfo.S3Url
        } else {
            renderResp, err := network.Request(fmt.Sprintf("%v-%v", m.Chat.ID, m.ID), m.Payload)
            if err != nil {
                log.Printf("Error: On command \"/render\": %v", err)
                return
            }
            err = dbClient.Put(m.Payload, &database.PictureInfo{
                S3Url:  renderResp.S3Url,
                Width:  renderResp.Width,
                Height: renderResp.Height,
            })
            if err != nil {
                log.Printf("Error: When PUT data: %v", err)
            }
            pictureUrl = renderResp.S3Url
        }

        a := &tb.Photo{
            File:    tb.FromURL(pictureUrl),
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

        queryID := q.ID

        if q.Text == "" {
            return
        }
        log.Printf("ID: %v, Query: %v", queryID, q.Text)

        pictureInfo, err := dbClient.Get(q.Text)
        if err != nil {
            log.Printf("Debug: OnQuery: %v", err)
        }

        if pictureInfo == nil {
            renderResp, err := network.Request(queryID, q.Text)
            if err != nil {
                log.Printf("Error: OnQuery: %v", err)
                return
            }
            pictureInfo = &database.PictureInfo{
                S3Url:  renderResp.S3Url,
                Width:  renderResp.Width,
                Height: renderResp.Height,
            }
            err = dbClient.Put(q.Text, pictureInfo)
            if err != nil {
                log.Printf("Error: When PUT data: %v", err)
            }
        }

        urls := []string{
            pictureInfo.S3Url,
        }

        results := make(tb.Results, len(urls)) // []tb.Result
        for i, url := range urls {
            result := &tb.PhotoResult{
                ResultBase:  tb.ResultBase{},
                URL:         url,
                Width:       pictureInfo.Width,
                Height:      pictureInfo.Height,
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
