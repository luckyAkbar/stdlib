package cmd

import (
	"os"

	runtime "github.com/banzaicloud/logrus-runtime-formatter"
	"github.com/evalphobia/logrus_sentry"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// CobraInitializer returns a new cobra command
func CobraInitializer() *cobra.Command {
	return &cobra.Command{
		Use:   "cobra-example",
		Short: "An example of cobra",
		Long:  "This application shows how to create modern CLI applications in go using Cobra CLI library",
	}
}

// SetupLogger sets up the logger. should be called within init() function from your app
func SetupLogger(env, logLevel, sentryDSN string) {
	formatter := runtime.Formatter{
		ChildFormatter: &logrus.JSONFormatter{},
		Line:           true,
		File:           true,
	}

	if env == "development" || env == "local" {
		formatter = runtime.Formatter{
			ChildFormatter: &logrus.TextFormatter{
				ForceColors:   true,
				FullTimestamp: true,
			},
			Line: true,
			File: true,
		}
	}

	logrus.SetFormatter(&formatter)
	logrus.SetOutput(os.Stdout)

	logrusLevel, err := logrus.ParseLevel(logLevel)
	if err != nil {
		logrusLevel = logrus.DebugLevel
	}
	logrus.SetLevel(logrusLevel)

	hook, err := logrus_sentry.NewSentryHook(sentryDSN, []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
	})

	if err != nil {
		logrus.Info("Logger configured to use only local stdout")
		return
	}

	hook.SetEnvironment(env)
	hook.Timeout = 0 // fire and forget
	hook.StacktraceConfiguration.Enable = true
	logrus.AddHook(hook)
}
