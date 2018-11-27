package graphql

import (
	"context"

	gql "github.com/shurcooL/graphql"
	"github.com/sirupsen/logrus"
)

//BuildCreateParams create build
type BuildCreateParams struct {
	Token          string
	Commit         string
	Branch         string
	Message        string
	UserEmail      string
	UserName       string
	PullRequestURL *string
	IsPullRequest  *gql.Boolean
	ErrorMessage   *string
}

//BuildUpdateParams update build
type BuildUpdateParams struct {
	Token          string
	Commit         string
	Branch         string
	PullRequestURL *string
	IsPullRequest  *gql.Boolean
}

//UpdateBuilds updates builds
func (c *Client) UpdateBuilds(params *BuildUpdateParams) {
	c.Log.WithFields(logrus.Fields{
		"commit":         params.Commit,
		"branch":         params.Branch,
		"pullRequestUrl": params.PullRequestURL,
		"isPullRequest":  params.IsPullRequest,
	}).Debug("updateBuild")
	var buildMutation struct {
		UpddateBuildByToken struct {
			Successful gql.Boolean
		} `graphql:"updateBuildsByToken(token: $token, commit: $commit, branch: $branch, pullRequestUrl: $pullRequestUrl, isPullRequest: $isPullRequest)"`
	}
	variables := map[string]interface{}{
		"token":          params.Token,
		"commit":         params.Commit,
		"branch":         params.Branch,
		"pullRequestUrl": params.PullRequestURL,
		"isPullRequest":  params.IsPullRequest,
	}
	err := c.GraphqlClient.Mutate(context.Background(), &buildMutation, variables)
	if err != nil {
		c.Log.Error(err)
	}
	c.Log.WithField("successful", buildMutation.UpddateBuildByToken.Successful).Debug("update result")
}

//CreateBuild creates a build
func (c *Client) CreateBuild(params *BuildCreateParams) {

	c.Log.WithFields(logrus.Fields{
		"commit":         params.Commit,
		"message":        params.Message,
		"branch":         params.Branch,
		"userEmail":      params.UserEmail,
		"userName":       params.UserName,
		"pullRequestUrl": params.PullRequestURL,
		"isPullRequest":  params.IsPullRequest,
	}).Debug("createBuild")
	var buildMutation struct {
		CreateBuildByToken struct {
			Successful gql.Boolean
		} `graphql:"createBuildByToken(token: $token, commit: $commit, message: $message, branch: $branch, userEmail: $userEmail, userName: $userName, pullRequestUrl: $pullRequestUrl, isPullRequest: $isPullRequest)"`
	}
	variables := map[string]interface{}{
		"token":          params.Token,
		"commit":         params.Commit,
		"message":        params.Message,
		"branch":         params.Branch,
		"userEmail":      params.UserEmail,
		"userName":       params.UserName,
		"pullRequestUrl": params.PullRequestURL,
		"isPullRequest":  params.IsPullRequest,
	}
	err := c.GraphqlClient.Mutate(context.Background(), &buildMutation, variables)
	if err != nil {
		c.Log.Error(err)
	}
	c.Log.WithField("successful", buildMutation.CreateBuildByToken.Successful).Debug("create result")
}
