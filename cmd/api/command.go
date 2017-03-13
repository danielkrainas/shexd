package api

import (
	"context"
	"fmt"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/danielkrainas/gobag/cmd"
	cfg "github.com/danielkrainas/gobag/configuration"
	"github.com/danielkrainas/gobag/context"

	"github.com/danielkrainas/shex/api/server"
	"github.com/danielkrainas/shex/registry/actions"
	"github.com/danielkrainas/shex/registry/configuration"
	storage "github.com/danielkrainas/shex/registry/storage/loader"
)

func init() {
	cmd.Register("api", Info)
}

func run(ctx context.Context, args []string) error {
	config, err := configuration.Resolve(args)
	if err != nil {
		return err
	}

	ctx, err = configureLogging(ctx, config)
	if err != nil {
		return fmt.Errorf("error configuring logging: %v", err)
	}

	log := acontext.GetLogger(ctx)
	log.Info("initializing server")
	actionPack, err := actions.FromConfig(config)
	if err != nil {
		return err
	}

	s, err := server.New(ctx, config.HTTP, actionPack)
	if err != nil {
		return err
	}

	log.Infof("using %q logging formatter", config.Log.Formatter)
	storage.LogSummary(ctx, config)
	return s.ListenAndServe()
}

var (
	Info = &cmd.Info{
		Use:   "api",
		Short: "run the api server",
		Long:  "Run the api server.",
		Run:   cmd.ExecutorFunc(run),
	}
)

func configureLogging(ctx context.Context, config *configuration.Config) (context.Context, error) {
	log.SetLevel(logLevel(config.Log.Level))
	formatter := config.Log.Formatter
	if formatter == "" {
		formatter = "text"
	}

	switch formatter {
	case "json":
		log.SetFormatter(&log.JSONFormatter{
			TimestampFormat: time.RFC3339Nano,
		})

	case "text":
		log.SetFormatter(&log.TextFormatter{
			TimestampFormat: time.RFC3339Nano,
		})

	default:
		if config.Log.Formatter != "" {
			return ctx, fmt.Errorf("unsupported log formatter: %q", config.Log.Formatter)
		}
	}

	if len(config.Log.Fields) > 0 {
		var fields []interface{}
		for k := range config.Log.Fields {
			fields = append(fields, k)
		}

		ctx = acontext.WithValues(ctx, config.Log.Fields)
		ctx = acontext.WithLogger(ctx, acontext.GetLogger(ctx, fields...))
	}

	ctx = acontext.WithLogger(ctx, acontext.GetLogger(ctx))
	return ctx, nil
}

func logLevel(level cfg.LogLevel) log.Level {
	l, err := log.ParseLevel(string(level))
	if err != nil {
		l = log.InfoLevel
		log.Warnf("error parsing level %q: %v, using %q", level, err, l)
	}

	return l
}
