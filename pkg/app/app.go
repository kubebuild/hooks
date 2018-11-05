package app

import (
	"os"
	"time"

	"github.com/kubebuild/webhooks/pkg/graphql"
	"github.com/kubebuild/webhooks/pkg/utils"
	"github.com/kubebuild/webhooks/pkg/web"
	"github.com/sirupsen/logrus"
)

// App provides a default app structure with Logger
type App struct {
	Config        Configurer
	Log           *logrus.Logger
	GraphqlClient *graphql.Client
}

// NewApp configures and returns an App
func NewApp(config Configurer) (*web.Web, error) {
	app := &App{
		Config: config,
	}

	logger, err := newLogger(config)
	if err != nil {
		return nil, err
	}
	app.Log = logger
	app.GraphqlClient = newGraphqlClient(config, app.Log)
	app.Log.WithFields(
		logrus.Fields{
			"version":     config.GetVersion(),
			"name":        config.GetName(),
			"grapqhl-url": config.GetGraphqlURL(),
		}).Info("Starting app ...")
	return web.NewWeb(app.Log, app.GraphqlClient), nil
}

func newLogger(config Configurer) (*logrus.Logger, error) {
	log := logrus.New()
	log.Formatter = utils.NewLogFormatter(config.GetName(), config.GetVersion(), &logrus.TextFormatter{TimestampFormat: time.RFC3339Nano, FullTimestamp: true})
	log.Level = config.GetLogLevel()
	log.Out = os.Stdout
	return log, nil
}

func newGraphqlClient(config Configurer, log *logrus.Logger) *graphql.Client {
	client := graphql.NewGraphqlClient(config.GetGraphqlURL(), log)
	return client
}
