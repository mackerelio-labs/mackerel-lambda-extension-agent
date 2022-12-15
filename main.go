package main

import (
	"context"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/mackerelio-labs/mackerel-lambda-extension-agent/dispatcher"
	"github.com/mackerelio-labs/mackerel-lambda-extension-agent/host/mackerel"
	"github.com/mackerelio-labs/mackerel-lambda-extension-agent/lambda/extension"
	"github.com/mackerelio-labs/mackerel-lambda-extension-agent/lambda/telemetry"
	"github.com/sirupsen/logrus"
)

var Logger *logrus.Entry

func init() {
	logLevelStr := strings.ToUpper(os.Getenv("EXT_LOG_LEVEL"))
	var logLevel logrus.Level
	switch logLevelStr {
	case "TRACE":
		logLevel = logrus.TraceLevel
	case "DEBUG":
		logLevel = logrus.DebugLevel
	case "INFO":
		logLevel = logrus.InfoLevel
	case "WARNING":
		logLevel = logrus.WarnLevel
	case "ERROR":
		logLevel = logrus.ErrorLevel
	case "FATAL":
		logLevel = logrus.FatalLevel
	case "PANIC":
		logLevel = logrus.PanicLevel
	default:
		logLevel = logrus.WarnLevel
	}
	logrus.SetLevel(logLevel)

	if logLevel >= logrus.DebugLevel {
		logrus.SetReportCaller(true)
	}

	extName := "mackerel-lambda-extension-agent"
	Logger = logrus.WithFields(logrus.Fields{"ext": extName, "pkg": "main"})
	mackerel.Logger = logrus.WithFields(logrus.Fields{"ext": extName, "pkg": "host/mackerel"})
	extension.Logger = logrus.WithFields(logrus.Fields{"ext": extName, "pkg": "lambda/extension"})
	telemetry.Logger = logrus.WithFields(logrus.Fields{"ext": extName, "pkg": "lambda/telemetry"})
	dispatcher.Logger = logrus.WithFields(logrus.Fields{"ext": extName, "pkg": "dispatcher"})
}

func main() {
	conf, err := GetConfig()
	if err != nil {
		Logger.Error(err)
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		s := <-sigs
		cancel()
		Logger.Info("received signal", s, "terminating")
	}()

	extCli := extension.NewClient(conf.AWSLambdaConfig.RuntimeApi)
	extID, err := extCli.Register(ctx, conf.AWSLambdaConfig.ExtensionName)
	if err != nil {
		Logger.Error(err)
		return
	}

	tlmListener := telemetry.NewTelemetryApiListener(conf.AWSLambdaConfig.IsSAMLocal)
	tlmListenerUri, err := tlmListener.Start()
	if err != nil {
		Logger.Error(err)
		return
	}

	tlmCli := telemetry.NewClient(conf.AWSLambdaConfig.RuntimeApi)
	if _, err = tlmCli.Subscribe(ctx, extID, tlmListenerUri); err != nil {
		Logger.Error(err)
		return
	}

	host, err := mackerel.CreateOrGetHost(&mackerel.CreateOrGetHostParam{
		MackerelApiKey: conf.MackerelConfig.ApiKey,
		RoleFullnames:  conf.MackerelConfig.RoleFullnames,
		FunctionName:   conf.AWSLambdaConfig.FunctionName,
		EnvironmentID:  conf.AWSLambdaConfig.EnvironmentID,
	})
	if err != nil {
		Logger.Error(err)
		return
	}

	if err := host.CreateGraphDefs(); err != nil {
		Logger.Error(err)
		return
	}

	dispatcher := dispatcher.NewDispatcher(host)

	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()
	go func() {
		for range ticker.C {
			dispatcher.Dispatch(ctx, tlmListener.LogEventsQueue, false)
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			Logger.Info("Waiting for next event...")

			// This is a blocking action
			res, err := extCli.NextEvent(ctx)
			if err != nil {
				Logger.Warning(err)
				return
			}

			// Dispatching log events from previous invocations
			dispatcher.Dispatch(ctx, tlmListener.LogEventsQueue, false)

			if res.EventType == extension.Shutdown {
				// Dispatch all remaining telemetry, handle shutdown
				Logger.Info("Shutdown event")
				dispatcher.Dispatch(ctx, tlmListener.LogEventsQueue, true)
				host.Retire()
				return
			}
		}
	}
}
