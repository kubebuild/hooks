package web

import (
	"fmt"
	"os"
	"regexp"

	"net/http"

	bugsnag "github.com/bugsnag/bugsnag-go"
	"github.com/cloudflare/cfssl/log"
	"github.com/kubebuild/webhooks/pkg/graphql"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/singleflight"
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
	githubHook, _    = github.New()
	gitlabHook, _    = gitlab.New()
	bitbucketHook, _ = bitbucket.New()
)

//Web interface for web client
type Web struct {
	log           *logrus.Logger
	graphqlClient *graphql.Client
}

var requestGroup singleflight.Group

// NewWeb instantiate web
func NewWeb(log *logrus.Logger, graphqlClient *graphql.Client) *Web {
	web := &Web{
		log:           log,
		graphqlClient: graphqlClient,
	}

	bugsnag.Configure(bugsnag.Configuration{
		APIKey: os.Getenv("BUGSNAG_API_KEY"),
		// The import paths for the Go packages
		// containing your source files
		ProjectPackages: []string{"main", "github.com/kubebuild/webhook"},
	})
	http.HandleFunc(healthPath, web.handleHealth)
	http.HandleFunc(githubPath, web.handleGithub)
	http.HandleFunc(gitlabPath, web.handleGitlab)
	http.HandleFunc(bitbucketPath, web.handleBitbucket)
	http.ListenAndServe(":9000", bugsnag.Handler(nil))
	return web
}

func parseRef(ref string) string {
	re := regexp.MustCompile("refs/\\w+/(.*)")
	branch := re.FindStringSubmatch(ref)[1]
	return branch
}

func (wb *Web) createBuild(token string, commit string, message string, branch string, userEmail string, userName string, pullRequestURL *string) {
	params := &graphql.BuildCreateParams{
		Token:          token,
		Commit:         commit,
		Message:        message,
		Branch:         branch,
		UserEmail:      userEmail,
		UserName:       userName,
		PullRequestURL: pullRequestURL,
	}
	wb.graphqlClient.CreateBuild(params)
}

func extractToken(r *http.Request) string {
	tokenArr, ok := r.URL.Query()["token"]
	if !ok || len(tokenArr[0]) < 1 {
		log.Error("Url Param 'token' is missing")
		return ""
	}
	return tokenArr[0]
}

func (wb *Web) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}
func (wb *Web) handleGithub(w http.ResponseWriter, r *http.Request) {
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
		branch := parseRef(push.Ref)
		group := fmt.Sprintf("%s-%s", push.After, branch)
		requestGroup.Do(group, func() (interface{}, error) {
			wb.createBuild(token,
				push.After,
				push.HeadCommit.Message,
				branch,
				push.HeadCommit.Author.Email,
				push.HeadCommit.Author.Name, nil)

			return nil, nil
		})

	case github.PullRequestPayload:
		pullRequest := payload.(github.PullRequestPayload)
		prLink := pullRequest.PullRequest.HTMLURL
		if pullRequest.Action == "synchronize" {
			branch := pullRequest.PullRequest.Head.Ref
			sha := pullRequest.PullRequest.Head.Sha
			group := fmt.Sprintf("%s-%s", sha, branch)
			requestGroup.Do(group, func() (interface{}, error) {
				wb.graphqlClient.UpdateBuild(&graphql.BuildUpdateParams{
					Token:          token,
					Commit:         sha,
					Branch:         branch,
					PullRequestURL: &prLink,
				})
				return nil, nil
			})

		}
		if pullRequest.Action == "opened" {
			wb.createBuild(token,
				pullRequest.PullRequest.Head.Sha,
				pullRequest.PullRequest.Title,
				pullRequest.PullRequest.Head.Ref,
				pullRequest.Sender.Login,
				pullRequest.Sender.Login, &prLink)
		}
	}
}

func (wb *Web) handleGitlab(w http.ResponseWriter, r *http.Request) {
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
		wb.createBuild(token,
			push.After,
			push.Commits[0].Message,
			parseRef(push.Ref),
			push.Commits[0].Author.Email,
			push.Commits[0].Author.Name, nil)
	}
}

func (wb *Web) handleBitbucket(w http.ResponseWriter, r *http.Request) {
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
		ref := parseRef(push.Push.Changes[0].New.Type)
		wb.createBuild(token, commit, message, ref, push.Actor.Username, push.Actor.Username, nil)
	}
}
