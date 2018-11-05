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
}

//BuildUpdateParams update build
type BuildUpdateParams struct {
	Token          string
	Commit         string
	Branch         string
	PullRequestURL *string
}

//UpdateBuild creates a build
func (c *Client) UpdateBuild(params *BuildUpdateParams) {
	c.Log.WithFields(logrus.Fields{
		"commit":         params.Commit,
		"branch":         params.Branch,
		"pullRequestUrl": params.PullRequestURL,
	}).Debug("createBuild")
	var buildMutation struct {
		UpddateBuildByToken struct {
			Successful gql.Boolean
		} `graphql:"updateBuildByToken(token: $token, commit: $commit, branch: $branch, pullRequestUrl: $pullRequestUrl)"`
	}
	variables := map[string]interface{}{
		"token":          params.Token,
		"commit":         params.Commit,
		"branch":         params.Branch,
		"pullRequestUrl": params.PullRequestURL,
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
	}).Debug("createBuild")
	var buildMutation struct {
		CreateBuildByToken struct {
			Successful gql.Boolean
		} `graphql:"createBuildByToken(token: $token, commit: $commit, message: $message, branch: $branch, userEmail: $userEmail, userName: $userName, pullRequestUrl: $pullRequestUrl)"`
	}
	variables := map[string]interface{}{
		"token":          params.Token,
		"commit":         params.Commit,
		"message":        params.Message,
		"branch":         params.Branch,
		"userEmail":      params.UserEmail,
		"userName":       params.UserName,
		"pullRequestUrl": params.PullRequestURL,
	}
	err := c.GraphqlClient.Mutate(context.Background(), &buildMutation, variables)
	if err != nil {
		c.Log.Error(err)
	}
	c.Log.WithField("successful", buildMutation.CreateBuildByToken.Successful).Debug("create result")
}
