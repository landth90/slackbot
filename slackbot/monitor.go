package slackbot

import (
	"context"
	"log"
	"time"
)

const maxErrors = 3
const tickDuration = 20 * time.Second

// 1m20s
const monitorErrorMarginDuration = (maxErrors + 1) * tickDuration

// Monitor polls Cloud Build until the build reaches completed status, then triggers the Slack event.
func Monitor(ctx context.Context, projectId string, buildId string, webhook string, jobName string) {
	NotifyStart(buildId, webhook, jobName)

	svc := gcbClient(ctx)
	errors := 0

	t := time.Tick(tickDuration)
	for {
		log.Printf("Polling build %s", buildId)
		getMonitoredBuild := svc.Projects.Builds.Get(projectId, buildId)
		monitoredBuild, err := getMonitoredBuild.Do()
		if err != nil {
			if errors <= maxErrors {
				log.Printf("Failed to get build details from Cloud Build.  Will retry in %s", tickDuration)
				errors++
				continue
			} else {
				log.Fatalf("Reached maximum number of errors (%d).  Exiting", maxErrors)
			}
		}
		switch monitoredBuild.Status {
		case "SUCCESS", "FAILURE", "INTERNAL_ERROR", "TIMEOUT", "CANCELLED":
			log.Printf("Terminal status reached.  Notifying")
			NotifyFinish(monitoredBuild, webhook, jobName)
			return
		}
		<-t
	}
}
