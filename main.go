package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
)

var key string
var secret string
var domain string
var domainType string
var name string
var ttl int
var url string
var newIp string

type GodaddyRecord struct {
	Data string `json:"data"`
	Name string `json:"name"`
	TTL  int    `json:"ttl"`
	Type string `json:"type"`
}

type GodaddyPutRecord struct {
	Data string `json:"data"`
	TTL  int    `json:"ttl"`
}

func main() {
	flag.StringVar(&key, "k", "", "godaddy developer key")
	flag.StringVar(&secret, "s", "", "godaddy developer secret")
	flag.StringVar(&domain, "d", "", "godaddy developer domain")
	flag.StringVar(&domainType, "t", "A", "domain type")
	flag.StringVar(&name, "n", "@", "domain name")
	flag.IntVar(&ttl, "T", 600, "domain name")
	flag.StringVar(&url, "u", "http://api.ipify.org", "default check url")
	flag.StringVar(&newIp, "N", "", "new ip")

	flag.Parse()

	if key == "" || secret == "" || domain == "" {
		flag.PrintDefaults()
		return
	}
	var ip string
	if newIp == "" {
		log.Println(key, secret, domain, domainType, name, ttl, url)
		resp, err := http.Get(url)
		if err != nil {
			log.Fatal(err)
		}
		defer func() {
			_ = resp.Body.Close()
		}()

		ipBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		ip = string(ipBytes)
	} else {
		ip = newIp
	}

	log.Printf("current ip = %s", ip)

	godaddyUrl := fmt.Sprintf("https://api.godaddy.com/v1/domains/%s/records/%s/%s", domain, domainType, name)
	request, err := http.NewRequest("GET", godaddyUrl, nil)
	if err != nil {
		log.Fatal(err)
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", fmt.Sprintf("sso-key %s:%s", key, secret))

	remoteResp, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		_ = remoteResp.Body.Close()
	}()

	remoteBytes, err := io.ReadAll(remoteResp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var respBody []GodaddyRecord
	err = json.Unmarshal(remoteBytes, &respBody)
	if err != nil {
		log.Fatal(err)
	}

	var hasCurrentIp = false
	for _, record := range respBody {
		log.Printf("remote ip = %s", record.Data)

		if record.Data == ip {
			hasCurrentIp = true
			break
		}
	}

	if hasCurrentIp {
		log.Println("No update required")
		return
	}

	putBody := make([]GodaddyPutRecord, 0)
	putBody = append(putBody, GodaddyPutRecord{
		Data: ip,
		TTL:  ttl,
	})

	putBodyBytes, err := json.Marshal(putBody)
	log.Println(string(putBodyBytes))
	updateRequest, err := http.NewRequest("PUT", godaddyUrl, bytes.NewReader(putBodyBytes))
	if err != nil {
		log.Fatal(err)
	}

	updateRequest.Header.Set("Content-Type", "application/json")
	updateRequest.Header.Set("Authorization", fmt.Sprintf("sso-key %s:%s", key, secret))

	resultResp, err := http.DefaultClient.Do(updateRequest)
	if err != nil {
		log.Fatal(err)
	}

	if resultResp.StatusCode == 200 {
		log.Println("update success")
		return
	}

	resultBodyBytes, err := io.ReadAll(resultResp.Body)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("update failed: " + string(resultBodyBytes))
}
