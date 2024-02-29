package main

import (
	"bytes"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func readinessProbe(cfg *config, status chan bool, client *http.Client) {
	var i int16
	for i = 1; i < cfg.MaxRetry; i++ {
		sts := checkAvailability(cfg.Hostname, cfg.AdminPort, client)
		if sts == 200 {
			log.Println("Openmetadata is available")
			status <- true
			break
		}
		log.Println("Waiting for openmetadata")
	}
	status <- false
}

func checkAvailability(hostname string, adminPort int32, client *http.Client) (status int) {
	resp, err := client.Get(fmt.Sprintf("%s:%d/healthcheck", hostname, adminPort))
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Printf("status code error: %d %s", resp.StatusCode, resp.Status)
	}
	return resp.StatusCode
}

func getTokenWithUserPass(cfg *config, token *jwtToken, client *http.Client) {
	login := login{
		Email:    cfg.Username,
		Password: b64.StdEncoding.EncodeToString([]byte(cfg.Password)),
	}
	marshalled, err := json.Marshal(login)
	if err != nil {
		log.Fatalf("impossible to marshall login: %s", err)
	}
	url := fmt.Sprintf("%s:%d/api/v1/users/login", cfg.Hostname, cfg.Port)
	req, err := http.NewRequest("POST", url, bytes.NewReader(marshalled))
	if err != nil {
		log.Fatalf("impossible to build request: %s", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("impossible to send request: %s", err)
	}
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("status code error: %d %s", resp.StatusCode, resp.Status)
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("impossible to read all body of response: %s", err)
	}
	// log.Printf("res body: %s", string(resBody))
	if err := json.Unmarshal(respBody, &token); err != nil {
		log.Fatalf("impossible to unmarshal json %s:", err)
	}
}

func getUserIdByName(cfg *config, token *jwtToken, client *http.Client) string {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s:%d/api/v1/users/name/%s", cfg.Hostname, cfg.Port, "ingestion-bot"), nil)
	if err != nil {
		log.Fatalf("impossible to build request: %s", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("impossible to send request: %s", err)
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("impossible to read all body of response: %s", err)
	}
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("status code error: %d %s", resp.StatusCode, resp.Status)
	}
	// log.Printf("res body: %s", string(resBody))
	user := userId{}
	if err := json.Unmarshal(respBody, &user); err != nil {
		log.Fatalf("impossible to unmarshal json %s:", err)
	}
	return user.Id
}

func getBotTokenById(cfg *config, token *jwtToken, client *http.Client, id string) string {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s:%d/api/v1/users/token/%s", cfg.Hostname, cfg.Port, id), nil)
	if err != nil {
		log.Fatalf("impossible to build request: %s", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("impossible to send request: %s", err)
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("impossible to read all body of response: %s", err)
	}
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("status code error: %d %s", resp.StatusCode, resp.Status)
	}
	// log.Printf("res body: %s", string(resBody))
	bot := botJwtToken{}
	if err := json.Unmarshal(respBody, &bot); err != nil {
		log.Fatalf("impossible to unmarshal json %s:", err)
	}
	return bot.JWTToken
}
