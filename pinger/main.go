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

var lastSuccessMap = make(map[string]time.Time)

func pingIP(ip string) bool {
	cmd := exec.Command("ping", "-c", "1", ip)
	err := cmd.Run()
	return err == nil
}

func sendPingResult(result PingResult) {
	jsonData, err := json.Marshal(result)
	if err != nil {
		log.Println("Ошибка сериализации:", err)
		return
	}

	resp, err := http.Post("http://backend:8080/add-ping-result", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println("Ошибка при отправке запроса:", err)
		return
	}
	defer resp.Body.Close()
}

func main() {
	ips := []string{"8.8.8.8", "77.88.8.8"}
	for {
		for _, ip := range ips {
			fmt.Println("Пингуется IP:", ip)
			success := pingIP(ip)
			now := time.Now()

			if success {
				lastSuccessMap[ip] = now
			}

			result := PingResult{
				IP:          ip,
				PingTime:    now,
				LastSuccess: lastSuccessMap[ip],
			}

			sendPingResult(result)
		}
		time.Sleep(10 * time.Second)
	}
}
