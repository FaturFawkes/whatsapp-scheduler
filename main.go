package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/enescakir/emoji"
	"github.com/robfig/cron/v3"
)

func main() {
	jktTime, _ := time.LoadLocation("Asia/Jakarta")
	scheduler := cron.New(cron.WithLocation(jktTime))
	defer scheduler.Stop()

	message := Absent{
		Level:     "reminder",
		Timestamp: time.Now().Format(time.RFC3339),
		Message:   fmt.Sprintf("Absen Clock Out guys %v", emoji.WinkingFace),
	}
	msg, _ := json.Marshal(message)

	scheduler.AddFunc("5 18 * * 1-5", func() {
		clockOut(string(msg))
	})
	scheduler.AddFunc("30 4 * * *", waSender)


	go scheduler.Start()
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
}

type Absent struct {
	Level     string `json:"event"`
	Timestamp string `json:"time"`
	Message   string `json:"message"`
}

func clockOut(msg string) {
	apiurl := "https://api.ultramsg.com/instance64023/messages/chat"
	data := url.Values{}
	data.Set("token", os.Getenv("TOKEN"))
	data.Set("to", os.Getenv("GROUP_ID"))
	data.Set("body", "```" + msg + "```")

	payload := strings.NewReader(data.Encode())

	req, _ := http.NewRequest("POST", apiurl, payload)

	req.Header.Add("content-type", "application/x-www-form-urlencoded")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	fmt.Println(string(body))
}

func waSender() {
	apiurl := "https://api.ultramsg.com/instance64023/messages/chat"
	data := url.Values{}
	data.Set("token", os.Getenv("TOKEN"))
	data.Set("to", os.Getenv("NUMBER"))
	data.Set("body", fmt.Sprintf("```" + "Selamat Pagi Bennica San %v```", emoji.SmilingFaceWithHalo))

	payload := strings.NewReader(data.Encode())

	req, _ := http.NewRequest("POST", apiurl, payload)

	req.Header.Add("content-type", "application/x-www-form-urlencoded")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	fmt.Println(string(body))

}