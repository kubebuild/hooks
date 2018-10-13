package main

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"time"

	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/shurcooL/graphql"
	"gopkg.in/go-playground/webhooks.v5/bitbucket"
	"gopkg.in/go-playground/webhooks.v5/github"
	"gopkg.in/go-playground/webhooks.v5/gitlab"
)

const (
	githubPath    = "/github"
	gitlabPath    = "/gitlab"
	bitbucketPath = "/bitbucket"
	healthPath    = "/"
)

var (
	log              = newLogger()
	githubHook, _    = github.New()
	gitlabHook, _    = gitlab.New()
	bitbucketHook, _ = bitbucket.New()
	// graphqlClient    = graphql.NewClient("https://api.kubebuild.com/graphql", nil)
	graphqlClient = graphql.NewClient("http://localhost:4000/graphql", nil)
)

func newLogger() *logrus.Logger {
	logLevel, _ := logrus.ParseLevel("info")
	log := logrus.New()
	log.Formatter = &logrus.JSONFormatter{TimestampFormat: time.RFC3339Nano}
	log.Level = logLevel
	log.Out = os.Stdout
	return log
}
func main() {
	http.HandleFunc(healthPath, handleHealth)
	http.HandleFunc(githubPath, handleGithub)
	http.HandleFunc(gitlabPath, handleGitlab)
	http.HandleFunc(bitbucketPath, handleBitbucket)
	http.ListenAndServe(":9000", nil)
}

func parseRef(ref string) string {
	re := regexp.MustCompile("refs/\\w+/(.*)")
	branch := re.FindStringSubmatch(ref)[1]
	return branch
}

func createBuild(token string, commit string, message string, ref string, userEmail string, userName string) {
	branch := parseRef(ref)
	log.WithFields(logrus.Fields{
		"commit":    commit,
		"message":   message,
		"branch":    branch,
		"userEmail": userEmail,
		"userName":  userName,
	}).Info("createBuild")
	var buildMutation struct {
		CreateBuildByToken struct {
			Successful graphql.Boolean
		} `graphql:"createBuildByToken(token: $token, commit: $commit, message: $message, branch: $branch, userEmail: $userEmail, userName: $userName)"`
	}
	variables := map[string]interface{}{
		"token":     token,
		"commit":    commit,
		"message":   message,
		"branch":    branch,
		"userEmail": userEmail,
		"userName":  userName,
	}
	err := graphqlClient.Mutate(context.Background(), &buildMutation, variables)
	if err != nil {
		log.Error(err)
	}
	log.WithField("successful", &buildMutation.CreateBuildByToken.Successful).Info("create result")
}

func extractToken(r *http.Request) string {
	tokenArr, ok := r.URL.Query()["token"]
	if !ok || len(tokenArr[0]) < 1 {
		log.Error("Url Param 'token' is missing")
		return ""
	}
	return tokenArr[0]
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}
func handleGithub(w http.ResponseWriter, r *http.Request) {
	token := extractToken(r)
	payload, err := githubHook.Parse(r, github.PushEvent, github.PullRequestEvent)
	if err != nil {
		if err == github.ErrEventNotFound {
			// ok event wasn;t one of the ones asked to be parsed
		}
	}
	switch payload.(type) {

	case github.PushPayload:
		push := payload.(github.PushPayload)
		createBuild(token,
			push.After,
			push.HeadCommit.Message,
			push.Ref,
			push.HeadCommit.Author.Email,
			push.HeadCommit.Author.Name)

	case github.PullRequestPayload:
		pullRequest := payload.(github.PullRequestPayload)
		if pullRequest.Action == "opened" {
			createBuild(token,
				pullRequest.PullRequest.Head.Sha,
				pullRequest.PullRequest.Body,
				pullRequest.PullRequest.Head.Ref,
				pullRequest.Sender.Login,
				pullRequest.Sender.Login)
		}
	}
}

func handleGitlab(w http.ResponseWriter, r *http.Request) {
	token := extractToken(r)
	payload, err := gitlabHook.Parse(r, gitlab.PushEvents)
	if err != nil {
		if err == github.ErrEventNotFound {
			// ok event wasn;t one of the ones asked to be parsed
		}
	}
	switch payload.(type) {

	case gitlab.PushEventPayload:
		push := payload.(gitlab.PushEventPayload)
		createBuild(token,
			push.After,
			push.Commits[0].Message,
			push.Ref,
			push.Commits[0].Author.Email,
			push.Commits[0].Author.Name)
	}
}

func handleBitbucket(w http.ResponseWriter, r *http.Request) {
	token := extractToken(r)
	fmt.Println(token)
	payload, err := bitbucketHook.Parse(r, bitbucket.RepoPushEvent)
	if err != nil {
		if err == github.ErrEventNotFound {
			// ok event wasn;t one of the ones asked to be parsed
		}
	}
	switch payload.(type) {

	case bitbucket.RepoPushPayload:
		push := payload.(bitbucket.RepoPushPayload)
		commit := push.Push.Changes[0].New.Target.Hash
		message := push.Push.Changes[0].New.Target.Message
		ref := push.Push.Changes[0].New.Type
		createBuild(token, commit, message, ref, push.Actor.Username, push.Actor.Username)
	}
}
