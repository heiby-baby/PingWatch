package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"time"
)

type PingResult struct {
	IP          string    `json:"ip"`
	PingTime    time.Time `json:"ping_time"`
	LastSuccess time.Time `json:"last_success"`
}

func pingIP(ip string) bool {
	cmd := exec.Command("ping", "-c", "1", ip)
	err := cmd.Run()
	return err == nil
}

func sendPingResult(result PingResult) {
	jsonData, err := json.Marshal(result)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := http.Post("http://backend:8080/add-ping-result", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
}

func main() {
	ips := []string{"8.8.8.8", "77.88.8.8"}

	for {
		for _, ip := range ips {
			fmt.Println("Пингуется ip: ", ip)
			success := pingIP(ip)
			result := PingResult{
				IP:          ip,
				PingTime:    time.Now(),
				LastSuccess: time.Now(),
			}
			if !success {
				result.LastSuccess = time.Time{}
			}
			sendPingResult(result)
		}
		time.Sleep(10 * time.Second)
	}
}
