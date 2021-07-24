package network

import (
    "banson.moe/t2f-bot/config"
    "bytes"
    "encoding/json"
    "github.com/pkg/errors"
    "io/ioutil"
    "log"
    "net/http"
)

type PictureRender struct {
    QueryID string `json:"query_id"`
    Formula string `json:"formula"`
}

type ResponseBody struct {
    S3Url  string `json:"s3Url"`
    Width  int    `json:"width"`
    Height int    `json:"height"`
}

func Request(queryID, formula string) (*ResponseBody, error) {

    pictureRender := PictureRender{
        QueryID: queryID,
        Formula: formula,
    }

    bodyBytes, _ := json.Marshal(pictureRender)
    requestBody := bytes.NewReader(bodyBytes)

    req, _ := http.NewRequest(http.MethodPost, config.RenderServerHost, requestBody)
    req.Header.Add("Content-Type", "application/json")
    req.Header.Add("Accept", "*/*")

    newClient := http.Client{}
    resp, _ := newClient.Do(req)

    if resp.StatusCode != http.StatusOK {
        log.Println("during network/Request: response get error")
        return nil, errors.New("response's status is not OK")
    }
    body, _ := ioutil.ReadAll(resp.Body)
    responseBody := ResponseBody{}
    err := json.Unmarshal(body, &responseBody)
    if err != nil {
        log.Printf("during network request unmarshal: %v", err)
        return nil, err
    }
    log.Println(responseBody.S3Url)
    return &responseBody, nil
}
