package graphql

import (
	"github.com/shurcooL/graphql"
	"github.com/sirupsen/logrus"
)

// Client hold mutation actions
type Client struct {
	GraphqlClient *graphql.Client
	Log           *logrus.Logger
}

// NewGraphqlClient Return a new graphql client
func NewGraphqlClient(graphqlURL string, log *logrus.Logger) *Client {
	gqlClient := graphql.NewClient(graphqlURL, nil)
	return &Client{
		GraphqlClient: gqlClient,
		Log:           log,
	}
}
