package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/joho/godotenv"
)

var CLOUDFLARE_API_KEY = ""
var ZONE_ID = ""

type IP struct {
	Query string
}

type ID struct {
	Result []struct {
		Id      string `json:"id"`
		Content string `json:"content"`
	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	CLOUDFLARE_API_KEY = os.Getenv("CLOUDFLARE_API_KEY")
	ZONE_ID = os.Getenv("ZONE_ID")
	fmt.Println("Starting...")
	s := gocron.NewScheduler(time.UTC)
	s.Every(60).Minute().Do(func() {
		ip := getip2()
		fmt.Println(ip)
		id, zip, err := getZoneId()

		if err != nil {
			fmt.Println(err)
		} else if zip != ip {
			updateCloudflare(ip, id)
		}

	})
	s.StartBlocking()
}

func getip2() string {
	req, err := http.Get("http://ip-api.com/json/")
	if err != nil {
		return err.Error()
	}
	defer req.Body.Close()

	body, err := io.ReadAll(req.Body)
	if err != nil {
		return err.Error()
	}
	var ip IP
	json.Unmarshal(body, &ip)
	return ip.Query
}

func getZoneId() (string, string, error) {
	url := "https://api.cloudflare.com/client/v4/zones/" + ZONE_ID + "/dns_records/"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		return "", "", err
	}
	req.Header.Set("Authorization", "Bearer "+CLOUDFLARE_API_KEY)

	res, err := client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", "", err
	}

	var id ID
	json.Unmarshal(body, &id)
	return id.Result[0].Id, id.Result[0].Content, nil

}

func updateCloudflare(ip string, id string) {
	ZONE_NAME := os.Getenv("ZONE_NAME")

	url := "https://api.cloudflare.com/client/v4/zones/" + ZONE_ID + "/dns_records/" + id
	fmt.Println(url)

	//payload := strings.NewReader("{content: " + ip + ",name: example.com,proxied: true,type: A, ttl: 1}")
	payload := []byte(`{"type":"A", "name":"` + ZONE_NAME + `", "content":"` + ip + `", "ttl":1,"proxied":true}`)
	fmt.Println(bytes.NewBuffer(payload))
	req, _ := http.NewRequest("PUT", url, bytes.NewBuffer(payload))

	req.Header.Add("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+CLOUDFLARE_API_KEY)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(res)
	fmt.Println(string(body))
}
