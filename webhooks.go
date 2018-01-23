package main

import (
	"gopkg.in/go-playground/webhooks.v3/github"
	"gopkg.in/go-playground/webhooks.v3"
	"os"
	"fmt"
)

var DefaultLog webhooks.Logger = webhooks.NewLogger(false)

func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		return "8080"
	}
	return port
}

func halfPipeInChanged(changes []string) bool {
	for _, file := range changes {
		if file == ".halfpipe.io" {
			return true
		}
	}
	return false
}

func isHalfPipeCommit(payload github.PushPayload) bool {
	for _, commit := range payload.Commits {
		if halfPipeInChanged(commit.Added) || halfPipeInChanged(commit.Modified) || halfPipeInChanged(commit.Removed) {
			return true
		}
	}
	return false
}

func main() {
	DefaultLog.Info("Starting webhook receiver")
	hook := github.New(&github.Config{Secret: "asdf"})
	hook.RegisterEvents(HandlePushEvent, github.PushEvent, github.PushEvent)

	err := webhooks.Run(hook, ":"+getPort(), "/webhooks")
	if err != nil {
		DefaultLog.Error(err.Error())
	}
}

func HandlePushEvent(payload interface{}, header webhooks.Header) {
	switch payload.(type) {
	case github.PushPayload:
		push := payload.(github.PushPayload)
		if isHalfPipeCommit(push) {
			DefaultLog.Info("Got push payload that touches .halfpipe.io")
		} else {
			DefaultLog.Info("Got push payload that doesnt touch .halfpipe.io")
		}
	case github.PingPayload:
		ping := payload.(github.PingPayload)
		DefaultLog.Info("Got ping payload")
		DefaultLog.Info(fmt.Sprintf("%+v", ping))
	}
}
