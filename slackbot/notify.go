package slackbot

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	cloudbuild "google.golang.org/api/cloudbuild/v1"
)

func NotifyStart(id string, webhook string, jobName string) {
	message := getMessage("START", id, jobName)
	postActually(message, webhook)
}

// Notify posts a notification to Slack that the build is complete.
func NotifyFinish(b *cloudbuild.Build, webhook string, jobName string) {
	message := getMessage(b.Status, b.Id, jobName)
	postActually(message, webhook)
}

func getMessage(status string, id string, jobName string) (message string) {
	url := fmt.Sprintf("https://console.cloud.google.com/cloud-build/builds/%s", id)
	iconUrl := "https://mpng.pngfly.com/20180329/qjq/kisspng-google-cloud-platform-google-compute-engine-kubern-container-5abc828e10c6a8.2707130315223036300687.jpg"

	var icon, color, desc string
	switch status {
	case "START":
		icon = "hourglass_flowing_sand"
		color = "#eac8fd"
		desc = "Build starting..."
	case "SUCCESS":
		icon = "white_check_mark"
		color = "#00ff1b"
		desc = "Build succeeded!!!"
	case "FAILURE":
		icon = "x"
		color = "#fa0000"
		desc = "Build failed!!!"
	case "CANCELLED":
		icon = "wastebasket"
		color = "#cacaca"
		desc = "Build cancelled..."
	case "TIMEOUT":
		icon = "hourglass"
		color = "#fa0000"
		desc = "Timeout..."
	case "STATUS_UNKNOWN", "INTERNAL_ERROR":
		icon = "interrobang"
		color = "#fa0000"
		desc = "Unknow error occured"
	default:
		icon = "question"
		color = "#cacaca"
	}

	msgFmt := `{
		"username":"Cloud Build",
		"icon_url":"%s",
		"attachments": [{
			"pretext": ":loudspeaker: DEPLOYMENT %s",
			"title": "Job: %s",
			"text":"● ID: %s \n:%s: %s \n[(Open details)](%s)",
			"color":"%s"
		}]
	}`
	message = fmt.Sprintf(msgFmt, iconUrl, status, jobName, id, icon, desc, url, color)
	return
}

func postActually(message string, webhook string) {
	reader := strings.NewReader(message)
	resp, err := http.Post(webhook, "application/json", reader)
	if err != nil {
		log.Fatalf("Failed to post to Slack: %v", err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	log.Printf("Posted message to Slack: [%v], got response [%s]", message, body)
}
