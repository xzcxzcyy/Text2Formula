package main

import (
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
	S3Url string `json:"s3Url"`
}

func Request(queryID, formula string) (string, error) {

	pictureRender := PictureRender{
		QueryID: "1234566666",
		Formula: "hahahahahx_2",
	}

	bodyBytes, _ := json.Marshal(pictureRender)
	requestBody := bytes.NewReader(bodyBytes)

	req, _ := http.NewRequest(http.MethodPost, "http://localhost:6001/render", requestBody)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "*/*")


	newClient := http.Client{}
	resp, _ := newClient.Do(req)

	if resp.StatusCode != http.StatusOK {
		log.Println("response get error")
		return "", errors.New("response's status is not OK")
	}
	body, _ := ioutil.ReadAll(resp.Body)
	responseBody := ResponseBody{}
	err := json.Unmarshal(body, &responseBody)
 	if err != nil {
 		log.Println(err)
	}
	log.Println(responseBody.S3Url)
	return responseBody.S3Url, nil
}
