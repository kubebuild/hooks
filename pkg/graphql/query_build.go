package graphql

import (
	"context"

	"github.com/shurcooL/graphql"
)

//Build Struct
type Build struct {
	ID             graphql.ID
	IsPullRequest  *graphql.Boolean
	PullRequestUrl *graphql.String
}

// BuildQuery query for builds
type BuildQuery struct {
	Build *Build `graphql:"buildForBranch(token: $token, branch: $branch)"`
}

//QueryBuildForBranch Return build for branch
func (c *Client) QueryBuildForBranch(branch string, token string) (*Build, error) {

	q := &BuildQuery{}
	variables := map[string]interface{}{
		"branch": branch,
		"token":  token,
	}
	err := c.GraphqlClient.Query(context.Background(), q, variables)
	if err != nil {
		c.Log.WithError(err).Error("BuildQuery Failed")
		return nil, err
	}
	return q.Build, nil
}
