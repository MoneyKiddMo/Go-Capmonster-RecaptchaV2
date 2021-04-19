package main

import (
	"fmt"
	"net/http"
	"log"
	"io/ioutil"
	"encoding/json"
	"strconv"
	"bytes"
	"time"
)

type capResponse struct {
	Errorid          int    `json:"errorId"`
	Errorcode        string `json:"errorCode"`
	Errordescription string `json:"errorDescription"`
	Taskid           int  `json:"taskId"`
}

type capSuccess struct {
	Errorid          int         `json:"errorId"`
	Errorcode        interface{} `json:"errorCode"`
	Errordescription interface{} `json:"errorDescription"`
	Solution         struct {
		Grecaptcharesponse string `json:"gRecaptchaResponse"`
	} `json:"solution"`
	Status string `json:"status"`
}

func capTask(reqURL string) int {	
	apiKey := "apikey" // Your CapMonster API Key
	siteKey := "sitekey" // Insert Your Recaptcha SiteKey

	client := &http.Client{}

	var jsonData = []byte(`{"clientKey":"`+ apiKey + `","task":{"type":"NoCaptchaTaskProxyless","websiteURL":"` + reqURL + `","websiteKey":"` + siteKey + `"}}`)
	capURL := "https://api.capmonster.cloud/createTask"
	req, err := http.NewRequest("POST", capURL, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatalf("Fatal Error %s", err)
	}

	resp, err := client.Do(req)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	jsonResponse := capResponse{}
	jsonErr := json.Unmarshal(body, &jsonResponse)
	if jsonErr != nil {
		log.Fatalf("Fatal error %s", jsonErr)
	}
	
	fmt.Println("Successfully Retrieved Captcha ID From CapMonster Servers: " + strconv.Itoa(jsonResponse.Taskid))
	capID := jsonResponse.Taskid
	return capID
}

func grabCapResponse(capID int) string {
	apiKey := "apikey" // Your CapMonster API Key

	client := &http.Client{Timeout: 20 * time.Second}
	
	
	reqURL := "https://api.capmonster.cloud/getTaskResult"
	var jsonData = []byte(`{"clientKey":"` + apiKey + `","taskId":` + strconv.Itoa(capID) + `}`)
	req, err := http.NewRequest("POST", reqURL, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatalf("Fatal Error %s", err)
	}
	resp, err := client.Do(req)
	defer resp.Body.Close()
	
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Fatal Error %s",err)
	}
	
	jsonResponse := capSuccess{}
	jsonErr := json.Unmarshal(body, &jsonResponse)
	if jsonErr != nil {
		log.Fatalf("Fatal error %s", jsonErr)
	}

	capResp := jsonResponse.Solution.Grecaptcharesponse
	return capResp
}

func main() {
	fmt.Println("initiating capmonster")

	reqURL := "pageurl" // Your Page URL
	capID := capTask(reqURL)
	for i := 0; i < 1000; i++ {
		recapID := grabCapResponse(capID)
		if recapID != "" {
		fmt.Println("[200] Retrived Captcha Token: " + grabCapResponse(capID))
		break
		}

	}	
}